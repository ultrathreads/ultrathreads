// src/services/tag-service.ts
import { apiFetch } from '@/lib/api/client';
import type { TagEntity } from '@/types/domain';

export interface TagPageData {
  tags: TagEntity[];
  error: string | null;
}

/**
 * 获取热门标签列表
 * - 失败时返回空数组兜底，保证侧边栏不白屏
 */
export async function getHotTags(limit = 20): Promise<TagPageData> {
  try {
    const data = await apiFetch<TagEntity[]>(`/tags/hot?limit=${limit}`, {
      auth: false,
      cacheStrategy: { next: { tags: ['tags'], revalidate: 120 } },
    });

    return {
      tags: Array.isArray(data) ? data : [],
      error: null,
    };
  } catch (err) {
    console.error('[TagService] Fetch hot tags failed:', err);
    return {
      tags: [],
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}