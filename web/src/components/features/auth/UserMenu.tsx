'use client';

import Link from 'next/link';
import { useAuth } from '@/hooks/use-auth';
import { useClickOutside } from '@/hooks/use-click-outside';
import Avatar from '@/components/ui/Avatar';

export default function UserMenu() {
  const { user, isLoggedIn, isLoading, error, logout, displayName, avatarUrl } = useAuth();
  
  const { isOpen, toggle, close, ref } = useClickOutside(false);

  // 数据未就绪或未登录时，不渲染任何内容
  if (isLoading || error || !isLoggedIn || !user) {
    return null;
  }

  return (

    <div className="user-menu-wrapper" ref={ref}>
      
      <div className="user-menu-trigger" onClick={toggle} style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
        <Avatar 
          className="user-avatar" 
          src={avatarUrl} 
          alt={displayName} 
        />
        <span className="user-name">{displayName}</span>
        <span className={`user-arrow ${isOpen ? 'active' : ''}`}>▼</span>
      </div>

      <div className={`user-dropdown ${isOpen ? 'show' : ''}`} id="userDropdown">
        <div className="dropdown-header">
          <span className="user-level">{user.levelName}</span>
          <span className="user-score">积分: {user.score}</span>
        </div>
        <div className="dropdown-divider" />

        <Link href="/settings/profile" className="dropdown-item" onClick={close}>
          👤 个人中心
        </Link>

        <Link href="/settings/account" className="dropdown-item" onClick={close}>
          ⚙️ 账号设置
        </Link>

        <div className="dropdown-divider" />

        <div
          className="dropdown-item danger"
          onClick={() => {
            close();
            logout();
          }}
        >
          🚪 退出登录
        </div>
      </div>
    </div>
  );
}