// src/services/my-post-service.ts
import { apiFetch } from '@/lib/api/client';
import type { UserEntity, NodeEntity } from '@/types/domain';
import type { PaginationMeta } from '@/types/api';
import { DEFAULT_LIMIT } from '@/constants';
import { assembleSideload } from '@/lib/utils/assemble-sideload';
import type { IncludedData } from '@/lib/utils/assemble-sideload';

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
  data: T[];
  page: PaginationMeta;
  lastReadAtMap: Record<string, number>;
  included?: IncludedData;
}

// ==================== 视图层类型 ====================

export interface UserPostsPageData<T> {
  posts: T[];
  paging: PaginationMeta;
  error: string | null;
}

// ==================== 服务函数 ====================

/**
 * 【内部底层函数】通用的用户帖子数据获取逻辑
 * 不对外暴露，仅用于复用请求、分页组装和错误处理
 */
async function fetchUserPostsByType<T>(
  userSlug: string,
  page: number,
  type: 'root' | 'reply',
  pageSize: number = DEFAULT_LIMIT,
): Promise<UserPostsPageData<T>> {
  const safePage = Math.max(1, Number.isNaN(page) ? 1 : page);

  const params = new URLSearchParams({
    page: String(safePage),
    pageSize: String(pageSize),
    type,
  });

  try {
    const rsp = await apiFetch<UserPostsApiResponse<T>>(
      `/users/${userSlug}/posts?${params.toString()}`,
      {
        auth: true,
        skipDataUnwrap: true,
        cacheStrategy: { next: { revalidate: 0 } },
      },
    );

    const assembledPosts = assembleSideload(rsp.data ?? [], rsp.included);
    return {
      posts: assembledPosts,
      paging: rsp.meta,
      error: null,
    };
  } catch (err) {
    console.error(`[UserPostService] Fetch ${type} posts failed:`, err);
    return {
      posts: [],
      paging: { currentPage: safePage, pageSize: pageSize, totalItems: 0 },
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
  pageSize?: number,
): Promise<UserPostsPageData<MyRootPostListItem>> {
  return fetchUserPostsByType<MyRootPostListItem>(userSlug, page, 'root', pageSize);
}

/**
 * 获取指定用户的【回帖】列表
 */
export function getUserReplyPostsPageData(
  userSlug: string,
  page: number,
  pageSize?: number,
): Promise<UserPostsPageData<MyReplyPostListItem>> {
  return fetchUserPostsByType<MyReplyPostListItem>(userSlug, page, 'reply', pageSize);
}