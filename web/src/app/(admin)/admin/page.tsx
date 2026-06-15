// src/app/(admin)/admin/page.tsx
import { getSystemInfo } from '@/services/admin.server';

export default async function AdminDashboardPage() {
  const sysInfo = await getSystemInfo();

  const infoItems: { label: string; value: string | number; mono?: boolean }[] = [
    { label: '应用名称', value: sysInfo.appName },
    { label: '版本号', value: sysInfo.appVersion },
    { label: '运行时长', value: sysInfo.upTime },
    { label: '操作系统', value: `${sysInfo.os} / ${sysInfo.arch}` },
    { label: 'CPU 核心数', value: sysInfo.numCpu },
    { label: 'Go 版本', value: sysInfo.goversion },
    { label: '构建时间', value: sysInfo.buildTime, mono: true },
    { label: 'Commit', value: sysInfo.buildCommit.slice(0, 8), mono: true },
  ];

  return (
    <div className="admin-dashboard">
      {/* 页面标题区 */}
      <div className="admin-page-header">
        <div>
          <h1 className="admin-page-title">仪表盘</h1>
          <p className="admin-page-desc">系统概览与运行状态</p>
        </div>
      </div>

      {/* 统计卡片占位区 */}
      <div className="admin-stats-grid">
        <div className="admin-stat-card">
          <span className="admin-stat-label">注册用户</span>
          <span className="admin-stat-value">{sysInfo.registerUserCount}</span>
        </div>
        <div className="admin-stat-card">
          <span className="admin-stat-label">帖子总数</span>
          <span className="admin-stat-value">{sysInfo.postTotalCount}</span>
        </div>
        <div className="admin-stat-card">
          <span className="admin-stat-label">今日新帖</span>
          <span className="admin-stat-value">{sysInfo.todayNewPostCount}</span>
        </div>
        <div className="admin-stat-card">
          <span className="admin-stat-label">系统状态</span>
          <span className="admin-stat-value" style={{ color: '#10b981' }}>正常</span>
        </div>
      </div>

      {/* 系统信息卡片 */}
      <div className="admin-card">
        <div className="admin-card-header">
          <h2 className="admin-card-title">系统信息</h2>
          {/* ✅ 已修正 Badge 类名 */}
          <span className="admin-badge admin-badge-success">Running</span>
        </div>

        <dl className="admin-info-grid">
          {infoItems.map((item) => (
            <div key={item.label} className="admin-info-item">
              <dt className="admin-info-label">{item.label}</dt>
              <dd className={`admin-info-value${item.mono ? ' mono' : ''}`}>
                {item.value}
              </dd>
            </div>
          ))}
        </dl>
      </div>
    </div>
  );
}