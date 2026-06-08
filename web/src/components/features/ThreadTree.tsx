// components/features/ThreadTree.tsx
'use client';

import { useState, useMemo } from 'react';
import Link from 'next/link';
import type { NodeEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';

import ThreadItem from './ThreadItem';
import NodeHeader from './NodeHeader';

/** 从列表页透传的回溯状态 */
export interface BackState {
  nodeId?: string;
  page?: string;
}

interface Props {
  threads: ThreadViewItem[];
  activeNode: NodeEntity | null;
  backState?: BackState;
}

/**
 * 构建带回溯参数的详情页链接
 * 无参数时返回干净 URL，避免冗余查询字符串
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

/** 客户端排序函数 */
function sortThreads(threads: ThreadViewItem[], sortType: string): ThreadViewItem[] {
  const sorted = [...threads];
  switch (sortType) {
    case 'latest':
      return sorted.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());
    case 'reply':
      return sorted.sort((a, b) => {
        const diff = b.lastCommentTime - a.lastCommentTime;
        return diff !== 0 ? diff : b.id - a.id;
      });
    case 'most':
      return sorted.sort((a, b) => (b.replies?.length || 0) - (a.replies?.length || 0));
    case 'hot':
      return sorted.sort((a, b) => {
        const scoreA = (a.replies?.length || 0) * 10 - new Date(a.date).getTime() / 1e12;
        const scoreB = (b.replies?.length || 0) * 10 - new Date(b.date).getTime() / 1e12;
        return scoreB - scoreA;
      });
    default:
      return sorted;
  }
}

export default function ThreadTree({ threads, activeNode, backState }: Props) {
  const [allCollapsed, setAllCollapsed] = useState(false);
  const [sort, setSort] = useState('reply');

  // ✅ 先排序，再构建树
  const tree = useMemo(() => {
    return sortThreads(threads, sort);
  }, [threads, sort]);

  const toggleAll = () => setAllCollapsed((prev) => !prev);

  return (
    <div className="thread-tree-container">
      <div className="thread-tree-header">
        <NodeHeader node={activeNode} />

        <div className="thread-tree-actions">
          <select
            className="sort-select"
            value={sort}
            onChange={(e) => setSort(e.target.value)}
          >
            <option value="latest">最新发布</option>
            <option value="reply">最新回复</option>
            <option value="most">最多回复</option>
            <option value="hot">综合热门</option>
          </select>
          <button
            className={`collapse-all-btn ${allCollapsed ? 'is-collapsed' : ''}`}
            onClick={toggleAll}
            aria-label={allCollapsed ? '展开所有回帖' : '折叠所有回帖'}
          >
            <svg width="12" height="12" viewBox="0 0 12 12">
              <path d="M2 4l4 4 4-4" fill="none" stroke="currentColor" strokeWidth="1.5" />
            </svg>
            <span className="collapse-all-text">{allCollapsed ? '展开回帖' : '折叠回帖'}</span>
          </button>
        </div>
      </div>

      <ul className="thread">
        {tree.map((t) => (
          <ThreadItem
            key={t.id}
            item={t}
            isRoot
            globalCollapsed={allCollapsed}
            backState={backState}
          />
        ))}
      </ul>
    </div>
  );
}