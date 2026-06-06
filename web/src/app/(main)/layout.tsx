// src/app/(main)/layout.tsx
import Header from '@/components/Header';
import Sidebar from '@/components/Sidebar';
import Footer from '@/components/Footer';
import BackToTop from '@/components/BackToTop';
import PreviewPopover from '@/components/PreviewPopover';
import { getSiteConfig } from '@/lib/api';

export default async function MainLayout({ children }: { children: React.ReactNode }) {
  const config = await getSiteConfig();

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