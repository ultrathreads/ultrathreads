// src/app/(auth)/layout.tsx
import Header from '@/components/layout/Header';
import Footer from '@/components/layout/Footer';

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="app-layout">
      <div className="auth-page-wrapper">
        {/* 顶部导航 */}
        <Header />
        
        {/* 中间内容区：占据剩余空间并居中表单 */}
        <div className="auth-content-area">
          {children}
        </div>

      </div>
    </div>
  );
}