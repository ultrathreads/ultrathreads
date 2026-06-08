'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import clsx from 'clsx'; // 推荐使用 clsx 或 tailwind-merge 处理类名拼接

const NAV_ITEMS = [
  { href: '/settings/profile', label: '👤 个人信息' },
  { href: '/settings/account', label: '🔒 账号设置' },
];

export function SettingsNav() {
  const pathname = usePathname();

  return (
    <nav className="settings-nav">
      {NAV_ITEMS.map((item) => {
        // 精确匹配当前路径，避免 /settings/profile/edit 错误高亮 /settings/profile
        const isActive = pathname === item.href;

        return (
          <Link
            key={item.href}
            href={item.href}
            className={clsx('settings-link', {
              'settings-link active': isActive, // 替换为你的高亮样式类名
            })}
            aria-current={isActive ? 'page' : undefined}
          >
            {item.label}
          </Link>
        );
      })}
    </nav>
  );
}