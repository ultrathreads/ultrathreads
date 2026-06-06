// src/app/(main)/settings/layout.tsx
import Link from 'next/link';

export default function SettingsLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="settings-container">
      {/* 左侧：设置导航菜单 */}
      <aside className="settings-sidebar">
        <h3 className="settings-title">个人设置</h3>
        <nav className="settings-nav">
          <Link href="/settings/profile" className="settings-link">
            👤 个人信息
          </Link>
          <Link href="/settings/account" className="settings-link">
            🔒 账号设置
          </Link>
        </nav>
      </aside>

      {/* 右侧：具体的设置页面内容 */}
      <main className="settings-content">
        {children}
      </main>
    </div>
  );
}