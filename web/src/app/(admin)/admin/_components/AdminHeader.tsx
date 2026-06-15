// src/components/admin/AdminHeader.tsx
'use client';

// ✅ 直接复用前台已有的 hook，无需改动 AuthProvider
import { useAuth } from '@/hooks/use-auth';
import Link from 'next/link';

export function AdminHeader() {
  // 🎉 可以直接拿到 displayName 和 avatarUrl，连 user?.nickname 的判断都省了
  const { displayName, avatarUrl, logout } = useAuth();

  return (
    <header className="admin-header">
      <div className="admin-header-left">
        <nav className="admin-breadcrumb">
          <Link href="/admin">后台</Link>
          <span className="admin-breadcrumb-separator">/</span>
          <span className="admin-breadcrumb-current">仪表盘</span>
        </nav>
      </div>

      <div className="admin-header-right">
        <Link href="/" className="admin-action-link" target="_blank" rel="noopener noreferrer">
          访问前台
        </Link>

        <div className="admin-user-info">
          {avatarUrl ? (
            <img src={avatarUrl} alt="" className="admin-user-avatar" />
          ) : (
            <div className="admin-user-avatar" style={{ background: 'var(--bg-muted)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '0.75rem' }}>
              {displayName[0] || 'A'}
            </div>
          )}
          <span className="admin-user-name">{displayName || 'Admin'}</span>
        </div>
      </div>
    </header>
  );
}