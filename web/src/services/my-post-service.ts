// src/services/my-post-service.ts
import { apiFetch } from '@/lib/api/client';
import type { UserEntity, NodeEntity } from '@/types/domain';
import type { PaginationMeta } from '@/types/api';

// ==================== 传输层类型 ====================
export interface MyPostListItem {
  id: number;
  threadId: number;
  parentId: number | null;
  parentTitle: string | null;
  title: string;
  createTime: number;
  user: Pick<UserEntity, 'id' | 'username' | 'nickname'>;
  node: Pick<NodeEntity, 'nodeId' | 'name'> | null;
}

interface UserPostsApiResponse {
  results: MyPostListItem[];
  page: {
    page: number;
    limit: number;
    total: number;
  };
}

// ==================== 视图层类型 ====================
export interface UserPostsPageData {
  posts: MyPostListItem[];
  paging: PaginationMeta;
  error: string | null;
}

// ==================== 服务函数 ====================
const DEFAULT_LIMIT = 20;

/**
 * 获取指定用户的帖子列表
 * @param userId 目标用户ID（当前用户或他人）
 * @param page   页码
 * @param limit  每页条数
 */
export async function getUserPostsPageData(
  userId: number,
  page: number,
  limit: number = DEFAULT_LIMIT,
): Promise<UserPostsPageData> {
  const safePage = Math.max(1, Number.isNaN(page) ? 1 : page);

  const params = new URLSearchParams({
    page: String(safePage),
    limit: String(limit),
  });

  try {
    // ✅ 动态拼接 userId，复用同一接口
    const data = await apiFetch<UserPostsApiResponse>(
      `/user/posts/${userId}?${params.toString()}`,
      {
        auth: true, // 保持鉴权，后端可根据 token 判断是否返回敏感字段
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
    console.error('[UserPostService] Fetch failed:', err);
    return {
      posts: [],
      paging: { currentPage: safePage, pageSize: limit, totalItems: 0 },
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}