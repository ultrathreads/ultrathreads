// app/(main)/page.tsx
import type { Metadata } from 'next';
import { getServerTranslation } from '@/lib/i18n/i18n-server';
import ThreadTree from '@/components/features/ThreadTree';
import TopicPagination from '@/components/TopicPagination';
import { getThreadPageData } from '@/services/thread-service';
import { getNodeDetail } from '@/services/node-service';

interface Props {
  searchParams: Promise<{ page?: string; nodeId?: string }>;
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
  
  // ✅ 将字符串 nodeId 转为数字或 undefined，方便下游服务处理
  const parsedNodeId = nodeId ? Number(nodeId) : undefined;

  // ✅ 并行获取：帖子列表现在携带了 nodeId 参数
  const [threadResult, nodeResult] = await Promise.all([
    getThreadPageData(currentPage, parsedNodeId), 
    parsedNodeId ? getNodeDetail(parsedNodeId) : Promise.resolve({ node: null, error: null }),
  ]);

  const { posts, paging, error } = threadResult;
  const { node } = nodeResult;

  if (error) {
    return <div className="p-8 text-center text-red-500">{t('common:loadFailed')}</div>;
  }

  // ✅ 严格保留原始 <> Fragment 结构，零新增标签
  return (
    <>
      {posts.length > 0 ? (
        <ThreadTree threads={posts} activeNode={node} />
      ) : (
        <div className="p-8 text-center text-gray-400">{t('home:noThreads')}</div>
      )}

      {posts.length > 0 && (
        <TopicPagination
          totalItems={paging.total}
          pageSize={paging.limit}
          currentPage={paging.page}
        />
      )}
    </>
  );
}