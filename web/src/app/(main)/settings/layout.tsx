import { SettingsNav } from './settings-nav';

export default function SettingsLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="profile-container">
      <main className="profile-content">
        {children}
      </main>
      <aside className="profile-sidebar">
        <h3 className="profile-title">个人设置</h3>
        {/* 使用客户端导航组件 */}
        <SettingsNav />
      </aside>
    </div>
  );
}