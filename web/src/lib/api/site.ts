// src/lib/api/site.ts
import { apiFetch, ApiBusinessError } from './client';
import type { SiteConfig } from '@/lib/types/site';

/** ⚠️ 仅描述 envelope.data 内部的结构，不包含 code/message */
interface SiteConfigRaw {
  appinfo: {
    name: string;
    version: string;
    user_level_admin: number;
  };
  setting: {
    siteTitle: string;
    siteDescription: string;
    siteKeywords: string | null;
    siteNavs: { label: string; href: string }[] | null;
    siteTips: string | null;
    siteNotification: string;
    siteIndexHtml: string;       // 👈 补充遗漏字段
    recommendTags: string[] | null;
    scoreConfig: { postTopicScore: number; postCommentScore: number };
    defaultNodeId: number;
  };
}

const FALLBACK_CONFIG: SiteConfig = {
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
  recommendTags: [],
  scoreConfig: { postTopicScore: 0, postCommentScore: 0 },
  defaultNodeId: 1,
};

function mapToSiteConfig(raw: SiteConfigRaw): SiteConfig {
  return {
    appName: raw.appinfo.name,
    appVersion: raw.appinfo.version,
    adminLevel: raw.appinfo.user_level_admin,
    siteTitle: raw.setting.siteTitle,
    siteDescription: raw.setting.siteDescription,
    siteKeywords: raw.setting.siteKeywords,
    navLinks: raw.setting.siteNavs ?? [
      { label: '首页', href: '/' },
      { label: '关于', href: '/about' },
    ],
    siteNotification: raw.setting.siteNotification,
    siteTips: raw.setting.siteTips,
    recommendTags: raw.setting.recommendTags ?? [],
    scoreConfig: raw.setting.scoreConfig,
    defaultNodeId: raw.setting.defaultNodeId,
  };
}

export async function getSiteConfig(): Promise<SiteConfig> {
  try {
    // ✅ apiFetch 已自动拆信封，这里拿到的直接是 SiteConfigRaw
    const raw = await apiFetch<SiteConfigRaw>('/config/site-config', {
      cacheStrategy: { next: { revalidate: 300 } },
    });
    return mapToSiteConfig(raw);
  } catch (error) {
    // 💡 可以精确区分业务错误和网络错误
    if (error instanceof ApiBusinessError) {
      console.error(`[getSiteConfig] Business error (code:${error.code}): ${error.message}`);
    } else {
      console.error('[getSiteConfig] Fetch failed, using fallback:', error);
    }
    return FALLBACK_CONFIG;
  }
}