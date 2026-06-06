'use client';
import { createContext, useContext, useState, useEffect, useCallback, useRef } from 'react';
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

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<CurrentUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // ✅ 新增：标记是否已触发登出，防止竞态条件下的无效刷新
  const isLoggedOutRef = useRef(false);

  const refreshUser = useCallback(async () => {
    // ✅ 如果已经登出，直接跳过，避免无意义的 401 请求
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
      // ✅ 核心修复：401 是正常的未登录状态，不是系统错误
      if (err instanceof ApiBusinessError && err.message === 'UNAUTHORIZED') {
        console.info('[AuthProvider] Token expired or logged out, clearing session silently.');
        setUser(null);
        // ⚠️ 不调用 setError，避免登录页闪现错误提示
        return;
      }

      // 只有非 401 的真正异常才记录到全局错误状态
      const message = err instanceof Error ? err.message : String(err);
      console.error('[AuthProvider] refreshUser failed:', message, err);
      setError(message);
      setUser(null);
    } finally {
      // ✅ 仅在未登出时更新 loading 状态，避免覆盖跳转流程
      if (!isLoggedOutRef.current) {
        setIsLoading(false);
      }
    }
  }, []);

  // 首次挂载时自动获取用户信息
  useEffect(() => {
    refreshUser();
  }, [refreshUser]);

  const logout = useCallback(async () => {
    // 🔑 第一时间设置标记，阻断所有后续 refreshUser 调用
    isLoggedOutRef.current = true;

    try {
      await fetch('/api/auth/logout', { method: 'POST' });
    } finally {
      // 无论接口是否成功，前端都应重置内存状态并跳转
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
      logout
    }}>
      {children}
    </AuthContext.Provider>
  );
}

// 自定义 Hook，方便组件调用
export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
}