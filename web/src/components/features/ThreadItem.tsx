// components/features/ThreadItem.tsx
'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import type { ThreadViewItem } from '@/types/view';
import type { BackState } from './ThreadTree'; // ✅ 1. 从 ThreadTree 导入共享类型
import { RelativeTime } from '@/components/RelativeTime';

interface Props {
  item: ThreadViewItem;
  isRoot?: boolean;
  currentPostId?: string | number;
  globalCollapsed?: boolean;
  backState?: BackState; // ✅ 2. 新增可选的回溯状态属性
}

/**
 * 构建带回溯参数的详情页链接
 * 与 ThreadTree 中的逻辑保持一致，确保所有入口生成的 URL 格式统一
 */
function buildPostUrl(postId: number | string, backState?: BackState): string {
  if (!backState || (!backState.nodeId && !backState.page)) {
    return `/post/${postId}`;
  }

  const params = new URLSearchParams();
  if (backState.nodeId) params.set('nodeId', backState.nodeId);
  if (backState.page) params.set('page', backState.page);

  return `/post/${postId}?${params.toString()}`;
}

export default function ThreadItem({
  item,
  isRoot,
  currentPostId,
  globalCollapsed,
  backState, // ✅ 3. 解构接收 backState
}: Props) {
  const [localFolded, setLocalFolded] = useState(false);

  const folded = globalCollapsed ?? localFolded;

  // ✅ 当全局折叠状态变化时，同步重置本地状态
  // 防止用户手动展开后，再次点击"全部折叠"时因本地状态覆盖而失效
  useEffect(() => {
    if (globalCollapsed !== undefined) {
      setLocalFolded(globalCollapsed);
    }
  }, [globalCollapsed]);

  const hasReplies = item.replies && item.replies.length > 0;
  const isActive = currentPostId !== undefined && String(item.id) === String(currentPostId);

  return (
    <li className={folded ? 'folded' : ''}>
      <div className={`entry ${isActive ? 'active' : ''}`}>
        {isRoot && (
          <span className="fold-expand">
            {hasReplies ? (
              <a onClick={(e) => { e.preventDefault(); setLocalFolded((prev) => !prev); }}>
                <svg width="12" height="12" viewBox="0 0 12 12">
                  <path d="M2 4l4 4 4-4" fill="none" stroke="#7f8c8d" strokeWidth="1.5" />
                </svg>
              </a>
            ) : (
              <svg className="fold-thread" width="12" height="12" viewBox="0 0 12 12" fill="#7f8c8d">
                <rect x="2" y="2" width="6" height="6" />
              </svg>
            )}
          </span>
        )}
        {isRoot ? (
          <svg className="icon-topic" width="14" height="14" viewBox="0 0 14 14" fill="#3498db">
            <circle cx="7" cy="7" r="5"></circle>
          </svg>
        ) : (
          <svg className="icon-reply-svg" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 264 264">
            <path d="M6,16 v108 a57 57, 0, 0, 0, 57, 57 h92 v27 l 45.5,-45.5 -45.5,-45.5 v27 h-92 a20 20, 0, 0, 1, -20,-20 v-108 z"></path>
          </svg>
        )}

        {/* ✅ 4. 核心改动：原生 <a> 替换为 Next.js <Link>，并注入回溯参数 */}
        <Link
          className={`subject ${isRoot ? '' : 'read'} ${isActive ? 'active' : ''}`}
          href={buildPostUrl(item.id, backState)}
        >
          {item.title}
        </Link>

        <span className="metadata">
          <span className="author-name">{item.author}</span>
          <span className="tail">
            <RelativeTime timestamp={item.date} />
          </span>
          {isRoot && item.nodeName && <span className="category">({item.nodeName})</span>}
        </span>
        <a
          className="preview-btn"
          data-title={item.title}
          title="预览"
        >
          <svg width="14" height="14" viewBox="0 0 16 16" fill="#95a5a6">
            <path d="M8 3C4 3 1 8 1 8s3 5 7 5 7-5 7-5-3-5-7-5zm0 8a3 3 0 110-6 3 3 0 010 6z" />
          </svg>
        </a>
      </div>

      {hasReplies && (
        <ul className={`reply ${folded ? 'collapsed' : ''}`}>
          {item.replies.map((r) => (
            /* ✅ 5. 递归渲染子回复时，继续向下透传 backState */
            <ThreadItem
              key={r.id}
              item={r}
              currentPostId={currentPostId}
              globalCollapsed={globalCollapsed}
              backState={backState}
            />
          ))}
        </ul>
      )}
    </li>
  );
}