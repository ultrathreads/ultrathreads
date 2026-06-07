// src/services/thread-service.ts
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
 * @param page 当前页码
 * @param nodeId 可选的板块ID，传入时仅获取该板块下的帖子
 */
export async function getThreadPageData(
  page: number,
  nodeId?: number,
): Promise<ThreadPageData> {
  const safePage = Math.max(1, Number.isNaN(page) ? 1 : page);

  // ✅ 动态构建查询参数，有 nodeId 时追加，无则保持原样
  const params = new URLSearchParams({
    page: String(safePage),
    limit: String(DEFAULT_LIMIT),
  });
  if (nodeId !== undefined && !Number.isNaN(nodeId)) {
    params.set('nodeId', String(nodeId));
  }

  // ✅ 缓存标签按节点隔离，避免切换板块时命中旧缓存
  const cacheTags = ['threads', ...(nodeId ? [`node-${nodeId}`] : [])];

  try {
    const data = await apiFetch<ThreadsResponse>(
      `/posts/threads?${params.toString()}`,
      {
        auth: false,
        cacheStrategy: { next: { tags: cacheTags, revalidate: 60 } },
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