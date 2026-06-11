// app/(main)/tags/[tagId]/page.tsx
import type { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { getServerTranslation } from '@/lib/i18n/i18n-server';
import ThreadTree from '@/components/features/ThreadTree';
import TopicPagination from '@/components/TopicPagination';
import { getTagPageData } from '@/services/thread-service';
import { buildThreadTree } from '@/lib/utils/thread-tree';
import { getTagDetail } from '@/services/tag-service';

interface Props {
  params: Promise<{ tagId: string }>;
  searchParams: Promise<{ page?: string }>;
}

/** 定义需要透传给子组件的回溯状态类型 */
interface BackState {
  tagId: string;
  page?: string;
}

async function parseSearchParams(searchParams: Props['searchParams']) {
  const { page: pageStr } = await searchParams;
  const currentPage = Math.max(1, parseInt(pageStr || '1', 10));
  return { currentPage };
}

export async function generateMetadata({ searchParams }: Props): Promise<Metadata> {
  const { currentPage, tagId } = await parseSearchParams(searchParams);
  let title = currentPage > 1 ? `帖子列表 - 第 ${currentPage} 页` : '帖子列表';

  if (tagId) {
    const { tag } = await getTagDetail(tagId);
    if (tag) title = `${tag.tagName} - ${title}`;
  }
  return { title };
}

export default async function TagPage({ params, searchParams }: Props) {
  const t = await getServerTranslation(['common', 'home']);
  const { tagId } = await params;
  const { currentPage } = await parseSearchParams(searchParams);

  const numericTagId = Number(tagId);
  if (!Number.isFinite(numericTagId) || numericTagId <= 0) {
    notFound();
  }

  const [threadResult, tagResult] = await Promise.all([
    getTagPageData(numericTagId, currentPage),
    getTagDetail(numericTagId),
  ]);

  const { posts, paging, lastReadAtMap, error } = threadResult;
  const { tag } = tagResult;

  if (error) {
    return <div className="p-8 text-center text-red-500">{t('common:loadFailed')}</div>;
  }

  const viewPosts = buildThreadTree(posts, { lastReadAtMap });

  const backState: BackState = { tagId: String(numericTagId) };
  if (currentPage > 1) {
    backState.page = String(currentPage);
  }

  return (
    <>
      {viewPosts.length > 0 ? (
        <ThreadTree threads={viewPosts} activeTag={tag} backState={backState} />
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