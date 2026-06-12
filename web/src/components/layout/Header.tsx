'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useSearchParams, useParams } from 'next/navigation';
import { useTranslation } from '@/lib/i18n/i18n-client';
import { useAuth } from '@/hooks/use-auth';
import UserMenu from '@/components/features/UserMenu';
import { useTheme } from '@/hooks/useTheme';

interface HeaderProps {
  siteTitle: string;
}

export default function Header({ siteTitle }: HeaderProps) {
  const { t } = useTranslation(['common']);
  const { isLoggedIn, isLoading } = useAuth();
  const params = useParams<{ slug?: string }>(); 
  const { theme, toggleTheme } = useTheme();

  // 2. 添加挂载状态，用于解决 SSR 水合不匹配问题
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  const nodeSlug = params.slug;
  const createHref = nodeSlug 
    ? `/create?nodeSlug=${encodeURIComponent(nodeSlug)}` 
    : '/create';

  // 3. 如果正在加载，显示骨架屏（保持你原有的逻辑）
  if (isLoading) {
    return (
      <header className="header">
        <div className="header-left">
          <Link href="/" className="header-title">{siteTitle}</Link>
          <div className="search-box">
            <span className="search-icon">🔍</span>
            <input className="search-input" id="searchInput" placeholder={t('common:search_default_value')} />
          </div>
        </div>
        <div className="header-actions">
          <span className="loading-placeholder">...</span>
        </div>
      </header>
    );
  }

  return (
    <header className="header">
      <div className="header-left">
        <Link href="/" className="header-title">
          <svg className="ut-logo" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" aria-hidden="true">
            <path d="M443.52 318.72h-166.4v136.96h303.36v-83.2H976v262.4H580.48v-83.2H277.36v156.8a80 80 0 0 0 80 80h223.36v-83.2H976v262.4H580.48v-83.2h-223.36a176 176 0 0 1-176-176V318.72H48v-262.4h395.52v262.4z m232.96 552.96h203.52v-70.4h-203.52v70.4z m0-332.8h203.52v-70.4h-203.52v70.4zM144 222.72h203.52v-70.4H144v70.4z"></path>
          </svg>
          {siteTitle}
        </Link>
        <div className="search-box">
          <span className="search-icon">🔍</span>
          <input className="search-input" id="searchInput" placeholder={t('common:search_default_value')} />
        </div>
      </div>
      <div className="header-actions">
        {isLoggedIn ? (
          <>
            <Link href={createHref} className="post-btn">✏️ {t('common:posting')}</Link>
            <UserMenu />
          </>
        ) : (
          <>
            <Link href="/auth/login" className="login-link">{t('common:login')}</Link>
            <Link href="/auth/register" className="register-link">{t('common:register')}</Link>
          </>
        )}
        
        {isMounted ? (
          <button
            onClick={toggleTheme}
            className="icon-btn theme-toggle-btn"
            aria-label={theme === 'light' ? '切换到深色模式' : '切换到亮色模式'}
          >
            {theme === 'light' ? '🌙' : '☀️'}
          </button>
        ) : (
          // 提供一个占位符，防止按钮出现时导致布局抖动（Layout Shift）
          <span className="icon-btn theme-toggle-btn" aria-hidden="true">
            🌙
          </span>
        )}
      </div>
    </header>
  );
}