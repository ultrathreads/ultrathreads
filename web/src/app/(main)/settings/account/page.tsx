// src/app/(main)/settings/account/page.tsx
export default function AccountPage() {
  return (
    <div className="settings-card">
      <h2 className="settings-card-title">账号设置</h2>
      <p className="settings-card-desc">管理你的邮箱、密码以及账号安全选项。</p>
      
      <form className="settings-form">
        <div className="form-group">
          <label className="form-label">绑定邮箱</label>
          <input type="email" className="form-input" defaultValue="user@example.com" />
        </div>
        <div className="form-group">
          <label className="form-label">修改密码</label>
          <input type="password" className="form-input" placeholder="输入新密码" />
        </div>
        <div className="form-group">
          <label className="form-label">确认密码</label>
          <input type="password" className="form-input" placeholder="再次输入新密码" />
        </div>
        <button type="submit" className="save-btn">更新账号信息</button>
      </form>
    </div>
  );
}