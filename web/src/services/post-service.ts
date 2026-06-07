// src/services/post-service.ts
import { notFound } from 'next/navigation';
import { apiFetch, ApiBusinessError } from '@/lib/api/client';
import type { PostEntity, PostWithThread } from '@/types/domain';

// ✅ 已有：获取主帖详情
export async function getPostDetail(postId: string): Promise<PostDetail> {
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
 * 对应接口: GET /post/{id}/with-thread
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
  content: string;       // Markdown 原文
  tags?: string[];
}

export interface CreatePostResponse {
  id: number;
  // ...其他后端返回的字段
}

/**
 * 创建新帖子
 */
export async function createPost(payload: CreatePostPayload): Promise<CreatePostResponse> {
  return apiFetch<CreatePostResponse>('/posts', {
    method: 'POST',
    auth: true,          // ✅ 发帖必须登录
    body: JSON.stringify(payload),
  });
}
