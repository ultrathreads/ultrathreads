import { SettingsNav } from './settings-nav';

export default function SettingsLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="settings-container">
      <aside className="settings-sidebar">
        <h3 className="settings-title">个人设置</h3>
        {/* 使用客户端导航组件 */}
        <SettingsNav />
      </aside>

      <main className="settings-content">
        {children}
      </main>
    </div>
  );
}