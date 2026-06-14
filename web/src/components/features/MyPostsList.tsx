// components/features/MyPostsList.tsx
'use client';

import Link from 'next/link';
import type { MyPostListItem } from '@/services/my-post-service';
import { RelativeTime } from '@/components/ui/RelativeTime';
import type { UserEntity } from '@/types/domains';

interface Props {
  initialPosts: MyPostListItem[];
  initialPaging: { page: number; pageSize: number; totalItems: number };
  user?: UserEntity;
  tab: 'root' | 'replies';
}

export default function MyPostsList({ initialPosts, user, tab }: Props) {
  // 根据 tab 类型决定标题
  const pageTitle = user
    ? (tab === 'replies' ? `${user.nickname} 的回帖` : `${user.nickname} 的帖子`)
    : (tab === 'replies' ? '我的回帖' : '我的主帖');

  return (
    <div className="thread-tree-container">
      <div className="thread-tree-header">
        <h1>{pageTitle}</h1>
      </div>

      <ul className="thread">
        {initialPosts.map((item) => {
          return (
            <li key={item.slug}>
              <div className="entry">
                {item.isRoot ? (
                  <svg className="icon-topic" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 14 14" fill="currentColor">
                    <circle cx="7" cy="7" r="5" />
                  </svg>
                ) : (
                  <svg className="icon-reply-svg" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 264 264" width="14" height="14">
                    <path d="M6,16 v108 a57 57,0 0 0 57,57 h92 v27 l45.5-45.5-45.5-45.5 v27 h-92 a20 20,0 0 1-20-20 v-108z" />
                  </svg>
                )}

                <Link
                  className={`subject ${!item.isRoot ? 'read' : ''}`}
                  href={`/threads/${item.slug}`}
                >
                  {item.title || '(无标题)'}
                </Link>

                <span className="metadata">
                  {item.parentTitle && (
                    <span className="author-name" title={`回复: ${item.parentTitle}`}>
                      → {item.parentTitle}
                    </span>
                  )}
                  {item.node && <span className="category">({item.node.name})</span>}
                  <span className="tail">
                    <RelativeTime timestamp={item.createTime} />
                  </span>
                </span>
              </div>
            </li>
          );
        })}
      </ul>
    </div>
  );
}