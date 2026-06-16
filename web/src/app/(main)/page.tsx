// app/(main)/page.tsx
import { getServerTranslation } from '@/lib/i18n/i18n-server';
import { getThreadPageData } from '@/services/thread-service';
import ThreadListView from '@/components/features/ThreadListView';

interface Props {
  searchParams: Promise<{ page?: string }>;
}

export default async function HomePage({ searchParams }: Props) {
  const t = await getServerTranslation(['common', 'home']);
  
  const { page: pageStr } = await searchParams;
  const currentPage = Math.max(1, parseInt(pageStr || '1', 10));

  // 首页固定为第1页，且没有 nodeSlug
  const [threadResult] = await Promise.all([
    getThreadPageData(currentPage, ""),
  ]);

  return (
    <ThreadListView 
      t={t} 
      threadResult={threadResult} 
      nodeResult={{ node: null, error: null }} 
      currentPage={1} 
      safeNodeSlug="" 
    />
  );
}