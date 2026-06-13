// src/services/post-service.ts
import { notFound } from 'next/navigation';
import { apiFetch, ApiBusinessError } from '@/lib/api/client';
import type { PostEntity, PostWithThread } from '@/types/domain';
import type { CreateRootPostPayload, UpdateRootPostPayload, CreateReplyPayload, CreatePostResponse } from '@/types/post'

interface PostServiceOptions {
  noCache?: boolean;
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
    if (error instanceof ApiBusinessError) {
      console.error(`[PostService] Biz Error: ${error.message} (code: ${error.code})`);
    }
    notFound();
  }
}

/**
 * ✅ 获取帖子详情及其所有回帖
 */
export async function getPostWithThread(
  postSlug: string,
  options?: PostServiceOptions
  ): Promise<PostWithThread> {
  try {
    return await apiFetch<PostWithThread>(`/posts/${postSlug}/tree`, {
      noCache: options?.noCache
    });
  } catch (error) {
    if (error instanceof ApiBusinessError) {
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
    return await apiFetch<PostWithFlat>(`/posts/${postSlug}/flat`);
  } catch (error) {
    if (error instanceof ApiBusinessError) {
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