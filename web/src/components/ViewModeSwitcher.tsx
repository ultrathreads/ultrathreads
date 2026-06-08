// src/components/ViewModeSwitcher.tsx
'use client';

import Link from 'next/link';
import { usePathname, useSearchParams } from 'next/navigation';

export function ViewModeSwitcher() {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  
  // ✅ 从 URL 读取当前模式，保持与服务端 page.tsx 的判断逻辑完全一致
  const currentView = searchParams.get('view') === 'flat' ? 'flat' : 'tree';

  const createHref = (targetView: 'tree' | 'flat') => {
    const params = new URLSearchParams(searchParams.toString());
    
    if (targetView === 'flat') {
      params.set('view', 'flat');
    } else {
      params.delete('view'); // tree 是默认值，删除参数保持 URL 整洁
    }
    
    const qs = params.toString();
    return qs ? `${pathname}?${qs}` : pathname;
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