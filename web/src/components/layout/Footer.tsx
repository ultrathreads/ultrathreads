// src/components/Footer.tsx
'use client'; 

import Link from 'next/link';
import { useSiteConfig } from '@/providers/SiteConfigProvider';

export default function Footer() {
  const config = useSiteConfig();

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

        {/* 右侧：快捷链接 */}
        <div className="footer-right">
          <nav className="footer-links">
            <Link href="/about">关于我们</Link>
            <Link href="/terms">服务条款</Link>
            <Link href="/privacy">隐私政策</Link>
            <Link href="/contact">联系我们</Link>
          </nav>
        </div>
      </div>
    </footer>
  );
}