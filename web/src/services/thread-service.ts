// src/services/thread-service.ts
import { apiFetch } from '@/lib/api/client';
import type { PostEntity, NodeEntity, UserEntity } from '@/types/domain';
import type { PaginationMeta } from '@/types/api';

// ==================== 传输层类型 ====================

/**
 * 帖子列表项（传输模型）
 * 从 PostEntity 派生，仅保留列表页渲染必需字段
 */
export interface ThreadListItem {
  id: number;
  threadId: number;
  parentId: number;
  title: string;
  createTime: number;
  lastCommentTime: number;
  viewCount: number;
  commentCount: number;
  user: Pick<UserEntity, 'id' | 'username' | 'nickname' | 'avatar'>;
  node: Pick<NodeEntity, 'nodeId' | 'name'>;
}

/** API 原始响应结构（不导出，仅内部使用） */
interface ThreadsApiResponse {
  results: ThreadListItem[];
  page: PaginationMeta;
}

// ==================== 视图层类型 ====================

export interface ThreadPageData {
  posts: ThreadListItem[];
  paging: PaginationMeta;
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

  const params = new URLSearchParams({
    page: String(safePage),
    limit: String(DEFAULT_LIMIT),
  });
  if (nodeId !== undefined && !Number.isNaN(nodeId)) {
    params.set('nodeId', String(nodeId));
  }

  // 缓存标签按节点隔离，避免切换板块时命中旧缓存
  const cacheTags = ['threads', ...(nodeId ? [`node-${nodeId}`] : [])];

  try {
    const data = await apiFetch<ThreadsApiResponse>(
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
      paging: { currentPage: safePage, pageSize: DEFAULT_LIMIT, totalItems: 0 },
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}