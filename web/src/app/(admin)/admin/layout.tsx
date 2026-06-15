// src/app/(admin)/admin/layout.tsx
import { redirect } from 'next/navigation';
import { getServerSession } from '@/lib/auth/session';
import { canAccessAdminPanel } from '@/lib/auth/permissions';
import { AdminSidebar } from './_components/AdminSidebar';
import { AdminHeader } from './_components/AdminHeader';
import { AuthProvider } from '@/providers/AuthProvider';
import { Toaster } from 'sonner';

// ✅ 引入独立的后台样式
import '@/styles/admin.css';

export default async function AdminLayout({ children }: { children: React.ReactNode }) {
  const user = await getServerSession();

  if (!canAccessAdminPanel(user)) {
    redirect('/auth/login?redirect=/admin');
  }

  return (
    // ✅ 移除 <html> 和 <body>，用 div 作为 admin-body 容器
    // 这样既保持了 CSS 隔离，又避免了 SSR/CSR 结构不一致
    <AuthProvider initialUser={user}>
      <div className="admin-body">
        <div className="admin-layout">
          <AdminSidebar />
          <div className="admin-main-wrapper">
            <AdminHeader />
            <main className="admin-content">
              {children}
            </main>
            <Toaster
              position="top-center"
              toastOptions={{
                className: 'admin-toast',       // 复用已有 admin 样式前缀
                duration: 3000,
              }}
            />
          </div>
        </div>
      </div>
    </AuthProvider>
  );
}