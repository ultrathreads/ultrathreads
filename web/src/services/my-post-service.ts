// src/services/my-post-service.ts
import { apiFetch } from '@/lib/api/client';
import type { UserEntity, NodeEntity } from '@/types/domain';
import type { PaginationMeta } from '@/types/api';

// ==================== 传输层类型 ====================

/**
 * 根帖列表项
 */
export interface MyRootPostListItem {
  slug: string;
  threadSlug: string;
  title: string;
  createTime: number;
  user: Pick<UserEntity, 'slug' | 'username' | 'nickname'>;
  node: Pick<NodeEntity, 'slug' | 'name'> | null;
}

/**
 * 回帖列表项
 */
export interface MyReplyPostListItem {
  slug: string;
  threadSlug: string;
  parentSlug: string | null;
  parentTitle: string | null;
  content: string;
  createTime: number;
  user: Pick<UserEntity, 'slug' | 'username' | 'nickname'>;
  node: Pick<NodeEntity, 'slug' | 'name'> | null;
}

/**
 * 通用 API 响应结构（支持泛型）
 */
interface UserPostsApiResponse<T> {
  results: T[];
  page: {
    page: number;
    limit: number;
    total: number;
  };
}

// ==================== 视图层类型 ====================

export interface UserPostsPageData<T> {
  posts: T[];
  paging: PaginationMeta;
  error: string | null;
}

// ==================== 服务函数 ====================

const DEFAULT_LIMIT = 20;

/**
 * 【内部底层函数】通用的用户帖子数据获取逻辑
 * 不对外暴露，仅用于复用请求、分页组装和错误处理
 */
async function fetchUserPostsByType<T>(
  userSlug: string,
  page: number,
  type: 'root' | 'reply',
  limit: number = DEFAULT_LIMIT,
): Promise<UserPostsPageData<T>> {
  const safePage = Math.max(1, Number.isNaN(page) ? 1 : page);

  const params = new URLSearchParams({
    page: String(safePage),
    limit: String(limit),
    type,
  });

  try {
    const data = await apiFetch<UserPostsApiResponse<T>>(
      `/users/${userSlug}/posts?${params.toString()}`,
      {
        auth: true,
        cacheStrategy: { next: { revalidate: 0 } },
      },
    );

    return {
      posts: data.results ?? [],
      paging: {
        currentPage: data.page?.page ?? safePage,
        pageSize: data.page?.limit ?? limit,
        totalItems: data.page?.total ?? 0,
      },
      error: null,
    };
  } catch (err) {
    console.error(`[UserPostService] Fetch ${type} posts failed:`, err);
    return {
      posts: [],
      paging: { currentPage: safePage, pageSize: limit, totalItems: 0 },
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}

/**
 * 获取指定用户的【根帖】列表
 */
export function getUserRootPostsPageData(
  userSlug: string,
  page: number,
  limit?: number,
): Promise<UserPostsPageData<MyRootPostListItem>> {
  return fetchUserPostsByType<MyRootPostListItem>(userSlug, page, 'root', limit);
}

/**
 * 获取指定用户的【回帖】列表
 */
export function getUserReplyPostsPageData(
  userSlug: string,
  page: number,
  limit?: number,
): Promise<UserPostsPageData<MyReplyPostListItem>> {
  return fetchUserPostsByType<MyReplyPostListItem>(userSlug, page, 'reply', limit);
}