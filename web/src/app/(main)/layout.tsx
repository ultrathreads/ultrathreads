// src/app/(main)/layout.tsx
import Header from '@/components/Header';
import Sidebar from '@/components/Sidebar';
import Footer from '@/components/Footer';
import BackToTop from '@/components/BackToTop';
import PreviewPopover from '@/components/PreviewPopover';

export default function MainLayout({ children }: { children: React.ReactNode }) {
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
        </main>
      </div>
    </div>
  );
}