// src/lib/utils/back-context.ts
import { headers } from 'next/headers';

export interface BackContext {
  backUrl: string;
}

/**
 * 从当前请求的 Referer 中智能提取返回链接
 * 仅当 Referer 是本站列表页时才使用，否则回退到首页
 */
export async function getBackContext(): Promise<BackContext> {
  try {
    const headerStore = await headers();
    const referer = headerStore.get('referer') || '';
    const host = headerStore.get('host') || '';

    // 安全检查：只信任同域 Referer
    if (!referer || !referer.includes(host)) {
      return { backUrl: '/' };
    }

    const url = new URL(referer);
    const pathname = url.pathname;

    // ✅ 白名单：只有这些路径才是合法的"列表页"
    const listPagePatterns = [
      /^\/$/,             // 首页
      /^\/nodes\/[^/]+$/, // 节点列表
      /^\/tags\/[^/]+$/,  // 标签列表
      /^\/my$/,           // 我的页面
    ];

    const isListPage = listPagePatterns.some((p) => p.test(pathname));

    if (isListPage) {
      // 仅透传已知的列表页状态参数，防止无关参数泄漏
      const ALLOWED_PARAMS = ['page', 'tab'];
      const params = new URLSearchParams();

      for (const key of ALLOWED_PARAMS) {
        const value = url.searchParams.get(key);
        if (value !== null) {
          params.set(key, value);
        }
      }

      const queryString = params.toString();
      return {
        backUrl: queryString ? `${pathname}?${queryString}` : pathname,
      };
    }
  } catch {
    // SSR 环境异常时安全降级
  }

  return { backUrl: '/' };
}