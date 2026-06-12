// app/(main)/page.tsx
import { getServerTranslation } from '@/lib/i18n/i18n-server';
import { getThreadPageData } from '@/services/thread-service';
import ThreadListView from '@/components/features/ThreadListView';

export default async function HomePage() {
  const t = await getServerTranslation(['common', 'home']);
  
  // 首页固定为第1页，且没有 nodeSlug
  const [threadResult] = await Promise.all([
    getThreadPageData(1, ""),
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