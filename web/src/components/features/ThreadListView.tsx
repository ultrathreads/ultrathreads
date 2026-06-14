// components/features/ThreadListView.tsx
import ThreadTree from '@/components/features/ThreadTree';
import TopicPagination from '@/components/features/TopicPagination';
import EmptyTip from '@/components/ui/EmptyTip';
import { buildThreadTree } from '@/lib/utils/thread-tree';
import type { BackState } from '@/types/views';

interface Props {
  t: Record<string, any>; // 根据你的 i18n 类型调整
  threadResult: any;      // 根据实际类型调整
  nodeResult: { node: any; error: any }; // 根据实际类型调整
  page: number;
  safeNodeSlug: string;
}

export default function ThreadListView({ t, threadResult, nodeResult, page, safeNodeSlug }: Props) {
  const { posts, paging, lastReadAtMap, error } = threadResult;
  const { node } = nodeResult;

  if (error) {
    return <EmptyTip text={t('common:loadFailed')} variant="error" />;
  }

  const viewPosts = buildThreadTree(posts, { lastReadAtMap });

  const backState: BackState = {};
  if (safeNodeSlug) backState.nodeSlug = safeNodeSlug;
  if (page > 1) backState.page = String(page);

  return (
    <>
      {viewPosts.length > 0 ? (
        <ThreadTree threads={viewPosts} activeNode={node} backState={backState} />
      ) : (
        <EmptyTip text={t('common:noThreads')} />
      )}

      {viewPosts.length > 0 && (
        <TopicPagination
          totalItems={paging.totalItems}
          pageSize={paging.pageSize}
          page={paging.page}
        />
      )}
    </>
  );
}