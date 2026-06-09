// components/features/MyPostsList.tsx
'use client';

import Link from 'next/link';
import type { MyPostListItem } from '@/services/my-post-service';
import { RelativeTime } from '@/components/RelativeTime';

interface Props {
  initialPosts: MyPostListItem[];
  initialPaging: { currentPage: number; pageSize: number; totalItems: number };
}

/**
 * 我的帖子列表（纯展示组件）
 * 分页由服务端 TopicPagination 驱动，此组件仅负责 SSR 首屏渲染
 * 零新增 CSS，完全复用 ThreadTree / ThreadItem 已有样式类
 */
export default function MyPostsList({ initialPosts }: Props) {
  return (
    <div className="thread-tree-container">
      <div className="thread-tree-header">
        <h1>我的帖子</h1>
      </div>

      <ul className="thread">
        {initialPosts.map((item) => (
          <li key={item.id}>
            <div className="entry">
              {item.parentId == null ? (
                <svg className="icon-topic" width="14" height="14" viewBox="0 0 14 14" fill="#3498db">
                  <circle cx="7" cy="7" r="5" />
                </svg>
              ) : (
                <svg className="icon-reply-svg" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 264 264" width="14" height="14">
                  <path d="M6,16 v108 a57 57,0 0 0 57,57 h92 v27 l45.5-45.5-45.5-45.5 v27 h-92 a20 20,0 0 1-20-20 v-108z" />
                </svg>
              )}

              <Link
                className={`subject ${item.parentId != null ? 'read' : ''}`}
                href={`/post/${item.id}`}
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
        ))}
      </ul>
    </div>
  );
}