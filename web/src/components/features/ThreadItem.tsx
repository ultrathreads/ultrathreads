'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import type { ThreadViewItem } from '@/types/view';
import type { BackState } from './ThreadTree';
import { RelativeTime } from '@/components/RelativeTime';
import AuthorLink from '@/components/ui/AuthorLink';

interface Props {
  item: ThreadViewItem;
  isRoot?: boolean;
  currentPostSlug?: string;
  globalCollapsed?: boolean;
  backState?: BackState;
}

function buildPostUrl(postSlug: string | string, backState?: BackState): string {
  if (!backState || (!backState.nodeSlug && !backState.tagSlug && !backState.page)) {
    return `/post/${postSlug}`;
  }
  const params = new URLSearchParams();
  if (backState.nodeSlug) params.set('nodeSlug', backState.nodeSlug);
  if (backState.tagSlug) params.set('tagSlug', backState.tagSlug);
  if (backState.page) params.set('page', backState.page);
  return `/post/${postSlug}?${params.toString()}`;
}

export default function ThreadItem({
  item,
  isRoot,
  currentPostSlug,
  globalCollapsed,
  backState,
}: Props) {
  const [userOverride, setUserOverride] = useState<boolean | null>(null);
  const folded = userOverride ?? globalCollapsed ?? false;

  const hasReplies = item.replies && item.replies.length > 0;
  const isActive = currentPostSlug !== undefined && String(item.slug) === String(currentPostSlug);

  useEffect(() => {
    setUserOverride(null);
  }, [globalCollapsed]);

  const handleToggleFold = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setUserOverride(!folded);
  };

  return (
    <li className={folded ? 'folded' : ''}>
      <div className={`entry ${isActive ? 'active' : ''}`}>
        {/* 折叠按钮 */}
        {isRoot ? (
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
          item.isPinned ? (
            <svg className="icon-pinned" width="14" height="14" viewBox="0 0 24 24" fill="#e74c3c">
              <path d="M16 2H8a1 1 0 0 0-1 1v3.27l-3.88 3.88a1 1 0 0 0-.29.7V12a1 1 0 0 0 1 1h7v5l-2 2v1h6v-1l-2-2v-5h7a1 1 0 0 0 1-1v-1.15a1 1 0 0 0-.29-.7L17 5.27V3a1 1 0 0 0-1-1zM9 4h6v1.5l3.5 3.5H5.5L9 5.5V4z" />
            </svg>
          ) : (
            <svg 
              className={`icon-topic ${item.isRead ? 'is-read' : 'is-unread'}`} 
              width="14" 
              height="14" 
              viewBox="0 0 14 14" 
              fill="currentColor"
            >
              <circle cx="7" cy="7" r="5" />
            </svg>
          )
        ) : (
          <svg className="icon-reply-svg" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 264 264" width="14" height="14">
            <path d="M6,16 v108 a57 57,0 0 0 57,57 h92 v27 l45.5-45.5-45.5-45.5 v27 h-92 a20 20,0 0 1-20-20 v-108z" />
          </svg>
        )}
        <Link
          className={`subject ${isRoot ? '' : 'read'} ${isActive ? 'active' : ''}`}
          href={buildPostUrl(item.slug, backState)}
        >
          {item.title} 
        </Link>

        <span className="metadata">
          <AuthorLink 
            author={item.author} 
            authorSlug={item.authorSlug} 
            className="author-name" 
          />
          <span className="tail">
            <RelativeTime timestamp={item.date} />
          </span>
          {isRoot && item.nodeName && <span className="category">({item.nodeName})</span>}
        </span>
        <button
          className="icon-btn preview-btn"
          data-post-slug={String(item.slug)}
          title={`回复 ${item.author}`}
          aria-label={`回复 ${item.author}`}
          type="button"
        >
          <svg width="14" height="14" viewBox="0 0 16 16" fill="#95a5a6" aria-hidden="true">
            <path d="M8 3C4 3 1 8 1 8s3 5 7 5 7-5 7-5-3-5-7-5zm0 8a3 3 0 110-6 3 3 0 010 6z" />
          </svg>
        </button>
        {isRoot && (
          <Link
            className="icon-btn flat-view-btn"
            href={`/post/${item.slug}?view=flat`}
            title={`平铺模式浏览 ${item.author} 的帖子`}
            aria-label={`平铺模式浏览 ${item.author} 的帖子`}
          >
            <svg width="14" height="14" viewBox="0 0 16 16" fill="#95a5a6">
              <path d="M2.5 3.5A.5.5 0 0 1 3 3h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5zm0 4A.5.5 0 0 1 3 7h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5zm0 4a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5z" />
            </svg>
          </Link>
        )}
      </div>

      {hasReplies && (
        <ul className={`reply ${folded ? 'collapsed' : ''}`}>
          {item.replies.map((r) => (
            <ThreadItem
              key={r.slug}
              item={r}
              currentPostSlug={currentPostSlug}
              globalCollapsed={globalCollapsed}
              backState={backState}
            />
          ))}
        </ul>
      )}
    </li>
  );
}