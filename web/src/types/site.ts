// src/lib/types/site.ts

/** 后端原始返回结构（严格对应最新 API Response） */
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
    siteNavs: Array<{ label: string; href: string }> | null;
    defaultNodeId: number;
  };
}

/** 前端标准消费结构（保证非空、字段命名统一） */
export interface SiteConfig {
  appName: string;
  appVersion: string;
  adminLevel: number;
  siteTitle: string;
  siteDescription: string;
  siteKeywords: string | null;
  navLinks: Array<{ label: string; href: string }>; // ✅ 前端保证为数组
  defaultNodeId: number;
}