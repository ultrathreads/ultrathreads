// components/Header.tsx
'use client';

import { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import { useTranslation } from '@/lib/i18n/i18n-client';
import { useAuth } from '@/hooks/use-auth';
import UserMenu from '@/components/features/UserMenu';
import { useTheme } from '@/hooks/useTheme';
import { useSiteConfig } from '@/providers/SiteConfigProvider';

export default function Header() {
  const config = useSiteConfig();
  const { t } = useTranslation();
  const { isLoggedIn, isLoading } = useAuth();
  const { theme, toggleTheme } = useTheme();

  const [isMounted, setIsMounted] = useState(false);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const createMenuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  useEffect(() => {
    if (!isCreateOpen) return;

    const handleClickOutside = (event: MouseEvent) => {
      if (createMenuRef.current && !createMenuRef.current.contains(event.target as Node)) {
        setIsCreateOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isCreateOpen]);

  if (isLoading) {
    return (
      <header className="header">
        <div className="header-left">
          <Link href="/" className="header-title">{config.siteTitle}</Link>
          <div className="search-box">
            <span className="search-icon">🔍</span>
            <input className="search-input" id="searchInput" placeholder={t('search_default_value')} />
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
            <path d="M443.52 318.72h-166.4v136.96h303.36v-83.2H976v262.4H580.48v-83.2H277.36v156.8a80 80 0 0 0 80 80h223.36v-83.2H976v262.4H580.48v-83.2h-223.36a176 176 0 0 1-176-176V318.72H48v-262.4h395.52v262.4z m232.96 552.96h203.52v-70.4h-203.52v70.4z m0-332.8h203.52v-70.4h-203.52v70.4zM144 222.72h203.52v-70.4H144v70.4z" />
          </svg>
          {config.siteTitle}
        </Link>
        <div className="search-box">
          <span className="search-icon">🔍</span>
          <input className="search-input" id="searchInput" placeholder={t('search_default_value')} />
        </div>
      </div>

      <div className="header-actions">
        {isLoggedIn ? (
          <>
            <div className="create-menu-wrapper" ref={createMenuRef}>
              <button
                className="post-btn"
                onClick={() => setIsCreateOpen((prev) => !prev)}
                aria-expanded={isCreateOpen}
                aria-haspopup="true"
              >
                ✍️ {t('create')}
                <span className={`dropdown-arrow ${isCreateOpen ? 'open' : ''}`}>▾</span>
              </button>

              <div className={`user-dropdown ${isCreateOpen ? 'show' : ''}`}>
                <Link href="/create" className="dropdown-item" onClick={() => setIsCreateOpen(false)}>
                  📝 {t('create_post')}
                </Link>
              </div>
            </div>

            <UserMenu />
          </>
        ) : (
          <>
            <Link href="/auth/login" className="login-link">{t('login')}</Link>
            <Link href="/auth/register" className="register-link">{t('register')}</Link>
          </>
        )}

        {isMounted ? (
          <button
            onClick={toggleTheme}
            className="icon-btn theme-toggle-btn"
            aria-label={theme === 'light' ? t('switch_to_dark') : t('switch_to_light')}
          >
            {theme === 'light' ? '🌙' : '☀️'}
          </button>
        ) : (
          <span className="icon-btn theme-toggle-btn" aria-hidden="true">🌙</span>
        )}
      </div>
    </header>
  );
}