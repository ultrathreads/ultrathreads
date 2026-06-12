// app/(main)/tags/[slug]/page.tsx
import type { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { getServerTranslation } from '@/lib/i18n/i18n-server';
import EmptyTip from '@/components/ui/EmptyTip';
import ThreadTree from '@/components/features/ThreadTree';
import TopicPagination from '@/components/TopicPagination';
import { getTagPageData } from '@/services/thread-service';
import { buildThreadTree } from '@/lib/utils/thread-tree';
import { getTagDetail } from '@/services/tag-service';

interface Props {
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ page?: string }>;
}

/** 定义需要透传给子组件的回溯状态类型 */
interface BackState {
  tagSlug: string;
  page?: string;
}

async function parseSearchParams(searchParams: Props['searchParams']) {
  const { page: pageStr } = await searchParams;
  const currentPage = Math.max(1, parseInt(pageStr || '1', 10));
  return { currentPage };
}

export async function generateMetadata({ searchParams }: Props): Promise<Metadata> {
  const { currentPage, slug } = await parseSearchParams(searchParams);
  let title = currentPage > 1 ? `帖子列表 - 第 ${currentPage} 页` : '帖子列表';

  if (slug) {
    const { tag } = await getTagDetail(slug);
    if (tag) title = `${tag.tagName} - ${title}`;
  }
  return { title };
}

export default async function TagPage({ params, searchParams }: Props) {
  const t = await getServerTranslation(['common', 'home']);
  const { slug } = await params;
  const { currentPage } = await parseSearchParams(searchParams);

  if (!slug || typeof slug !== 'string' || slug.trim() === '') {
    notFound();
  }

  const [threadResult, tagResult] = await Promise.all([
    getTagPageData(slug, currentPage),
    getTagDetail(slug),
  ]);

  const { posts, paging, lastReadAtMap, error } = threadResult;
  const { tag } = tagResult;

  if (error) {
    return <EmptyTip text={t('common:loadFailed')} variant="error" />;
  }

  const viewPosts = buildThreadTree(posts, { lastReadAtMap });

  const backState: BackState = { tagSlug: String(slug) };
  if (currentPage > 1) {
    backState.page = String(currentPage);
  }

  return (
    <>
      {viewPosts.length > 0 ? (
        <ThreadTree threads={viewPosts} activeTag={tag} backState={backState} />
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