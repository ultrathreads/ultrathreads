// src/components/ViewModeSwitcher.tsx
'use client';

import Link from 'next/link';
import { usePathname, useSearchParams } from 'next/navigation';

export function ViewModeSwitcher() {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const currentView = searchParams.get('view') === 'flat' ? 'flat' : 'tree';

  const postId = pathname.split('/').filter(Boolean)[1];
  const hash = postId ? `#post-${postId}` : '';

  const createHref = (targetView: 'tree' | 'flat') => {
    const params = new URLSearchParams(searchParams.toString());

    if (targetView === 'flat') {
      params.set('view', 'flat');
    } else {
      params.delete('view');
    }

    const qs = params.toString();
    const base = qs ? `${pathname}?${qs}` : pathname;

    return targetView === 'flat' ? `${base}${hash}` : base;
  };

  return (
    <div className="view-mode-switcher" role="group" aria-label="浏览模式切换">
      <Link
        href={createHref('tree')}
        className={`mode-btn ${currentView === 'tree' ? 'active' : ''}`}
        title="树形模式"
        aria-current={currentView === 'tree' ? 'true' : undefined}
      >
        {/* 树形图标：层级分支结构 */}
        <svg
          className="icon-tree"
          xmlns="http://www.w3.org/2000/svg"
          width="14"
          height="14"
          viewBox="0 0 14 14"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.5"
          strokeLinecap="round"
          strokeLinejoin="round"
          aria-hidden="true"
        >
          <circle cx="3" cy="3" r="1.5" fill="currentColor" stroke="none" />
          <path d="M3 4.5V7h4" />
          <circle cx="9" cy="7" r="1.5" fill="currentColor" stroke="none" />
          <path d="M3 7v4h4" />
          <circle cx="9" cy="11" r="1.5" fill="currentColor" stroke="none" />
        </svg>
        <span>树形</span>
      </Link>
      <Link
        href={createHref('flat')}
        className={`mode-btn ${currentView === 'flat' ? 'active' : ''}`}
        title="平铺模式"
        aria-current={currentView === 'flat' ? 'true' : undefined}
      >
        {/* 平铺图标：对齐列表结构 */}
        <svg
          className="icon-flat"
          xmlns="http://www.w3.org/2000/svg"
          width="14"
          height="14"
          viewBox="0 0 14 14"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.5"
          strokeLinecap="round"
          strokeLinejoin="round"
          aria-hidden="true"
        >
          <line x1="5" y1="3" x2="12" y2="3" />
          <line x1="5" y1="7" x2="12" y2="7" />
          <line x1="5" y1="11" x2="12" y2="11" />
          <circle cx="2.5" cy="3" r="1" fill="currentColor" stroke="none" />
          <circle cx="2.5" cy="7" r="1" fill="currentColor" stroke="none" />
          <circle cx="2.5" cy="11" r="1" fill="currentColor" stroke="none" />
        </svg>
        <span>平铺</span>
      </Link>
    </div>
  );
}