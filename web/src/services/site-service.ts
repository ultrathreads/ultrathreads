// src/lib/services/site.service.ts
import { apiFetch, ApiBusinessError } from '@/lib/api/client';
import type { SiteConfig, SiteConfigRaw } from '@/lib/types/site';

/** 
 * ✅ 全局统一的兜底配置
 * 导出为常量，方便单元测试和特殊场景单独引用
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
  siteNotification: '',
  siteTips: null,
  siteIndexHtml: '',
  recommendTags: [],
  scoreConfig: { postTopicScore: 0, postCommentScore: 0 },
  defaultNodeId: 1,
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
    navLinks: raw.setting.siteNavs ?? [], // 👈 Service层不提供UI兜底，空则返回空数组
    siteNotification: raw.setting.siteNotification,
    siteTips: raw.setting.siteTips,
    recommendTags: raw.setting.recommendTags ?? [],
    scoreConfig: raw.setting.scoreConfig,
    defaultNodeId: raw.setting.defaultNodeId,
    // ✅ 补全遗漏字段的映射
    siteIndexHtml: raw.setting.siteIndexHtml, 
  };
}

/**
 * 获取站点配置
 * - 成功：返回真实数据
 * - 失败：打印错误日志 + 返回 Fallback，保证调用方永远拿到有效值
 */
export async function fetchSiteConfig(): Promise<SiteConfig> {
  try {
    const raw = await apiFetch<SiteConfigRaw>('/config/site-config', {
      cacheStrategy: { next: { revalidate: 300 } },
    });
    return mapToSiteConfig(raw);
  } catch (error) {
    if (error instanceof ApiBusinessError) {
      console.error(`[SiteService] Business error (code:${error.code}): ${error.message}`);
    } else {
      console.error('[SiteService] Fetch failed, using fallback:', error);
    }
    return FALLBACK_SITE_CONFIG;
  }
}