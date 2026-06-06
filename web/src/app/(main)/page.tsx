// app/page.tsx
import type { Metadata } from 'next';
import { getServerTranslation } from '@/lib/i18n-server';
import ThreadTree from '@/components/ThreadTree';
import TopicPagination from '@/components/TopicPagination';
import { getThreadPageData } from '@/services/thread-service';

interface Props {
  searchParams: Promise<{ page?: string }>;
}

export async function generateMetadata({ searchParams }: Props): Promise<Metadata> {
  const { page: pageStr } = await searchParams;
  const page = Math.max(1, parseInt(pageStr || '1', 10));

  return {
    title: page > 1 ? `帖子列表 - 第 ${page} 页` : '帖子列表',
  };
}

export default async function HomePage({ searchParams }: Props) {
  const t = await getServerTranslation(['common', 'home']);
  const { page: pageStr } = await searchParams;
  const currentPage = Math.max(1, parseInt(pageStr || '1', 10));

  // 👇 页面只负责调用 Service 并渲染，不感知 API 细节
  const { posts, paging, error } = await getThreadPageData(currentPage);

  if (error) {
    return <div className="p-8 text-center text-red-500">{t('common:loadFailed')}</div>;
  }

  return (
    <>
      {posts.length > 0 ? (
        <ThreadTree threads={posts} />
      ) : (
        <div className="p-8 text-center text-gray-400">{t('home:noThreads')}</div>
      )}

      <TopicPagination
        totalItems={paging.total}
        pageSize={paging.limit}
        currentPage={paging.page}
      />
    </>
  );
}