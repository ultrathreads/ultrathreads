// components/features/ThreadTree.tsx
'use client';

import { useState, useMemo, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import type { NodeEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';
import { markNodeAsRead } from '@/services/node-service';

import ThreadItem from '@/components/features/ThreadItem';
import NodeHeader, { type HeaderDisplayData } from '@/components/features/NodeHeader';

interface Props {
  threads: ThreadViewItem[];
  activeNode: NodeEntity | null;
  activeTag?: HeaderDisplayData | null;
}

/** 客户端排序函数 (保持不变) */
function sortThreads(threads: ThreadViewItem[], sortType: string): ThreadViewItem[] {
  const sorted = [...threads];
  sorted.sort((a, b) => {
    const pinA = a.isPinned ? 1 : 0;
    const pinB = b.isPinned ? 1 : 0;
    if (pinA !== pinB) return pinB - pinA;

    switch (sortType) {
      case 'latest':
        return new Date(b.date).getTime() - new Date(a.date).getTime();
      case 'reply': {
        const diff = b.lastCommentTime - a.lastCommentTime;
        return diff !== 0 ? diff : b.createTime - a.createTime;
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

export default function ThreadTree({ threads, activeNode, activeTag }: Props) { // ✅ 解构移除 backState
  const router = useRouter();

  const [allCollapsed, setAllCollapsed] = useState(false);
  const [sort, setSort] = useState('reply');
  const [markingRead, setMarkingRead] = useState(false);

  const effectiveSlug = activeNode?.slug;

  const tree = useMemo(() => sortThreads(threads, sort), [threads, sort]);

  const toggleAll = useCallback(() => setAllCollapsed((prev) => !prev), []);

  // ✅ 依赖项移除 backState
  const handleMarkAsRead = useCallback(async () => {
    if (!effectiveSlug) {
      console.warn('[ThreadTree] 标记已读跳过: 无法获取有效 Slug', { activeNode });
      return;
    }

    setMarkingRead(true);
    try {
      await markNodeAsRead(effectiveSlug);
      toast.success('标记已读成功');
      router.refresh();
    } catch (err) {
      console.error('标记已读失败:', err);
      toast.error('标记已读失败，请重试');
    } finally {
      setMarkingRead(false);
    }
  }, [effectiveSlug, activeNode, router]); // ✅ 依赖项精简

  const isMarkReadDisabled = markingRead || !effectiveSlug;

  const headerTagData: HeaderDisplayData | null = activeTag
    ? { name: activeTag.tagName }
    : null;

  // 版块页下 activeNode 存在 → 不显示；首页/标签页/我的 → 显示
  const showNodeLink = !activeNode;

  return (
    <div className="thread-tree-container">
      <div className="thread-tree-header">
        <NodeHeader node={activeNode} tag={headerTagData} />

        <div className="thread-tree-actions">
          <button
            className={`detail-action-btn ${isMarkReadDisabled ? 'is-disabled' : ''}`}
            onClick={handleMarkAsRead}
            disabled={isMarkReadDisabled}
            aria-label="标记当前节点/标签为已读"
            title={!effectiveSlug
              ? "当前无有效节点或标签，无法标记已读"
              : "将当前内容标记为已读"
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
            key={t.slug}
            item={t}
            isRoot
            globalCollapsed={allCollapsed}
            showNodeLink={showNodeLink}
          />
        ))}
      </ul>
    </div>
  );
}