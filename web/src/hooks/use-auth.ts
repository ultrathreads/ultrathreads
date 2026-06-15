// src/hooks/use-auth.ts
import { useContext } from 'react';
import { AuthContext } from '@/providers/AuthProvider';

export function useAuth() {
  const ctx = useContext(AuthContext);

  if (!ctx) {
    throw new Error('useAuth must be used within an AuthProvider');
  }

  const displayName = ctx.user?.nickname || ctx.user?.username || '';
  const avatarUrl = ctx.user?.avatar || '';
  
  // ✅ 集中管理角色相关的派生状态
  const roles = ctx.user?.roles ?? [];
  const isAdmin = roles.includes('admin');

  return {
    ...ctx,
    displayName,
    avatarUrl,
    roles,      // 原始角色数组
    isAdmin,    // 便捷的管理员标识
  };
}