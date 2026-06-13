// components/features/ThreadTree.tsx
'use client';

import { useState, useMemo, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import type { NodeEntity } from '@/types/domain';
import type { ThreadViewItem, BackState } from '@/types/view'; 
import { markNodeAsRead } from '@/services/node-service';

import ThreadItem from '@/components/features/ThreadItem';
import NodeHeader, { type HeaderDisplayData } from '@/components/features/NodeHeader';

interface Props {
  threads: ThreadViewItem[];
  activeNode: NodeEntity | null;
  activeTag?: HeaderDisplayData | null;
  backState?: BackState;
}

/**
 * 构建带回溯参数的详情页链接
 */
function buildPostUrl(postSlug: string, backState?: BackState): string {
  if (!backState || (!backState.nodeSlug && !backState.tagSlug && !backState.page)) {
    return `/threads/${postSlug}`;
  }

  const params = new URLSearchParams();
  if (backState.nodeSlug) params.set('nodeSlug', backState.nodeSlug);
  if (backState.tagSlug) params.set('tagSlug', backState.tagSlug);
  if (backState.page) params.set('page', backState.page);

  return `/threads/${postSlug}?${params.toString()}`;
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

export default function ThreadTree({ threads, activeNode, activeTag, backState }: Props) {
  const router = useRouter();

  const [allCollapsed, setAllCollapsed] = useState(false);
  const [sort, setSort] = useState('reply');
  const [markingRead, setMarkingRead] = useState(false);

  const effectiveSlug = activeNode?.slug;

  // 排序逻辑本身较复杂，保留 useMemo 是正确的
  const tree = useMemo(() => sortThreads(threads, sort), [threads, sort]);

  // 将纯逻辑函数提取出组件外部，避免每次渲染都重新创建
  const toggleAll = useCallback(() => setAllCollapsed((prev) => !prev), []);

  // 修复 handleMarkAsRead 的依赖项问题
  // 原代码依赖了 markingRead，导致每次 markingRead 变化时都会重新创建该函数
  // 实际上函数内部通过 setMarkingRead(true) 已经更新了状态，无需将其作为依赖
  const handleMarkAsRead = useCallback(async () => {
    if (!effectiveSlug) {
      console.warn('[ThreadTree] 标记已读跳过: 无法获取有效 Slug', { activeNode, backState });
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
  }, [effectiveSlug, activeNode, backState, router]);

  const isMarkReadDisabled = markingRead || !effectiveSlug;

  // 简化 headerTagData 的映射逻辑
  // 原代码使用 useMemo 仅为了做简单的属性映射，开销大于收益，直接赋值即可
  const headerTagData: HeaderDisplayData | null = activeTag 
    ? { name: activeTag.tagName } 
    : null;

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
            backState={backState}
          />
        ))}
      </ul>
    </div>
  );
}