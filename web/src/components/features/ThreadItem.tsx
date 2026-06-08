'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import type { ThreadViewItem } from '@/types/view';
import type { BackState } from './ThreadTree';
import { RelativeTime } from '@/components/RelativeTime';

interface Props {
  item: ThreadViewItem;
  isRoot?: boolean;
  currentPostId?: string | number;
  globalCollapsed?: boolean;
  backState?: BackState;
  onReplyClick?: (postId: string | number, postTitle: string) => void;
}

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
  backState,
  onReplyClick,
}: Props) {
  const [userOverride, setUserOverride] = useState<boolean | null>(null);
  const folded = userOverride ?? globalCollapsed ?? false;

  const hasReplies = item.replies && item.replies.length > 0;
  const isActive = currentPostId !== undefined && String(item.id) === String(currentPostId);

  useEffect(() => {
    setUserOverride(null);
  }, [globalCollapsed]);

  const handleToggleFold = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setUserOverride(!folded);
  };

  // ✅ 处理回复图标点击
  const handleReplyClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    onReplyClick?.(item.id, `${item.title}(${item.author})`);
  };

  return (
    <li className={folded ? 'folded' : ''}>
      <div className={`entry ${isActive ? 'active' : ''}`}>
        {/* 折叠按钮 */}
        {hasReplies ? (
          <span className="fold-expand">
            <a onClick={handleToggleFold} role="button" tabIndex={0}>
              <svg width="12" height="12" viewBox="0 0 12 12">
                <path
                  d={folded ? 'M4 2l4 4-4 4' : 'M2 4l4 4 4-4'}
                  fill="none"
                  stroke="#7f8c8d"
                  strokeWidth="1.5"
                />
              </svg>
            </a>
          </span>
        ) : (
          isRoot && (
            <span className="fold-expand">
              <svg className="fold-thread" width="12" height="12" viewBox="0 0 12 12" fill="#7f8c8d">
                <rect x="2" y="2" width="6" height="6" />
              </svg>
            </span>
          )
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

        <a className="preview-btn" data-title={item.title} title="预览">
          <svg width="14" height="14" viewBox="0 0 16 16" fill="#95a5a6">
            <path d="M8 3C4 3 1 8 1 8s3 5 7 5 7-5 7-5-3-5-7-5zm0 8a3 3 0 110-6 3 3 0 010 6z" />
          </svg>
        </a>

        {/* ✅ 标题区域：hover 时显示回复图标 */}
        <span className="subject-wrapper">
          {/* ✅ 悬浮显示的回复按钮 */}
          {onReplyClick && (
            <button
              className="inline-reply-btn"
              onClick={handleReplyClick}
              title={`回复 ${item.author}`}
              aria-label={`回复 ${item.author}`}
              type="button"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
              </svg>
            </button>
          )}
        </span>
      </div>

      {hasReplies && (
        <ul className={`reply ${folded ? 'collapsed' : ''}`}>
          {item.replies.map((r) => (
            <ThreadItem
              key={r.id}
              item={r}
              currentPostId={currentPostId}
              globalCollapsed={globalCollapsed}
              backState={backState}
              onReplyClick={onReplyClick}
            />
          ))}
        </ul>
      )}
    </li>
  );
}