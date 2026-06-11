// app/(main)/my/posts/page.tsx
import type { Metadata } from 'next';
import { revalidatePath } from 'next/cache'; 
import { redirect } from 'next/navigation';
import { getServerTranslation } from '@/lib/i18n/i18n-server';
import MyPostsList from '@/components/features/MyPostsList';
import MyFavoritesList from '@/components/features/MyFavoritesList';
import { getFavoritesPageData, deleteFavorite } from '@/services/favorite-service';
import TopicPagination from '@/components/TopicPagination';
import EmptyTip from '@/components/ui/EmptyTip';
import MyPostsTabs from '@/components/ui/MyPostsTabs';
import { 
  getUserRootPostsPageData, 
  getUserReplyPostsPageData 
} from '@/services/my-post-service';
import { getCurrentUser } from '@/services/auth.server';

interface Props {
  searchParams: Promise<{ page?: string; tab?: string }>;
}

export async function generateMetadata({ searchParams }: Props): Promise<Metadata> {
  const { page: pageStr, tab } = await searchParams;
  const page = Math.max(1, parseInt(pageStr || '1', 10));
  
  let title = '我的主帖';
  if (tab === 'bookmarks') {
    title = page > 1 ? `我的书签 - 第 ${page} 页` : '我的书签';
  } else if (tab === 'replies') {
    title = page > 1 ? `我的回帖 - 第 ${page} 页` : '我的回帖';
  } else {
    // 默认情况（root 或其他）
    title = page > 1 ? `我的主帖 - 第 ${page} 页` : '我的主帖';
  }

  return { title };
}

export default async function MyPostsPage({ searchParams }: Props) {
  const t = await getServerTranslation(['common']);
  const params = await searchParams;
  const currentPage = Math.max(1, parseInt(params.page || '1', 10));
  
  // 默认 Tab 设为 'root'
  const currentTab = params.tab || 'root';

  const currentUser = await getCurrentUser();
  if (!currentUser?.id) {
    redirect('/login?callback=/my/posts');
  }

  // 根据当前 Tab 渲染不同的内容
  const renderContent = async () => {
    if (currentTab === 'bookmarks') {
      return <BookmarksContent currentPage={currentPage} t={t} />;
    } else if (currentTab === 'replies') {
      return <ReplyPostsContent userId={currentUser.id} currentPage={currentPage} t={t} tab={currentTab} />;
    }

    // 默认渲染根帖列表（无论是 'root' 还是其他未匹配的值，都兜底显示根帖）
     return <RootPostsContent userId={currentUser.id} currentPage={currentPage} t={t} tab={currentTab} />;
  };

  return (
    <>
      {/* 渲染 Tab 导航 */}
      <MyPostsTabs />

      {/* 渲染对应内容 */}
      {await renderContent()}
    </>
  );
}

// 主帖内容组件
async function RootPostsContent({ userId, currentPage, t, tab }: { userId: string; currentPage: number; t: any; tab: 'root' | 'replies' }) {
  const { posts, paging, error } = await getUserRootPostsPageData(userId, currentPage);

  if (error) {
    return <EmptyTip text={t('common:loadFailed')} variant="error" />;
  }

  return (
    <>
      {posts.length > 0 ? (
        <MyPostsList initialPosts={posts} initialPaging={paging} tab={tab} />
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

// 回帖内容组件
async function ReplyPostsContent({ userId, currentPage, t, tab }: { userId: string; currentPage: number; t: any; tab: 'root' | 'replies' }) {
  const { posts, paging, error } = await getUserReplyPostsPageData(userId, currentPage);

  if (error) {
    return <EmptyTip text={t('common:loadFailed')} variant="error" />;
  }

  return (
    <>
      {posts.length > 0 ? (
        // TODO: 后续可以替换为专门的回帖列表组件 <MyRepliesList />
        <MyPostsList initialPosts={posts} initialPaging={paging} tab={tab} />
      ) : (
        <EmptyTip text="暂无回帖记录" />
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

// 更新：书签内容组件
async function BookmarksContent({ currentPage, t }: { currentPage: number; t: any }) {
  const { favorites, paging, error } = await getFavoritesPageData(currentPage);

  if (error) {
    return <EmptyTip text={t('common:loadFailed')} variant="error" />;
  }

  // 定义 Server Action
  async function handleDeleteFavoriteAction(formData: FormData) {
    'use server';
    
    const entityType = formData.get('entityType') as string;
    const entityId = Number(formData.get('entityId'));

    if (!entityType || !entityId) {
      // 1. 参数校验失败，直接抛出错误，客户端会捕获到它
      throw new Error('参数缺失：无法执行删除操作');
    }

    try {
      await deleteFavorite(entityType, entityId);
      
      // 2. 删除成功后，重新验证当前页面的数据
      revalidatePath(`/my/posts?tab=bookmarks&page=${currentPage}`);
      
      // 3. 如果没有抛出错误，代表操作成功
      // 此时不需要显式 return，客户端的 await 会正常走完
    } catch (err) {
      console.error('Delete favorite failed:', err);
      
      // 4. 捕获数据库等底层错误，并重新抛出给客户端
      throw new Error('操作失败，请重试');
    }
  }

  return (
    <>
      {favorites.length > 0 ? (
        // 将 Server Action 传递给客户端组件
        <MyFavoritesList 
          initialFavorites={favorites} 
          onDeleteFavoriteAction={handleDeleteFavoriteAction} 
        />
      ) : (
        <EmptyTip text="暂无书签记录" />
      )}

      {favorites.length > 0 && (
        <TopicPagination
          totalItems={paging.totalItems}
          pageSize={paging.pageSize}
          currentPage={paging.currentPage}
        />
      )}
    </>
  );
}
