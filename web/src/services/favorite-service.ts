// src/services/favorite-service.ts
import { apiFetch } from '@/lib/api/client';
import type { UserEntity } from '@/types/domain';
import type { PaginationMeta } from '@/types/api';

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

/** 对应后端 DataEnvelope<T> */
interface FavoritesData {
  results: FavoriteItem[];
}

/** 对应后端完整 ListResponse 信封 */
interface FavoritesEnvelope {
  data: FavoritesData;
  meta: PaginationMeta;
  context?: Record<string, unknown>;
}

// ==================== 视图层类型 ====================
export interface FavoritesPageData {
  favorites: FavoriteItem[];
  paging: PaginationMeta;
  error: string | null;
}

// ==================== 服务函数 ====================
const DEFAULT_LIMIT = 20;

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
    pageSize: String(pageSize),
  });

  try {
    const envelope = await apiFetch<FavoritesEnvelope>('/users/me/favorites', {
      auth: true,
      skipDataUnwrap: true,
      params,
      cacheStrategy: { next: { revalidate: 0 } },
    });

    return {
      favorites: envelope.data.results ?? [],
      paging: envelope.meta,
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
  const params = new URLSearchParams({
    entityType,
    entityId: String(entityId),
  });

  await apiFetch(`/users/me/favorites?${params.toString()}`, {
    method: 'DELETE',
    auth: true,
  });
}