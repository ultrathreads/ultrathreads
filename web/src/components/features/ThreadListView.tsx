// components/features/ThreadListView.tsx
import ThreadTree from '@/components/features/ThreadTree';
import TopicPagination from '@/components/features/TopicPagination';
import EmptyTip from '@/components/ui/EmptyTip';
import { buildThreadTree } from '@/lib/utils/thread-tree';

interface BackState {
  nodeSlug?: string;
  page?: string;
}

interface Props {
  t: Record<string, any>; // 根据你的 i18n 类型调整
  threadResult: any;      // 根据实际类型调整
  nodeResult: { node: any; error: any }; // 根据实际类型调整
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

  const backState: BackState = {};
  if (safeNodeSlug) backState.nodeSlug = safeNodeSlug;
  if (currentPage > 1) backState.page = String(currentPage);

  return (
    <>
      {viewPosts.length > 0 ? (
        <ThreadTree threads={viewPosts} activeNode={node} backState={backState} />
      ) : (
        <EmptyTip text={t('common:noThreads')} />
      )}

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