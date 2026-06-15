// src/services/thread-service.ts
import { apiFetch } from '@/lib/api/client';
import type { PostEntity, NodeEntity, UserEntity } from '@/types/domain';
import type { PaginationMeta } from '@/types/api';
import { DEFAULT_LIMIT } from '@/constants';
import { assembleSideload } from '@/lib/utils/assemble-sideload';
import type { IncludedData } from '@/lib/utils/assemble-sideload';

// ==================== 传输层类型 ====================

/**
 * 帖子列表项（传输模型）
 * 从 PostEntity 派生，仅保留列表页渲染必需字段
 */
export interface ThreadListItem {
  slug: string;
  threadSlug: string;
  parentSlug: string;
  title: string;
  createTime: number;
  lastCommentTime: number;
  viewCount: number;
  commentCount: number;

  // Sideload 外键（新接口返回）
  userSlug?: string;
  nodeSlug?: string;
  tagSlugs?: string[]; // 标签slug列表，后端保证返回 [] 而非 null
}

/** API 原始响应结构（不导出，仅内部使用） */
interface ThreadsApiResponse {
  data: ThreadListItem[];
  page: PaginationMeta;
  lastReadAtMap: Record<string, number>;
  included?: IncludedData;
}

// ==================== 视图层类型 ====================

export interface ThreadPageData {
  posts: ThreadListItem[];
  paging: PaginationMeta;
  error: string | null;
  lastReadAtMap: Record<string, number>;
}

// ==================== 服务函数 ====================

/**
 * 获取帖子列表页数据
 * @param page 当前页码
 * @param nodeSlug 可选的板块ID，传入时仅获取该板块下的帖子；未传或无效时默认为 0
 */
export async function getThreadPageData(
  page: number,
  nodeSlug?: string,
): Promise<ThreadPageData> {
  const safePage = Math.max(1, Number.isNaN(page) ? 1 : page);

  const params = new URLSearchParams({
    page: String(safePage),
    pageSize: String(DEFAULT_LIMIT),
  });

  // 根据 nodeSlug 是否存在，动态选择 RESTful 路径
  // 空值 → /threads（首页全局列表）
  // 有值 → /nodes/:nodeSlug/threads（板块列表）
  const basePath = nodeSlug
    ? `/nodes/${encodeURIComponent(nodeSlug)}/threads`
    : '/threads';

  // 缓存标签按节点隔离，避免切换板块时命中旧缓存
  const cacheTags = ['threads', ...(nodeSlug ? [`node-${nodeSlug}`] : [])];

  try {
    const rsp = await apiFetch<ThreadsApiResponse>(
      `${basePath}?${params.toString()}`,
      {
        auth: true,
        skipDataUnwrap: true,
        cacheStrategy: { next: { tags: cacheTags } },
        // cacheStrategy: { next: { tags: cacheTags, revalidate: 60 } },
      },
    );
    const assembledPosts = assembleSideload(rsp.data ?? [], rsp.included);

    return {
      posts: assembledPosts,
      paging: rsp.meta,
      lastReadAtMap: rsp.context?.lastReadAtMap ?? {},
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
 * @param tagSlug 标签ID
 * @param page 当前页码
 */
export async function getTagPageData(
  tagSlug: string,
  page: number,
): Promise<ThreadPageData> {
  const safePage = Math.max(1, Number.isNaN(page) ? 1 : page);

  const params = new URLSearchParams({
    page: String(safePage),
    pageSize: String(DEFAULT_LIMIT),
  });

  const cacheTags = ['threads', `tag-${tagSlug}`];

  try {
    const rsp = await apiFetch<ThreadsApiResponse>(
       `/tags/${encodeURIComponent(tagSlug)}/threads?${params.toString()}`,
      {
        auth: true,
        skipDataUnwrap: true,
        cacheStrategy: { next: { tags: cacheTags } },
      },
    );

    const assembledPosts = assembleSideload(rsp.data ?? [], rsp.included);

    return {
      posts: assembledPosts,
      paging: rsp.meta,
      lastReadAtMap: rsp.context?.lastReadAtMap ?? {},
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
