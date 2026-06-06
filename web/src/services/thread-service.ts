// src/lib/services/thread-service.ts
import { apiFetch } from '@/lib/api/client';

// ==================== 类型定义 ====================

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
interface ThreadsResponse {
  results: SimplePost[];
  page: PaginationInfo;
}

export interface ThreadPageData {
  posts: SimplePost[];
  paging: PaginationInfo;
  error: string | null;
}

// ==================== 服务函数 ====================

const DEFAULT_LIMIT = 20;

/**
 * 获取帖子列表页数据
 * - 成功：返回真实数据
 * - 失败：打印错误日志 + 返回空数据兜底，保证页面不白屏
 */
export async function getThreadPageData(page: number): Promise<ThreadPageData> {
  const safePage = Math.max(1, Number.isNaN(page) ? 1 : page);

  try {
    const data = await apiFetch<ThreadsResponse>(
      `/posts/threads?page=${safePage}&limit=${DEFAULT_LIMIT}`,
      {
        auth: false,
        cacheStrategy: { next: { tags: ['threads'], revalidate: 60 } },
      },
    );

    return {
      posts: data.results ?? [],
      paging: data.page,
      error: null,
    };
  } catch (err) {
    console.error('[ThreadService] Fetch failed:', err);
    return {
      posts: [],
      paging: { page: safePage, limit: DEFAULT_LIMIT, total: 0 },
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}