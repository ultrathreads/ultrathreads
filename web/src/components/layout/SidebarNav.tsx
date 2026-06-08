'use client';

import Link from 'next/link';
import { usePathname, useSearchParams } from 'next/navigation';
import { useTranslation } from '@/lib/i18n/i18n-client';
import clsx from 'clsx';

const NAV_ITEMS = [
  { href: '/', labelKey: 'common:home', matchType: 'exact' as const },
  { href: '/settings/profile', labelKey: 'common:mine', matchType: 'startsWith' as const },
];

export function SidebarNav() {
  const { t } = useTranslation(['common']);
  const pathname = usePathname();
  const searchParams = useSearchParams();

  return (
    <ul className="forum-list">
      {NAV_ITEMS.map((item) => {
        let isActive: boolean;

        if (item.matchType === 'exact') {
          isActive = pathname === item.href && !searchParams.has('nodeId');
        } else {
          isActive = pathname.startsWith('/settings');
        }

        return (
          <li key={item.labelKey}>
            <Link
              href={item.href}
              className={clsx('forum-item', { active: isActive })}
              aria-current={isActive ? 'page' : undefined}
            >
              {t(item.labelKey)}
            </Link>
          </li>
        );
      })}
    </ul>
  );
}