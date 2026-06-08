// src/services/post-service.ts
import { notFound } from 'next/navigation';
import { apiFetch, ApiBusinessError } from '@/lib/api/client';
import type { PostEntity, PostWithThread } from '@/types/domain';

// ✅ 修复：返回类型从 PostDetail 改为 PostEntity
export async function getPostDetail(postId: string): Promise<PostEntity> {
  try {
    return await apiFetch<PostEntity>(`/post/${postId}`);
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
export async function getPostWithThread(postId: string): Promise<PostWithThread> {
  try {
    return await apiFetch<PostWithThread>(`/post/${postId}/with-thread`);
  } catch (error) {
    if (error instanceof ApiBusinessError) {
      console.error(`[PostService] Biz Error: ${error.message} (code: ${error.code})`);
    }
    notFound();
  }
}

export interface CreatePostPayload {
  title: string;
  nodeId: number;
  parentId?: number;
  content: string;
  tags?: string[];
}

export interface CreatePostResponse {
  id: number;
}

/**
 * 创建新帖子
 */
export async function createPost(payload: CreatePostPayload): Promise<CreatePostResponse> {
  return apiFetch<CreatePostResponse>('/posts', {
    method: 'POST',
    auth: true,
    body: JSON.stringify(payload),
  });
}

// ==================== 👇 新增：点赞 & 收藏 ====================

/**
 * 点赞帖子
 * POST /post/:id/like
 */
export async function likePost(postId: string | number): Promise<void> {
  await apiFetch<null>(`/post/${postId}/like`, {
    method: 'POST',
    auth: true,
    cacheStrategy: undefined, // ✅ POST 写操作禁用缓存
  });
}

/**
 * 收藏帖子
 * POST /post/:id/favorite
 */
export async function favoritePost(postId: string | number): Promise<void> {
  await apiFetch<null>(`/post/${postId}/favorite`, {
    method: 'POST',
    auth: true,
    cacheStrategy: undefined,
  });
}