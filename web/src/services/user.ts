// src/services/user.ts
import { apiFetch } from '@/lib/api/client';
import type { CurrentUser } from '@/types/auth';

/**
 * 用户资料更新请求体
 * ✅ 直接从 CurrentUser 派生，与 ProfilePage 的 ProfileForm 结构完全一致
 */
export type UserUpdatePayload = Pick<CurrentUser, 'nickname' | 'avatar' | 'website' | 'description'>;

/**
 * 更新指定用户的资料
 *
 * @param userId - 用户 ID，用于拼接 /users/:id 路由
 * @param payload - 更新数据
 *
 * - 后端通过 JWT 校验权限，但路由仍需合法的数字 ID
 * - 返回值为 null（c.Success(ctx, nil)）
 * - cache: 'no-store' 确保写入操作绕过缓存
 */
export async function updateUserProfile(
  userId: number,
  payload: UserUpdatePayload,
): Promise<void> {
  await apiFetch<null>(`/users/${userId}`, {
    method: 'PUT',
    auth: true,
    body: JSON.stringify(payload),
    cacheStrategy: { cache: 'no-store' },
  });
}