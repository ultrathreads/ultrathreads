// src/providers/SiteConfigProvider.tsx
'use client';

import { createContext, useContext } from 'react';
import type { SiteConfig } from '@/lib/types/site';
import { FALLBACK_SITE_CONFIG } from '@/services/site-service';

const SiteConfigContext = createContext<SiteConfig>(FALLBACK_SITE_CONFIG);

/**
 * 全局站点配置 Hook
 * 在任何 'use client' 组件中直接使用：const { recommendTags } = useSiteConfig();
 */
export function useSiteConfig() {
  return useContext(SiteConfigContext);
}

export function SiteConfigProvider({ 
  config, 
  children 
}: { 
  config: SiteConfig; 
  children: React.ReactNode; 
}) {
  return (
    <SiteConfigContext.Provider value={config}>
      {children}
    </SiteConfigContext.Provider>
  );
}