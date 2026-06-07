// src/hooks/use-auth.ts
import { useContext } from 'react';
import { AuthContext } from '@/providers/AuthProvider';

export function useAuth() {
  const ctx = useContext(AuthContext);

  if (!ctx) {
    throw new Error('useAuth must be used within an AuthProvider');
  }

  // 💡 集中管理派生状态，一处修改全局生效
  const displayName = ctx.user?.nickname || ctx.user?.username || '';
  const avatarUrl = ctx.user?.avatar
    || `https://api.dicebear.com/7.x/avataaars/svg?seed=${ctx.user?.username}`;

  return {
    ...ctx,
    displayName,
    avatarUrl,
  };
}