// src/services/post-service.ts
import { notFound } from 'next/navigation';
import { apiFetch, ApiError } from '@/lib/api/client';
import type { PostEntity, PostWithTree } from '@/types/domain';
import type { CreateRootPostPayload, UpdateRootPostPayload, CreateReplyPayload, CreatePostResponse } from '@/types/post'

interface PostServiceOptions {
  noCache?: boolean;
  /** 为 true 时，请求失败不会触发 notFound()，而是将错误重新抛出 */
  throwOnError?: boolean;
}

export async function getPostDetail(
  postSlug: string,
  options?: PostServiceOptions
  ): Promise<PostEntity> {
  try {
    return await apiFetch<PostEntity>(`/posts/${postSlug}`, {
      noCache: options?.noCache
    });
  } catch (error) {
    if (error instanceof ApiError) {
      console.error(`[PostService] Biz Error: ${error.message} (code: ${error.code})`);
    }
    // 如果调用方选择自行处理错误，则重新抛出
    if (options?.throwOnError) {
      throw error;
    }
    notFound();
  }
}

/**
 * ✅ 获取帖子详情及其所有回帖
 */
export async function getPostTree(
  postSlug: string,
  options?: PostServiceOptions
  ): Promise<PostWithTree> {
  try {
    return await apiFetch<PostWithTree>(`/posts/${postSlug}/tree`, {
      noCache: options?.noCache,
      skipDataUnwrap: true,
    });
  } catch (error) {
    if (error instanceof ApiError) {
      console.error(`[PostService] Biz Error: ${error.message} (code: ${error.code})`);
    }
    notFound();
  }
}

/**
 * ✅ 获取帖子详情及其所有回帖（扁平列表）
 */
export async function getPostFlat(postSlug: string): Promise<PostWithFlat> {
  try {
    return await apiFetch<PostWithFlat>(`/posts/${postSlug}/flat`, {
      skipDataUnwrap: true,
    });
  } catch (error) {
    if (error instanceof ApiError) {
      console.error(`[PostService] Biz Error: ${error.message} (code: ${error.code})`);
    }
    notFound();
  }
}

/**
 * 创建根帖（主帖）
 * POST /api/posts
 */
export async function createRootPost(
  payload: CreateRootPostPayload
): Promise<CreatePostResponse> {
  return apiFetch<CreatePostResponse>('/posts', {
    method: 'POST',
    auth: true,
    skipDataUnwrap: true,
    body: JSON.stringify(payload),
  });
}

/** 更新根帖 */
export async function updateRootPost(
  slug: string,
  payload: UpdateRootPostPayload
): Promise<CreatePostResponse> {
  return apiFetch<CreatePostResponse>(`/posts/${slug}`, {
    method: 'POST',
    auth: true,
    skipDataUnwrap: true,
    body: JSON.stringify(payload),
  });
}

/**
 * 创建回复
 * POST /api/posts/:parentSlug/replies
 */
export async function createReply(
  parentSlug: string,
  payload: CreateReplyPayload
): Promise<CreatePostResponse> {
  return apiFetch<CreatePostResponse>(`/posts/${parentSlug}/replies`, {
    method: 'POST',
    auth: true,
    body: JSON.stringify(payload),
  });
}

/** 更新回帖 */
export async function updateReply(
  slug: string,
  payload: UpdateReplyPayload
): Promise<CreatePostResponse> {
  return apiFetch<CreatePostResponse>(`/replies/${slug}`, {
    method: 'POST',
    auth: true,
    body: JSON.stringify(payload),
  });
}

// ==================== 👇 新增：点赞 & 收藏 ====================

/**
 * 点赞帖子
 * POST /post/:slug/like
 */
export async function likePost(postSlug: string | number): Promise<void> {
  await apiFetch<null>(`/posts/${postSlug}/like`, {
    method: 'POST',
    auth: true,
    cacheStrategy: undefined,
  });
}

/**
 * 收藏帖子
 * POST /post/:slug/favorite
 */
export async function favoritePost(postSlug: string | number): Promise<void> {
  await apiFetch<null>(`/posts/${postSlug}/favorite`, {
    method: 'POST',
    auth: true,
    cacheStrategy: undefined,
  });
}