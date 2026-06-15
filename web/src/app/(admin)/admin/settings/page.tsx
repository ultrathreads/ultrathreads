// src/app/(admin)/admin/settings/page.tsx
import { getSettings, saveSettings } from './actions';
import { SettingsForm } from './SettingsForm';

export default async function AdminSettingsPage() {
  const settings = await getSettings();

  return (
    <div className="admin-dashboard">
      <div className="admin-page-header">
        <div>
          <h1 className="admin-page-title">站点设置</h1>
          <p className="admin-page-desc">管理站点基础信息与推荐标签</p>
        </div>
      </div>

      <div className="admin-card">
        <div className="admin-card-header">
          <h2 className="admin-card-title">基础配置</h2>
        </div>
        {/* ✅ Next.js App Router 允许直接将 Server Action 传入 Client Component */}
        <SettingsForm initialData={settings} saveAction={saveSettings} />
      </div>
    </div>
  );
}