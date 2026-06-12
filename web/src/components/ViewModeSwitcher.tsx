// src/components/ViewModeSwitcher.tsx
'use client';

import Link from 'next/link';
import { usePathname, useSearchParams } from 'next/navigation';

export function ViewModeSwitcher() {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const currentView = searchParams.get('view') === 'flat' ? 'flat' : 'tree';

  // ✅ 直接从路径中提取 ID，硬拼 #post-id，不再依赖 window.location.hash
  // 假设路径格式为 /threads/123 或 /threads/123/xxx
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

    // 只有 flat 带 hash，tree 不带
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
        🌳 树形
      </Link>
      <Link
        href={createHref('flat')}
        className={`mode-btn ${currentView === 'flat' ? 'active' : ''}`}
        title="平铺模式"
        aria-current={currentView === 'flat' ? 'true' : undefined}
      >
        📋 平铺
      </Link>
    </div>
  );
}