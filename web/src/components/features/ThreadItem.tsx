'use client';
import { useState } from 'react';
import { Reply } from '@/types';
import { RelativeTime } from '@/components/RelativeTime';

interface Props {
  item: Reply;
  isRoot?: boolean;
  currentPostId?: string | number;
}

export default function ThreadItem({ item, isRoot, currentPostId }: Props) {
  const [folded, setFolded] = useState(false);
  const hasReplies = item.replies && item.replies.length > 0;

  // 在组件内部计算当前节点是否激活
  const isActive = currentPostId !== undefined && String(item.id) === String(currentPostId);

  return (
    <li className={folded ? 'folded' : ''}>
      <div className={`entry ${isActive ? 'active' : ''}`}>
        {isRoot && (
          <span className="fold-expand">
            {hasReplies ? (
              <a onClick={(e) => { e.preventDefault(); setFolded(!folded); }}>
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
        <a 
          className={`subject ${isRoot ? '' : 'read'} ${isActive ? 'active' : ''}`} 
          href={`/post/${item.id}`}
        >
          {item.title}
        </a>
        <span className="metadata">
          <span className="author-name">{item.author}</span>
          <span className="tail">
            <RelativeTime timestamp={item.date} />
          </span>
          {item.category && <span className="category">({item.category})</span>}
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
            <ThreadItem key={r.id} item={r} currentPostId={currentPostId} />
          ))}
        </ul>
      )}
    </li>
  );
}