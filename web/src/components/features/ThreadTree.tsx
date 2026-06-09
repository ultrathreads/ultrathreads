// components/features/ThreadTree.tsx
'use client';

import { useState, useMemo, useCallback } from 'react';
import type { NodeEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';
import { markNodeAsRead } from '@/services/node-service';

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

  sorted.sort((a, b) => {
    // ✅ 1. 最高优先级：置顶帖始终排在非置顶帖前面
    const pinA = a.isPinned ? 1 : 0;
    const pinB = b.isPinned ? 1 : 0;
    if (pinA !== pinB) {
      return pinB - pinA; // 降序：true(1) > false(0)
    }

    // ✅ 2. 次级优先级：仅在置顶状态相同时，才应用业务排序规则
    switch (sortType) {
      case 'latest':
        return new Date(b.date).getTime() - new Date(a.date).getTime();
      case 'reply': {
        const diff = b.lastCommentTime - a.lastCommentTime;
        return diff !== 0 ? diff : b.id - a.id;
      }
      case 'most':
        return (b.replies?.length || 0) - (a.replies?.length || 0);
      case 'hot': {
        const scoreA = (a.replies?.length || 0) * 10 - new Date(a.date).getTime() / 1e12;
        const scoreB = (b.replies?.length || 0) * 10 - new Date(b.date).getTime() / 1e12;
        return scoreB - scoreA;
      }
      default:
        return 0;
    }
  });

  return sorted;
}

export default function ThreadTree({ threads, activeNode, backState }: Props) {
  const [allCollapsed, setAllCollapsed] = useState(false);
  const [sort, setSort] = useState('reply');
  const [markingRead, setMarkingRead] = useState(false);

  // ✅ 提前计算有效的 nodeId，避免重复逻辑
  const effectiveNodeId = useMemo(() => {
    return activeNode?.nodeId ?? activeNode?.id;
  }, [activeNode]);

  // ✅ 排序逻辑保持不变
  const tree = useMemo(() => sortThreads(threads, sort), [threads, sort]);

  const toggleAll = () => setAllCollapsed((prev) => !prev);

  // ✅ 标记已读回调：使用预计算的 effectiveNodeId
  const handleMarkAsRead = useCallback(async () => {
    console.log('[ThreadTree] markAsRead clicked', { 
      nodeId: effectiveNodeId, 
      markingRead 
    });

    if (!effectiveNodeId) {
      console.warn('[ThreadTree] 标记已读跳过: 无法获取有效 nodeId', activeNode);
      return;
    }

    setMarkingRead(true);
    try {
      await markNodeAsRead(effectiveNodeId);
      // TODO: 成功后刷新未读状态
    } catch (err) {
      console.error('标记已读失败:', err);
    } finally {
      setMarkingRead(false);
    }
  }, [effectiveNodeId, markingRead, activeNode]);

  // ✅ 判断按钮是否应该禁用
  const isMarkReadDisabled = markingRead || !effectiveNodeId;

  return (
    <div className="thread-tree-container">
      <div className="thread-tree-header">
        <NodeHeader node={activeNode} />

        <div className="thread-tree-actions">
          {/* ✅ 标记已读按钮：增加视觉禁用态 + 精确的 disabled 条件 */}
          <button
            className={`detail-action-btn ${isMarkReadDisabled ? 'is-disabled' : ''}`}
            onClick={handleMarkAsRead}
            disabled={isMarkReadDisabled}
            aria-label="标记当前节点为已读"
            title={!effectiveNodeId 
              ? "当前无有效节点，无法标记已读" 
              : "将本节点所有帖子标记为已读"
            }
          >
            {markingRead ? (
              <span className="mark-read-loading">处理中…</span>
            ) : (
              <>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <polyline points="20 6 9 17 4 12" />
                </svg>
                <span className="mark-read-text">标记已读</span>
              </>
            )}
          </button>

          {/* 排序和折叠按钮保持不变 */}
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