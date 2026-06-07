// src/providers/AuthProvider.tsx
'use client';

import { createContext, useState, useEffect, useCallback, useRef } from 'react';
import { getCurrentUser, type CurrentUser } from '@/services/auth-service';
import { ApiBusinessError } from '@/lib/api/client';

interface AuthContextType {
  user: CurrentUser | null;
  isLoggedIn: boolean;
  isLoading: boolean;
  error: string | null;
  refreshUser: () => Promise<void>;
  logout: () => void;
}

// ✅ 改为具名导出，供 hooks/use-auth.ts 导入
export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<CurrentUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const isLoggedOutRef = useRef(false);

  const refreshUser = useCallback(async () => {
    if (isLoggedOutRef.current) return;

    try {
      setIsLoading(true);
      setError(null);
      console.log('[AuthProvider] Fetching current user...');

      const userInfo = await getCurrentUser();
      console.log('[AuthProvider] Got user:', userInfo);

      if (!userInfo) {
        throw new Error('getCurrentUser returned null/undefined after unwrap');
      }

      setUser(userInfo);
    } catch (err) {
      if (err instanceof ApiBusinessError && err.message === 'UNAUTHORIZED') {
        console.info('[AuthProvider] Token expired or logged out, clearing session silently.');
        setUser(null);
        return;
      }

      const message = err instanceof Error ? err.message : String(err);
      console.error('[AuthProvider] refreshUser failed:', message, err);
      setError(message);
      setUser(null);
    } finally {
      if (!isLoggedOutRef.current) {
        setIsLoading(false);
      }
    }
  }, []);

  useEffect(() => {
    refreshUser();
  }, [refreshUser]);

  const logout = useCallback(async () => {
    isLoggedOutRef.current = true;

    try {
      await fetch('/api/auth/logout', { method: 'POST' });
    } finally {
      setUser(null);
      setError(null);
      window.location.href = '/';
    }
  }, []);

  return (
    <AuthContext.Provider value={{
      user,
      isLoggedIn: !!user,
      isLoading,
      error,
      refreshUser,
      logout,
    }}>
      {children}
    </AuthContext.Provider>
  );
}
