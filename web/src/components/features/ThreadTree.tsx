// components/ThreadTree.tsx
'use client';

import { useState, useMemo } from 'react';
import type { SimplePost } from '@/lib/services/thread-service';
import type { ForumNode } from '@/lib/services/node-service';
import { buildThreadTree } from '@/lib/utils/thread-utils';
import ThreadItem from './ThreadItem';
import NodeHeader from './NodeHeader';

interface Props {
  threads: SimplePost[];
  activeNode: ForumNode | null;
}

export default function ThreadTree({ threads, activeNode }: Props) {
  const [allCollapsed, setAllCollapsed] = useState(false);
  const [sort, setSort] = useState('latest');

  const tree = useMemo(() => buildThreadTree(threads), [threads]);
  const toggleAll = () => setAllCollapsed((prev) => !prev);

  return (
    <div className="thread-tree-container">
      {/* ✅ 原始 thread-tree-header 结构 100% 保留 */}
      <div className="thread-tree-header">
        {/* ✅ 动态注入节点信息，输出的 DOM 和 className 与原来完全一致 */}
        <NodeHeader node={activeNode} />

        <div className="thread-tree-actions">
          <select className="sort-select" value={sort} onChange={(e) => setSort(e.target.value)}>
            <option value="latest">最新发布</option>
            <option value="reply">最新回复</option>
            <option value="most">最多回复</option>
            <option value="hot">综合热门</option>
          </select>
          <button
            className={`collapse-all-btn ${allCollapsed ? 'is-collapsed' : ''}`}
            onClick={toggleAll}
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
          <ThreadItem key={t.id} item={t} isRoot />
        ))}
      </ul>
    </div>
  );
}