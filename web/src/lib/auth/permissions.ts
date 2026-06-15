// src/lib/auth/permissions.ts (或 src/utils/auth.ts)
// 🎯 专门存放鉴权相关的运行时逻辑
import type { CurrentUser } from '@/types/auth';

/**
 * 判断用户是否有权访问后台面板
 * 优先检查显式权限码，兜底检查 level >= 10（管理员）
 */
export function canAccessAdminPanel(user: CurrentUser | null): boolean {
  if (!user) return false;
  return user.permissions?.includes('admin:panel:access') === true 
    || user.level >= 10;
}

// 未来可以方便地扩展更多权限判断函数
// export function canEditPost(user: CurrentUser | null, post: Post): boolean { ... }