// src/components/Footer.tsx
import Link from 'next/link';

interface FooterProps {
  appVersion: string;
}

export default function Footer({ appVersion }: FooterProps) {
  return (
    <footer className="main-footer">
      <div className="footer-container">
        {/* 左侧：版权信息 */}
        <div className="footer-left">
          <p className="footer-copyright">
            &copy; {new Date().getFullYear()} UltraThreads v{appVersion}. All rights reserved.
          </p>
          <p className="footer-slogan">连接每一个有趣的灵魂</p>
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