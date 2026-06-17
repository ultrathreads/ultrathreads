// src/components/Footer.tsx
'use client';

import Link from 'next/link';
import { useSiteConfig } from '@/providers/SiteConfigProvider';
import { useTranslation } from '@/lib/i18n/i18n-client';
import i18next, { languages, type Locale } from '@/lib/i18n/i18n-client';

export default function Footer() {
  const config = useSiteConfig();
  const { t } = useTranslation();

  // ✅ 语言切换核心逻辑
  const handleLanguageChange = async (lng: Locale) => {
    if (i18next.language === lng) return;
    await i18next.changeLanguage(lng);
    window.location.reload();
  };

  return (
    <footer className="main-footer">
      <div className="footer-container">
        {/* 左侧：版权信息 */}
        <div className="footer-left">
          <p className="footer-copyright">
            &copy; 2026-{new Date().getFullYear()} UltraThreads v{config.appVersion}. All rights reserved.
          </p>
          <p className="footer-slogan">{config.siteTitle} | {config.siteDescription}</p>
        </div>

        {/* 右侧：快捷链接 + 语言切换 */}
        <div className="footer-right">
          <nav className="footer-links">
            <Link href="/about">{t('about_us', '关于我们')}</Link>
            <Link href="/terms">{t('terms_of_service', '服务条款')}</Link>
            <Link href="/privacy">{t('privacy_policy', '隐私政策')}</Link>
            <Link href="/contact">{t('contact_us', '联系我们')}</Link>

            {/* ✅ 分隔符 + 语言切换按钮组 */}
            <span className="footer-divider" aria-hidden="true">|</span>
            <div className="footer-lang-switcher" role="group" aria-label={t('switch_language', '切换语言')}>
              {languages.map((lng) => (
                <button
                  key={lng}
                  onClick={() => handleLanguageChange(lng)}
                  className={`footer-lang-btn ${i18next.language === lng ? 'active' : ''}`}
                  aria-current={i18next.language === lng ? 'true' : undefined}
                >
                  {lng === 'zh' ? '中文' : 'English'}
                </button>
              ))}
            </div>
          </nav>
        </div>
      </div>
    </footer>
  );
}