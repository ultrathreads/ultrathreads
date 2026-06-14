// src/services/favorite-service.ts
import { apiFetch } from '@/lib/api/client';
import type { UserEntity } from '@/types/domain';
import type { PaginationMeta } from '@/types/api';
import { DEFAULT_LIMIT } from '@/constants';

// ==================== 传输层类型 ====================
export interface FavoriteItem {
  favoriteId: number;
  entityType: string;
  entityId: number;
  deleted: boolean;
  title: string;
  content: string;
  user: Pick<UserEntity, 'id' | 'username' | 'nickname' | 'avatar' | 'level' | 'levelName' | 'createTime'>;
  url: string;
  createTime: number;
}

interface FavoritesApiResponse {
  data: {
    page: {
      page: number;
      pageSize: number;
      total: number;
    };
    results: FavoriteItem[];
  };
}

// ==================== 视图层类型 ====================
export interface FavoritesPageData {
  favorites: FavoriteItem[];
  paging: PaginationMeta;
  error: string | null;
}

/**
 * 获取当前用户的书签列表
 * @param page 页码
 * @param pageSize 每页条数
 */
export async function getFavoritesPageData(
  page: number,
  pageSize: number = DEFAULT_LIMIT,
): Promise<FavoritesPageData> {
  const safePage = Math.max(1, Number.isNaN(page) ? 1 : page);

  const params = new URLSearchParams({
    page: String(safePage),
    pageSize: String(pageSize), // 根据后端代码，参数名为 pageSize
  });

  try {
    const data = await apiFetch<FavoritesApiResponse>('/user/favorites', {
      auth: true,
      params, // apiFetch 会自动处理 URLSearchParams
      cacheStrategy: { next: { revalidate: 0 } },
    });

    return {
      favorites: data.results ?? [],
      paging: {
        currentPage: data.page?.page ?? safePage,
        pageSize: data.page?.pageSize ?? pageSize,
        totalItems: data.page?.total ?? 0,
      },
      error: null,
    };
  } catch (err) {
    console.error('[FavoriteService] Fetch failed:', err);
    return {
      favorites: [],
      paging: { currentPage: safePage, pageSize: pageSize, totalItems: 0 },
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}

/**
 * 删除书签（服务端调用）
 * 注意：后端接收的是 entityType 和 entityId，而不是 favoriteId
 */
export async function deleteFavorite(entityType: string, entityId: number): Promise<void> {
  // 使用 DELETE 请求，参数通过 URLSearchParams 传递
  const params = new URLSearchParams({
    entityType,
    entityId: String(entityId),
  });

  await apiFetch(`/favorite/delete?${params.toString()}`, {
    method: 'DELETE',
    auth: true, // 必须携带 Token，因为后端需要获取当前用户
  });
}
