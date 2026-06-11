// app/(main)/my/posts/page.tsx
import type { Metadata } from 'next';
import { redirect } from 'next/navigation';
import { getServerTranslation } from '@/lib/i18n/i18n-server';
import MyPostsList from '@/components/features/MyPostsList';
import TopicPagination from '@/components/TopicPagination';
import EmptyTip from '@/components/ui/EmptyTip';
import { getUserPostsPageData } from '@/services/my-post-service';
import { getCurrentUser } from '@/services/auth.server';

interface Props {
  searchParams: Promise<{ page?: string }>;
}

export async function generateMetadata({ searchParams }: Props): Promise<Metadata> {
  const { page: pageStr } = await searchParams;
  const page = Math.max(1, parseInt(pageStr || '1', 10));
  const title = page > 1 ? `我的帖子 - 第 ${page} 页` : '我的帖子';
  return { title };
}

export default async function MyPostsPage({ searchParams }: Props) {
  const t = await getServerTranslation(['common']);
  const { page: pageStr } = await searchParams;
  const currentPage = Math.max(1, parseInt(pageStr || '1', 10));

  const currentUser = await getCurrentUser();
  if (!currentUser?.id) {
    redirect('/login?callback=/my/posts');
  }

  const { posts, paging, error } = await getUserPostsPageData(
    currentUser.id,
    currentPage
  );

  if (error) {
    return <EmptyTip text={t('common:loadFailed')} variant="error" />;
  }

  return (
    <>
      {posts.length > 0 ? (
        <MyPostsList
          initialPosts={posts}
          initialPaging={paging}
        />
      ) : (
        <EmptyTip text={t('common:noPosts')} />
      )}

      {posts.length > 0 && (
        <TopicPagination
          totalItems={paging.totalItems}
          pageSize={paging.pageSize}
          currentPage={paging.currentPage}
        />
      )}
    </>
  );
}