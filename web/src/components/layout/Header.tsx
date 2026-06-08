'use client';

import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { useTranslation } from '@/lib/i18n/i18n-client';
import { useAuth } from '@/hooks/use-auth';
import UserMenu from '@/components/features/auth/UserMenu';

interface HeaderProps {
  siteTitle: string;
}

export default function Header({ siteTitle }: HeaderProps) {
  const { t } = useTranslation(['common']);
  const { isLoggedIn, isLoading } = useAuth();
  const searchParams = useSearchParams();

  const nodeId = searchParams.get('nodeId');
  const createHref = nodeId ? `/create?nodeId=${encodeURIComponent(nodeId)}` : '/create';

  // 加载中保持原有布局骨架，避免闪烁
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
            <path d="M443.52 318.72h-166.4v136.96h303.36v-83.2H976v262.4H580.48v-83.2H277.12v156.8a80 80 0 0 0 80 80h223.36v-83.2H976v262.4H580.48v-83.2h-223.36a176 176 0 0 1-176-176V318.72H48v-262.4h395.52v262.4z m232.96 552.96h203.52v-70.4h-203.52v70.4z m0-332.8h203.52v-70.4h-203.52v70.4zM144 222.72h203.52v-70.4H144v70.4z"></path>
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
      </div>
    </header>
  );
}