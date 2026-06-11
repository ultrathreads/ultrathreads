// app/(main)/page.tsx
import type { Metadata } from 'next';
import { getServerTranslation } from '@/lib/i18n/i18n-server';
import ThreadTree from '@/components/features/ThreadTree';
import TopicPagination from '@/components/TopicPagination';
import EmptyTip from '@/components/ui/EmptyTip';
import { getThreadPageData } from '@/services/thread-service';
import { getNodeDetail } from '@/services/node-service';
import { buildThreadTree } from '@/lib/utils/thread-tree';

interface Props {
  searchParams: Promise<{ page?: string; nodeId?: string }>;
}

/** 定义需要透传给子组件的回溯状态类型 */
interface BackState {
  nodeId?: string;
  page?: string;
}

async function parseSearchParams(searchParams: Props['searchParams']) {
  const { page: pageStr, nodeId } = await searchParams;
  const currentPage = Math.max(1, parseInt(pageStr || '1', 10));

  const rawNodeId = nodeId ? Number(nodeId) : NaN;
  const validNodeId = Number.isFinite(rawNodeId) && rawNodeId > 0 
    ? rawNodeId 
    : undefined;

  return { currentPage, validNodeId };
}

export async function generateMetadata({ searchParams }: Props): Promise<Metadata> {
  const { currentPage, validNodeId } = await parseSearchParams(searchParams);
  let title = currentPage > 1 ? `帖子列表 - 第 ${currentPage} 页` : '帖子列表';

  if (validNodeId) {
    const { node } = await getNodeDetail(validNodeId);
    if (node) title = `${node.name} - ${title}`;
  }
  return { title };
}

export default async function HomePage({ searchParams }: Props) {
  const t = await getServerTranslation(['common', 'home']);
  const { currentPage, validNodeId } = await parseSearchParams(searchParams);

  const [threadResult, nodeResult] = await Promise.all([
    getThreadPageData(currentPage, validNodeId),
    // ✅ 现在 -1、0、NaN、undefined 全部被过滤，不会再发起无效请求
    validNodeId 
      ? getNodeDetail(validNodeId) 
      : Promise.resolve({ node: null, error: null }),
  ]);

  const { posts, paging, lastReadAtMap, error } = threadResult;
  const { node } = nodeResult;

  if (error) {
    return <EmptyTip text={t('common:loadFailed')} variant="error" />;
  }

  const viewPosts = buildThreadTree(posts, { lastReadAtMap });

  // 构建回溯状态：仅当有实际值时才传递，避免生成无意义的空参数
  const parsedNodeId = validNodeId ? Number(validNodeId) : undefined;
  const backState: BackState = {};
  if (parsedNodeId !== undefined) {
    backState.nodeId = String(parsedNodeId);
  }
  if (currentPage > 1) {
    backState.page = String(currentPage);
  }

  return (
    <>
      {viewPosts.length > 0 ? (
        // ✅ 将回溯状态透传给 ThreadTree，由其内部 Link 拼接到详情页 URL
        <ThreadTree threads={viewPosts} activeNode={node} backState={backState} />
      ) : (
        <EmptyTip text={t('common:noThreads')} />
      )}

      {viewPosts.length > 0 && (
        <TopicPagination
          totalItems={paging.total}
          pageSize={paging.limit}
          currentPage={paging.page}
        />
      )}
    </>
  );
}