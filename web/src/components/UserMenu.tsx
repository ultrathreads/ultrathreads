'use client';
import { useState, useRef, useEffect } from 'react';
import Link from 'next/link'; // 1. 引入 Next.js Link 组件

export default function UserMenu() {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener('click', handler);
    return () => document.removeEventListener('click', handler);
  }, []);

  // 2. 封装一个阻止冒泡的点击处理函数，防止点击菜单项时触发外层 wrapper 的 open/close 切换
  const handleItemClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    setOpen(false); // 点击后主动关闭菜单
  };

  return (
    <div className="user-menu-wrapper" ref={ref} onClick={() => setOpen(!open)}>
      <img
        className="user-avatar"
        src="https://api.dicebear.com/7.x/avataaars/svg?seed=Felix"
        alt="avatar"
      />
      <span className="user-name">张三</span>
      <span className={`user-arrow ${open ? 'active' : ''}`}>▼</span>
      
      <div className={`user-dropdown ${open ? 'show' : ''}`} id="userDropdown">
        {/* 3. 将 div 替换为 Link，并绑定 handleItemClick */}
        <Link 
          href="/settings/profile" 
          className="dropdown-item" 
          onClick={handleItemClick}
        >
          👤 个人中心
        </Link>
        
        <Link 
          href="/settings/account" 
          className="dropdown-item" 
          onClick={handleItemClick}
        >
          ⚙️ 账号设置
        </Link>
        
        <div className="dropdown-divider" />
        
        {/* 退出登录保持 div，后续可接入真实的登出逻辑 */}
        <div 
          className="dropdown-item danger" 
          onClick={(e) => {
            e.stopPropagation();
            setOpen(false);
            // TODO: 在此处添加退出登录的逻辑
            console.log('退出登录');
          }}
        >
          🚪 退出登录
        </div>
      </div>
    </div>
  );
}