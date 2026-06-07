// components/features/ThreadTree.tsx
'use client';

import { useState, useMemo } from 'react';
import type { ForumNode } from '@/lib/services/node-service';
import type { ThreadViewItem } from '@/types/view';
import { buildThreadTree } from '@/lib/utils/thread-utils';

import ThreadItem from './ThreadItem';
import NodeHeader from './NodeHeader';

interface Props {
  threads: ThreadViewItem[];
  activeNode: ForumNode | null;
}

/** 客户端排序函数 */
function sortThreads(threads: ThreadViewItem[], sortType: string): ThreadViewItem[] {
  const sorted = [...threads];
  switch (sortType) {
    case 'latest':
      return sorted.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());
    case 'reply':
      // 假设有 lastReplyTime 字段，若无则降级为 date
      return sorted.sort((a, b) => 
        new Date((b as any).lastReplyTime || b.date).getTime() - 
        new Date((a as any).lastReplyTime || a.date).getTime()
      );
    case 'most':
      return sorted.sort((a, b) => (b.replies?.length || 0) - (a.replies?.length || 0));
    case 'hot':
      // 简单热度算法：回复数权重 + 时间衰减
      return sorted.sort((a, b) => {
        const scoreA = (a.replies?.length || 0) * 10 - new Date(a.date).getTime() / 1e12;
        const scoreB = (b.replies?.length || 0) * 10 - new Date(b.date).getTime() / 1e12;
        return scoreB - scoreA;
      });
    default:
      return sorted;
  }
}

export default function ThreadTree({ threads, activeNode }: Props) {
  const [allCollapsed, setAllCollapsed] = useState(false);
  const [sort, setSort] = useState('latest');

  // ✅ 先排序，再构建树
  const tree = useMemo(() => {
    const sorted = sortThreads(threads, sort);
    return buildThreadTree(sorted);
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
          // ✅ 关键修复：将 allCollapsed 传递给子组件
          <ThreadItem 
            key={t.id} 
            item={t} 
            isRoot 
            globalCollapsed={allCollapsed} 
          />
        ))}
      </ul>
    </div>
  );
}