// app/(main)/page.tsx
import type { Metadata } from 'next';
import { getServerTranslation } from '@/lib/i18n/i18n-server';
import ThreadTree from '@/components/features/ThreadTree';
import TopicPagination from '@/components/TopicPagination';
import { getThreadPageData } from '@/services/thread-service';
import { getNodeDetail } from '@/services/node-service';
import { adaptToThreadView } from '@/lib/utils/thread-adapter';

interface Props {
  searchParams: Promise<{ page?: string; nodeId?: string }>;
}

/** 定义需要透传给子组件的回溯状态类型 */
interface BackState {
  nodeId?: string;
  page?: string;
}

export async function generateMetadata({ searchParams }: Props): Promise<Metadata> {
  const { page: pageStr, nodeId } = await searchParams;
  const page = Math.max(1, parseInt(pageStr || '1', 10));
  let title = page > 1 ? `帖子列表 - 第 ${page} 页` : '帖子列表';

  if (nodeId) {
    const { node } = await getNodeDetail(Number(nodeId));
    if (node) title = `${node.name} - ${title}`;
  }
  return { title };
}

export default async function HomePage({ searchParams }: Props) {
  const t = await getServerTranslation(['common', 'home']);
  const { page: pageStr, nodeId } = await searchParams;
  const currentPage = Math.max(1, parseInt(pageStr || '1', 10));

  const parsedNodeId = nodeId ? Number(nodeId) : undefined;

  const [threadResult, nodeResult] = await Promise.all([
    getThreadPageData(currentPage, parsedNodeId),
    parsedNodeId ? getNodeDetail(parsedNodeId) : Promise.resolve({ node: null, error: null }),
  ]);

  const { posts, paging, error } = threadResult;
  const { node } = nodeResult;

  if (error) {
    return <div className="p-8 text-center text-red-500">{t('common:loadFailed')}</div>;
  }

  // ✅ 将传输模型适配为视图模型，使 nodeName 等字段正确传递到 ThreadItem
  const viewPosts = posts.map(adaptToThreadView);

  // ✅ 构建回溯状态：仅当有实际值时才传递，避免生成无意义的空参数
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
        <div className="p-8 text-center text-gray-400">{t('home:noThreads')}</div>
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