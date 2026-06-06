'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useTranslation } from '@/lib/i18n-client';
import UserMenu from './UserMenu';

export default function Header() {
  const { t } = useTranslation(['common']);
  // 2. 定义一个状态来存储用户的登录状态
  const [isLoggedIn, setIsLoggedIn] = useState<boolean | null>(null);

  // 3. 在组件挂载后，随机决定登录状态，模拟真实情况
  useEffect(() => {
    // Math.random() 会生成一个 0 到 1 之间的随机数
    // 如果随机数大于 0.5，就视为已登录，否则为未登录
    const randomStatus = Math.random() > 0.5;
    setIsLoggedIn(randomStatus);
  }, []);

  // 4. 在状态确定之前，可以先渲染一个占位符，避免页面闪烁
  if (isLoggedIn === null) {
    return <header className="header">加载中...</header>;
  }

  return (
    <header className="header">
      <div className="header-left">
        <Link href="/" className="header-title">
          <svg className="ut-logo" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" aria-hidden="true">
            <path d="M443.52 318.72h-166.4v136.96h303.36v-83.2H976v262.4H580.48v-83.2H277.12v156.8a80 80 0 0 0 80 80h223.36v-83.2H976v262.4H580.48v-83.2h-223.36a176 176 0 0 1-176-176V318.72H48v-262.4h395.52v262.4z m232.96 552.96h203.52v-70.4h-203.52v70.4z m0-332.8h203.52v-70.4h-203.52v70.4zM144 222.72h203.52v-70.4H144v70.4z"></path>
          </svg>
          UltraThreads
        </Link>
        <div className="search-box">
          <span className="search-icon">🔍</span>
          <input className="search-input" id="searchInput" placeholder={t('common:search_default_value')} />
        </div>
      </div>
      <div className="header-actions">
        {/* 5. 根据登录状态，条件渲染不同的内容 */}
        {isLoggedIn ? (
          // --- 已登录状态 ---
          <>
            <Link href="/create" className="post-btn">✏️ {t('common:posting')}</Link>
            <UserMenu />
          </>
        ) : (
          // --- 未登录状态 ---
          <>
            <Link href="/auth/login" className="login-link">{t('common:login')}</Link>
            <Link href="/auth/register" className="register-link">{t('common:register')}</Link>
          </>
        )}
      </div>
    </header>
  );
}