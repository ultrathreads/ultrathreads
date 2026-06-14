// src/lib/services/site-service.ts
import { unstable_cache } from 'next/cache';
import { apiFetch, ApiError } from '@/lib/api/client';
import type { SiteConfig, SiteConfigRaw } from '@/lib/types/site';

/** 
 * 全局统一的兜底配置
 */
export const FALLBACK_SITE_CONFIG: SiteConfig = {
  appName: 'UltraThreads',
  appVersion: '0.0.0-unknown',
  adminLevel: 0,
  siteTitle: 'UltraThreads',
  siteDescription: '小而美的开发者社区',
  siteKeywords: null,
  navLinks: [
    { label: '首页', href: '/' },
    { label: '关于', href: '/about' },
  ],
  defaultNodeId: 1,
  recommendTags: [],
};

/** 将后端原始数据转换为前端标准结构 */
function mapToSiteConfig(raw: SiteConfigRaw): SiteConfig {
  return {
    appName: raw.appinfo.name,
    appVersion: raw.appinfo.version,
    adminLevel: raw.appinfo.user_level_admin,
    siteTitle: raw.setting.siteTitle,
    siteDescription: raw.setting.siteDescription,
    siteKeywords: raw.setting.siteKeywords,
    navLinks: raw.setting.siteNavs ?? [],
    defaultNodeId: raw.setting.defaultNodeId,
    recommendTags: raw.setting.recommendTags ?? [],
  };
}

const getCachedSiteConfig = unstable_cache(
  async (): Promise<SiteConfig> => {
    try {
      const raw = await apiFetch<SiteConfigRaw>('/site/config');
      return mapToSiteConfig(raw);
    } catch (error) {
      if (error instanceof ApiError) {
        if (error.code !== undefined) {
          console.error(
            `[SiteService] Business error [HTTP ${error.status}] (code:${error.code}): ${error.message}`
          );
        } else {
          console.error(
            `[SiteService] Infrastructure error [HTTP ${error.status}]: ${error.message}`
          );
        }
      } else {
        // 理论上不会进入此分支，保留作为类型安全兜底
        console.error('[SiteService] Unexpected error, using fallback:', error);
      }
      return FALLBACK_SITE_CONFIG;
    }
  },
  ['global-site-config'],
  { revalidate: 300, tags: ['site-config'] }
);

export async function fetchSiteConfig(): Promise<SiteConfig> {
  return getCachedSiteConfig();
}