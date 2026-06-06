'use client';
import { useState } from 'react';
import { Thread } from '@/types';
import ThreadItem from './ThreadItem';

interface Props {
  threads: Thread[];
}

export default function ThreadTree({ threads }: Props) {
  const [allCollapsed, setAllCollapsed] = useState(false);
  const [sort, setSort] = useState('latest');

  // 通过 CSS 类批量控制折叠（与原始 JS 行为一致）
  const toggleAll = () => {
    setAllCollapsed(!allCollapsed);
    // 实际 DOM 操作由子组件响应 prop 或 CSS 完成
    // 这里简化为重新渲染
  };

  return (
    <div className="thread-tree-container">
      <div className="thread-tree-header">
        <div className="board-title-wrapper">
          <span className="board-title-icon">💻</span>
          <div className="board-title-text">
            <div className="board-title-name">技术交流</div>
            <div className="board-title-desc">前端框架、后端架构、DevOps 等技术话题讨论区</div>
          </div>
        </div>
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
        {threads.map((t) => (
          <ThreadItem key={t.id} item={t} isRoot />
        ))}
      </ul>
    </div>
  );
}