import { ApiResponse, SiteConfigData, SiteConfig } from '@/types/config';

const GO_API_BASE = process.env.GO_API_BASE_URL || 'http://localhost:8080';

/**
 * 将 Go API 原始响应映射为前端友好的扁平结构
 */
function mapToSiteConfig(raw: SiteConfigData): SiteConfig {
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

/** 默认降级配置 */
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

/**
 * 从 Go 后端获取站点配置（SSR 专用）
 */
export async function getSiteConfig(): Promise<SiteConfig> {
  try {
    const res = await fetch(`${GO_API_BASE}/api/config/site-config`, {
      next: { revalidate: 60 },
      headers: { 'Content-Type': 'application/json' },
    });

    if (!res.ok) {
      throw new Error(`HTTP ${res.status}: ${res.statusText}`);
    }

    const json: ApiResponse<SiteConfigData> = await res.json();

    // ✅ 校验业务状态码
    if (json.code !== 0) {
      throw new Error(`API business error [${json.code}]: ${json.message}`);
    }

    return mapToSiteConfig(json.data);
  } catch (err) {
    console.error('[getSiteConfig] Failed:', err);
    return FALLBACK_CONFIG;
  }
}