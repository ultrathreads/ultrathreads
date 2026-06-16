// components/features/ThreadListView.tsx
import ThreadTree from '@/components/features/ThreadTree';
import TopicPagination from '@/components/features/TopicPagination';
import EmptyTip from '@/components/ui/EmptyTip';
import { buildThreadTree } from '@/lib/utils/thread-tree';

interface Props {
  t: Record<string, any>;
  threadResult: any;
  nodeResult: { node: any; error: any };
  currentPage: number;
  safeNodeSlug: string;
}

export default function ThreadListView({ t, threadResult, nodeResult, currentPage, safeNodeSlug }: Props) {
  const { posts, paging, lastReadAtMap, error } = threadResult;
  const { node } = nodeResult;

  if (error) {
    return <EmptyTip text={t('common:loadFailed')} variant="error" />;
  }

  const viewPosts = buildThreadTree(posts, { lastReadAtMap });

  return (
    <>
      {/* ✅ 直接传数据，空状态由 ThreadTree 内部处理 */}
      <ThreadTree
        threads={viewPosts}
        activeNode={node}
        emptyText={t('common:noThreads')}
      />

      {/* 分页仍然在外层控制，因为即使没帖子也可能需要显示"第1页" */}
      {viewPosts.length > 0 && (
        <TopicPagination
          totalItems={paging.total}
          pageSize={paging.pageSize}
          currentPage={paging.page}
        />
      )}
    </>
  );
}