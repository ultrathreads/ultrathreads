import { notFound } from 'next/navigation';
import { apiFetch, ApiError } from '@/lib/api/client';
import type { UserEntity } from '@/types/domain';

/**
 * 获取用户公开信息（通过用户ID）
 * GET /user/:id
 */
export async function getUserById(userId: string | number): Promise<UserEntity> {
  try {
    return await apiFetch<UserEntity>(`/user/${userId}`);
  } catch (error) {
    if (error instanceof ApiError) {
      console.error(`[UserService] Biz Error: ${error.message} (code: ${error.code})`);
    }
    notFound();
  }
}

/**
 * 获取用户公开信息（通过 Slug / 用户名）
 * GET /user/slug/:slug
 */
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

/**
 * 获取当前登录用户的个人信息（需鉴权）
 * GET /user/me
 */
export async function getCurrentUser(): Promise<UserEntity> {
  try {
    return await apiFetch<UserEntity>('/user/me', {
      auth: true, // ✅ 需要携带 Token
    });
  } catch (error) {
    if (error instanceof ApiError) {
      console.error(`[UserService] Biz Error: ${error.message} (code: ${error.code})`);
    }
    notFound();
  }
}

/**
 * 更新当前用户个人信息（需鉴权）
 * PUT /user/me
 */
export async function updateCurrentUser(payload: Partial<UserEntity>): Promise<UserEntity> {
  try {
    return await apiFetch<UserEntity>('/user/me', {
      method: 'PUT',
      auth: true,
      body: JSON.stringify(payload),
      cacheStrategy: undefined, // ✅ 写操作禁用缓存
    });
  } catch (error) {
    if (error instanceof ApiError) {
      console.error(`[UserService] Biz Error: ${error.message} (code: ${error.code})`);
    }
    throw error; // 更新失败抛出异常，由表单组件捕获并提示
  }
}