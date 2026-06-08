'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import clsx from 'clsx';

const NAV_ITEMS = [
  { href: '/settings/profile', label: '👤 个人信息' },
  { href: '/settings/account', label: '🔒 账号设置' },
];

export function SettingsNav() {
  const pathname = usePathname();

  return (
    <nav className="settings-nav">
      {NAV_ITEMS.map((item) => {
        const isActive = pathname === item.href;

        return (
          <Link
            key={item.href}
            href={item.href}
            className={clsx('settings-link', {
              'settings-link active': isActive,
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