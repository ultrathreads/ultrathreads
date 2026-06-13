// src/lib/types/site.ts

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
    recommendTags: string[];
  };
}

export interface SiteConfig {
  appName: string;
  appVersion: string;
  adminLevel: number;
  siteTitle: string;
  siteDescription: string;
  siteKeywords: string | null;
  navLinks: Array<{ label: string; href: string }>;
  defaultNodeId: number;
  recommendTags: string[];
}