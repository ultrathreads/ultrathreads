// src/services/tag-service.ts
import { apiFetch } from '@/lib/api/client';
import type { TagEntity } from '@/types/domain';

export interface TagPageData {
  tags: TagEntity[];
  error: string | null;
}

/**
 * 获取单个板块详情
 * - 用于点击板块时获取名称和简介，供右侧面板展示
 */
export async function getTagDetail(tagSlug: string): Promise<TagDetailData> {
  try {
    const data = await apiFetch<tagEntity>(`/tags/${tagSlug}`, {
      auth: false,
      cacheStrategy: { next: { tags: [`tag-${tagSlug}`], revalidate: 30 } },
    });

    return {
      tag: data ?? null,
      error: null,
    };
  } catch (err) {
    console.error(`[NodeService] Fetch tag ${tagSlug} failed:`, err);
    return {
      tag: null,
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}

/**
 * 获取热门标签列表
 * - 失败时返回空数组兜底，保证侧边栏不白屏
 */
export async function getHotTags(pageSize = 20): Promise<TagPageData> {
  try {
    const data = await apiFetch<TagEntity[]>(`/tags/hot?pageSize=${pageSize}`, {
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

/**
 * 获取标签自动补全建议
 * - 纯客户端调用，无缓存
 * - 失败时返回空数组，保证输入框不报错
 */
export async function fetchTagSuggestions(input: string): Promise<TagEntity[]> {
  if (!input.trim()) return [];

  try {
    // ✅ 复用 apiFetch，自动解包 envelope.data
    // Go 后端使用 FormValue("input")，query string 即可满足
    const data = await apiFetch<TagEntity[]>(`/tags/auto-complete?input=${encodeURIComponent(input.trim())}`, {
      method: 'POST',
      auth: true,
      cacheStrategy: undefined, // ⚠️ 关键：禁用 Next.js 缓存，确保每次输入都实时请求
    });

    return Array.isArray(data) ? data : [];
  } catch (err) {
    console.error('[TagService] Auto-complete failed:', err);
    return [];
  }
}
