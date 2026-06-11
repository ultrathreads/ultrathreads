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
  lastReadAtMap: Record<string, number>;
}

// ==================== 视图层类型 ====================

export interface ThreadPageData {
  posts: ThreadListItem[];
  paging: PaginationMeta;
  error: string | null;
  lastReadAtMap: Record<string, number>;
}

// ==================== 服务函数 ====================

const DEFAULT_LIMIT = 20;

/**
 * 获取帖子列表页数据
 * @param page 当前页码
 * @param nodeId 可选的板块ID，传入时仅获取该板块下的帖子；未传或无效时默认为 0
 */
export async function getThreadPageData(
  page: number,
  nodeId?: number,
): Promise<ThreadPageData> {
  const safePage = Math.max(1, Number.isNaN(page) ? 1 : page);
  // ✅ 确保 nodeId 始终存在：undefined / NaN / 未传 → 0
  const safeNodeId =
    nodeId !== undefined && !Number.isNaN(nodeId) ? nodeId : 0;

  const params = new URLSearchParams({
    page: String(safePage),
    limit: String(DEFAULT_LIMIT),
    nodeId: String(safeNodeId), // ✅ 始终携带，不再条件判断
  });

  // 缓存标签按节点隔离，避免切换板块时命中旧缓存
  const cacheTags = ['threads', ...(safeNodeId ? [`node-${safeNodeId}`] : [])];

  try {
    const data = await apiFetch<ThreadsApiResponse>(
      `/threads?${params.toString()}`,
      {
        auth: true,
        cacheStrategy: { next: { tags: cacheTags } },
        // cacheStrategy: { next: { tags: cacheTags, revalidate: 60 } },
      },
    );

    return {
      posts: data.results ?? [],
      paging: data.page,
      lastReadAtMap: data.lastReadAtMap ?? {},
      error: null,
    };
  } catch (err) {
    console.error('[ThreadService] Fetch failed:', err);
    return {
      posts: [],
      paging: { currentPage: safePage, pageSize: DEFAULT_LIMIT, totalItems: 0 },
      lastReadAtMap: {},
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}

/**
 * 获取指定标签下的帖子列表页数据
 * @param tagId 标签ID
 * @param page 当前页码
 */
export async function getTagPageData(
  tagId: number,
  page: number,
): Promise<ThreadPageData> {
  const safePage = Math.max(1, Number.isNaN(page) ? 1 : page);
  const safeTagId = tagId;

  const params = new URLSearchParams({
    page: String(safePage),
    limit: String(DEFAULT_LIMIT),
    tagId: String(safeTagId),
  });

  // ✅ 缓存标签按标签ID隔离，避免切换标签时命中旧缓存
  const cacheTags = ['threads', `tag-${safeTagId}`];

  try {
    const data = await apiFetch<ThreadsApiResponse>(
      `/threads/tag?${params.toString()}`, // 对应你后端的 GetTagPosts 接口
      {
        auth: true,
        cacheStrategy: { next: { tags: cacheTags } },
      },
    );

    return {
      posts: data.results ?? [],
      paging: data.page,
      lastReadAtMap: data.lastReadAtMap ?? {},
      error: null,
    };
  } catch (err) {
    console.error('[ThreadService] Fetch tag posts failed:', err);
    return {
      posts: [],
      paging: { currentPage: safePage, pageSize: DEFAULT_LIMIT, totalItems: 0 },
      lastReadAtMap: {},
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}
