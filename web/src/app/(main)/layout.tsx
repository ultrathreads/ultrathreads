// src/app/(main)/layout.tsx
import { Toaster } from 'sonner';
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import Footer from '@/components/layout/Footer';
import BackToTop from '@/components/BackToTop';
import PreviewPopover from '@/components/ui/PreviewPopover';
import { fetchSiteConfig } from '@/services/site-service';
import { SiteConfigProvider } from '@/providers/SiteConfigProvider'; 

export default async function MainLayout({ children }: { children: React.ReactNode }) {
  // 全站唯一调用点（unstable_cache 保证无重复请求）
  const config = await fetchSiteConfig();

  return (
    <SiteConfigProvider config={config}>
      <div className="app-layout">
        <Header 
          siteTitle={config.siteTitle} 
        />
        <div className="content-area">
          <Sidebar />
          <main className="main-content" id="mainContent">
            <div className="main-body">{children}</div>
            <Footer 
              appVersion={config.appVersion}
            />
            <BackToTop />
            <PreviewPopover />
            <Toaster position="top-center" richColors />
          </main>
        </div>
      </div>
    </SiteConfigProvider>
  );
}