/** Go API 统一响应包装 */
export interface ApiResponse<T> {
  code: number;
  data: T;
  message: string;
}

/** 应用信息 */
export interface AppInfo {
  name: string;
  version: string;
  user_level_admin: number;
}

/** 积分配置 */
export interface ScoreConfig {
  postTopicScore: number;
  postCommentScore: number;
}

/** 站点设置 */
export interface SiteSetting {
  siteTitle: string;
  siteDescription: string;
  siteKeywords: string | null;
  siteNavs: { label: string; href: string }[] | null;
  siteTips: string | null;
  siteNotification: string;
  siteIndexHtml: string;
  recommendTags: string[] | null;
  scoreConfig: ScoreConfig;
  defaultNodeId: number;
}

/** data 字段完整结构 */
export interface SiteConfigData {
  appinfo: AppInfo;
  setting: SiteSetting;
}

/** 前端使用的扁平化配置（从原始数据映射而来） */
export interface SiteConfig {
  // 来自 appinfo
  appName: string;
  appVersion: string;
  adminLevel: number;
  // 来自 setting
  siteTitle: string;
  siteDescription: string;
  siteKeywords: string | null;
  navLinks: { label: string; href: string }[];
  siteNotification: string;
  siteTips: string | null;
  recommendTags: string[];
  scoreConfig: ScoreConfig;
  defaultNodeId: number;
}

/** ⚠️ 仅描述后端 envelope.data 内部的原始结构 */
export interface SiteConfigRaw {
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
    siteIndexHtml: string;
    recommendTags: string[] | null;
    scoreConfig: { postTopicScore: number; postCommentScore: number };
    defaultNodeId: number;
  };
}