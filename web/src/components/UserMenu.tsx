'use client';
import { useState, useRef, useEffect } from 'react';
import Link from 'next/link';
import { useAuth } from '@/components/providers/AuthProvider';

export default function UserMenu() {
  const { user, isLoading, error, logout } = useAuth();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  // ✅ Click Outside 逻辑保持不变
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (ref.current && !ref.current.contains(event.target as Node)) {
        setOpen(false);
      }
    }
    if (open) {
      document.addEventListener('mousedown', handleClickOutside);
    }
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [open]);

  // ✅ 删除了原有的 isLoading / error / !user 调试占位符
  // 在数据未就绪或未登录时，不渲染任何内容
  if (isLoading || error || !user) {
    return null;
  }

  return (
    <div className="user-menu-wrapper" ref={ref} onClick={() => setOpen(prev => !prev)}>
      <img
        className="user-avatar"
        src={user.avatar || `https://api.dicebear.com/7.x/avataaars/svg?seed=${user.username}`}
        alt={user.nickname}
      />
      
      <span className="user-name">{user.nickname || user.username}</span>
      <span className={`user-arrow ${open ? 'active' : ''}`}>▼</span>
      
      <div className={`user-dropdown ${open ? 'show' : ''}`} id="userDropdown">
        <div className="dropdown-header">
          <span className="user-level">{user.levelName}</span>
          <span className="user-score">积分: {user.score}</span>
        </div>
        <div className="dropdown-divider" />

        <Link href="/settings/profile" className="dropdown-item" onClick={(e) => { e.stopPropagation(); setOpen(false); }}>
          👤 个人中心
        </Link>
        
        <Link href="/settings/account" className="dropdown-item" onClick={(e) => { e.stopPropagation(); setOpen(false); }}>
          ⚙️ 账号设置
        </Link>
        
        <div className="dropdown-divider" />
        
        <div 
          className="dropdown-item danger" 
          onClick={(e) => {
            e.stopPropagation();
            setOpen(false);
            logout();
          }}
        >
          🚪 退出登录
        </div>
      </div>
    </div>
  );
}