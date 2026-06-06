// src/lib/api/posts.ts
import { apiFetch } from './client';

/**
 * 后端 /api/posts/threads 返回的真实数据结构
 * ⚠️ 注意：results 中的元素是 SimplePost（由 converter.ToSimplePosts 生成）
 *    而非完整的 model.Post，字段以实际序列化结果为准
 */

export interface PostUser {
  id: number;
  username: string;
  email: string;
  nickname: string;
}

export interface SimplePost {
  id: number;
  parentId: number;
  threadId: number;
  user: PostUser; 
  title: string;
  content: string;
  createTime: number;
  updateTime: number;
  replyCount: number;
  viewCount: number;
  isPinned: boolean;
  isLocked: boolean;
}

export interface PaginationInfo {
  page: number;
  limit: number;
  total: number;
}

/** apiFetch 自动拆解 envelope.data 后的纯净业务数据 */
export interface ThreadsData {
  results: SimplePost[];
  page: PaginationInfo;
}

/**
 * 获取主帖及其扁平化回帖列表
 */
export async function fetchThreads(page: number, limit = 20): Promise<ThreadsData> {
  return apiFetch<ThreadsData>(`/posts/threads?page=${page}&limit=${limit}`, {
    auth: false,
    cacheStrategy: { next: { tags: ['threads'], revalidate: 60 } },
  });
}