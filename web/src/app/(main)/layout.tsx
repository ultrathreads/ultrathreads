// src/app/(main)/layout.tsx
import { Toaster } from 'sonner';
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import Footer from '@/components/layout/Footer';
import BackToTop from '@/components/BackToTop';
import PreviewPopover from '@/components/ui/PreviewPopover';
import { fetchSiteConfig } from '@/services/site-service';

export default async function MainLayout({ children }: { children: React.ReactNode }) {

  return (
    <div className="app-layout">
      <Header />
      <div className="content-area">
        <Sidebar />
        <main className="main-content" id="mainContent">
          <div className="main-body">{children}</div>
          <Footer />
          <BackToTop />
          <PreviewPopover />
          <Toaster position="top-center" richColors />
        </main>
      </div>
    </div>
  );
}