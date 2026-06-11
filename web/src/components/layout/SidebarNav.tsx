'use client';

import { useRouter } from 'next/navigation';
import { usePathname, useSearchParams } from 'next/navigation';
import { useTranslation } from '@/lib/i18n/i18n-client';
import clsx from 'clsx';

const NAV_ITEMS = [
  { href: '/', labelKey: 'common:home', matchType: 'exact' as const },
  { href: '/my/posts', labelKey: 'common:mine', matchType: 'exact' as const },
];

export function SidebarNav() {
  const { t } = useTranslation(['common']);
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const router = useRouter();

  return (
    <ul className="forum-list" role="list">
      {NAV_ITEMS.map((item) => {
        const isActive =
          item.matchType === 'exact'
            ? pathname === item.href && !searchParams.has('nodeId')
            : pathname.startsWith('/settings');

        return (
          <li
            key={item.labelKey}
            role="link"
            tabIndex={0}
            aria-current={isActive ? 'page' : undefined}
            className={clsx('forum-item cursor-pointer', { active: isActive })}
            onClick={() => router.push(item.href)}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                router.push(item.href);
              }
            }}
          >
            {t(item.labelKey)}
          </li>
        );
      })}
    </ul>
  );
}