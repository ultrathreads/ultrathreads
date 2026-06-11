// components/ui/MyPostsTabs.tsx
'use client';

import Link from 'next/link';
import { usePathname, useSearchParams } from 'next/navigation';

const TABS = [
  { key: 'root', label: '📝 我的主帖' },
  { key: 'replies', label: '💬 我的回帖' },
  { key: 'bookmarks', label: '📑 我的书签' },
];

export default function MyPostsTabs() {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const currentTab = searchParams.get('tab') || 'root';

  return (
    <div className="view-mode-switcher" style={{ marginBottom: '20px' }}>
      {TABS.map((tab) => (
        <Link
          key={tab.key}
          href={`${pathname}?tab=${tab.key}`}
          className={`mode-btn ${currentTab === tab.key ? 'active' : ''}`}
        >
          {tab.label}
        </Link>
      ))}
    </div>
  );
}