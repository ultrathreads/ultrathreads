// src/app/(main)/users/[slug]/page.tsx
import type { Metadata } from 'next';
import { getServerTranslation } from '@/lib/i18n/i18n-server';
import { getUserBySlug } from '@/services/user-service';
import { getUserRootPostsPageData } from '@/services/my-post-service';
import MyPostsList from '@/components/features/MyPostsList';
import TopicPagination from '@/components/features/TopicPagination';
import EmptyTip from '@/components/ui/EmptyTip';

interface Props {
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ page?: string }>;
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const user = await getUserBySlug(slug).catch(() => null);
  if (!user) return { title: '用户不存在' };
  return { title: `${user.nickname || user.username} 的公开帖子` };
}

export default async function UserPublicPostsPage({ params, searchParams }: Props) {
  const t = await getServerTranslation(['common']);
  const { slug } = await params;
  const { page: pageStr } = await searchParams;
  const page = Math.max(1, parseInt(pageStr || '1', 10));

  const user = await getUserBySlug(slug);
  const { posts, paging, error } = await getUserRootPostsPageData(user.slug, page);

  if (error) {
    return <EmptyTip text={t('common:loadFailed')} variant="error" />;
  }

  return (
    <>
      {posts.length > 0 ? (
        <MyPostsList
          initialPosts={posts}
          initialPaging={paging}
          user={user}
        />
      ) : (
        <EmptyTip text={t('common:noPosts')} />
      )}

      {posts.length > 0 && (
        <TopicPagination
          totalItems={paging.totalItems}
          pageSize={paging.pageSize}
          page={paging.page}
        />
      )}
    </>
  );
}