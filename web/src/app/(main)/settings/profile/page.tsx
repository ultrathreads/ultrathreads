// src/app/(main)/settings/profile/page.tsx
export default function ProfilePage() {
  return (
    <div className="settings-card">
      <h2 className="settings-card-title">个人信息</h2>
      <p className="settings-card-desc">在这里管理你的公开资料，让其他用户更好地了解你。</p>
      
      <form className="settings-form">
        <div className="form-group">
          <label className="form-label">头像</label>
          <div className="avatar-upload">
            <div className="avatar-preview">👤</div>
            <button type="button" className="upload-btn">更换头像</button>
          </div>
        </div>
        <div className="form-group">
          <label className="form-label">昵称</label>
          <input type="text" className="form-input" defaultValue="UltraThreads 用户" />
        </div>
        <div className="form-group">
          <label className="form-label">个人简介</label>
          <textarea className="form-textarea" rows={4} placeholder="介绍一下你自己吧..." />
        </div>
        <button type="submit" className="save-btn">保存修改</button>
      </form>
    </div>
  );
}