// src/providers/AuthProvider.tsx
'use client';

import { createContext, useState, useCallback, useRef } from 'react';
import { getCurrentUser, type CurrentUser } from '@/services/auth.server';
import { ApiError } from '@/lib/api/client';
import { toast } from 'sonner'; 

interface AuthContextType {
  user: CurrentUser | null;
  isLoggedIn: boolean;
  isLoading: boolean;
  refreshUser: () => Promise<void>;
  logout: () => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  /** ✅ 关键：服务端注入的初始用户状态 */
  initialUser: CurrentUser | null;
  children: React.ReactNode;
}

export function AuthProvider({ initialUser, children }: AuthProviderProps) {
  // ✅ 直接使用服务端注入的状态作为初始值
  const [user, setUser] = useState<CurrentUser | null>(initialUser);
  const [isLoading, setIsLoading] = useState(false); // 初始不再需要 loading
  const isLoggedOutRef = useRef(false);

  const refreshUser = useCallback(async () => {
    if (isLoggedOutRef.current) return;

    try {
      setIsLoading(true);
      const userInfo = await getCurrentUser();
      setUser(userInfo ?? null);
    } catch (err) {
      if (err instanceof ApiError && err.message === 'UNAUTHORIZED') {
        console.debug('[AuthProvider] Session expired, clearing user.');
        setUser(null);
        return;
      }
      console.error('[AuthProvider] refreshUser failed:', err);
      setUser(null);
    } finally {
      if (!isLoggedOutRef.current) setIsLoading(false);
    }
  }, []);

  const logout = useCallback(async () => {
    isLoggedOutRef.current = true;
    try {
      await fetch('/api/auth/logout', { method: 'POST' });
    } finally {
      setUser(null);
      toast.success('已成功退出登录');
      window.location.href = '/';
    }
  }, []);

  return (
    <AuthContext.Provider value={{
      user,
      isLoggedIn: !!user,
      isLoading,
      refreshUser,
      logout,
    }}>
      {children}
    </AuthContext.Provider>
  );
}