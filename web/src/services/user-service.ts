// src/services/user.ts
import { notFound } from 'next/navigation';
import { apiFetch, ApiError } from '@/lib/api/client';
import type { UserEntity } from '@/types/domain';
import type { CurrentUser } from '@/types/auth';

// ============ Types ============

/** 用户资料更新请求体 */
export type UserUpdatePayload = Pick<CurrentUser, 'nickname' | 'avatar' | 'website' | 'description'>;

/** GET /profile/:slug - 通过用户名获取公开信息 */
export async function getUserBySlug(slug: string): Promise<UserEntity> {
  try {
    return await apiFetch<UserEntity>(`/profile/${slug}`);
  } catch (error) {
    if (error instanceof ApiError) {
      console.error(`[UserService] Biz Error: ${error.message} (code: ${error.code})`);
    }
    notFound();
  }
}

/** GET /user/me - 获取当前登录用户信息 */
export async function getCurrentUser(): Promise<UserEntity> {
  try {
    return await apiFetch<UserEntity>('/user/me', { auth: true });
  } catch (error) {
    if (error instanceof ApiError) {
      console.error(`[UserService] Biz Error: ${error.message} (code: ${error.code})`);
    }
    notFound();
  }
}

// ============ Write ============

/** PUT /user/me - 更新当前用户资料 */
export async function updateCurrentUser(payload: UserUpdatePayload): Promise<void> {
  try {
    await apiFetch<null>('/user/me', {
      method: 'PUT',
      auth: true,
      body: JSON.stringify(payload),
      cacheStrategy: { cache: 'no-store' },
    });
  } catch (error) {
    if (error instanceof ApiError) {
      console.error(`[UserService] Update Error: ${error.message} (code: ${error.code})`);
    }
    throw error; // 写操作抛出，由表单组件捕获提示
  }
}

/** PUT /users/:slug - 更新指定用户资料（管理员场景） */
export async function updateUserProfile(userSlug: string, payload: UserUpdatePayload): Promise<void> {
  try {
    await apiFetch<null>(`/users/${userSlug}`, {
      method: 'PUT',
      auth: true,
      body: JSON.stringify(payload),
      cacheStrategy: { cache: 'no-store' },
    });
  } catch (error) {
    if (error instanceof ApiError) {
      console.error(`[UserService] Update Error: ${error.message} (code: ${error.code})`);
    }
    throw error;
  }
}