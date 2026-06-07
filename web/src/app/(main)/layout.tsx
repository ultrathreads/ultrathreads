// src/app/(main)/layout.tsx
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import Footer from '@/components/layout/Footer';
import BackToTop from '@/components/BackToTop';
import PreviewPopover from '@/components/ui/PreviewPopover';
import { fetchSiteConfig } from '@/services/site-service';

export default async function MainLayout({ children }: { children: React.ReactNode }) {
  const config = await fetchSiteConfig();

  return (
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
        </main>
      </div>
    </div>
  );
}