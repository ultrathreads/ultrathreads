// services/thread-service.ts
import { fetchThreads, type ThreadsData } from '@/lib/api/posts';

export interface ThreadPageData {
  posts: ThreadsData['results'];
  paging: ThreadsData['page'];
  error: string | null;
}

export async function getThreadPageData(page: number): Promise<ThreadPageData> {
  const safePage = Math.max(1, Number.isNaN(page) ? 1 : page);

  try {
    const data = await fetchThreads(safePage);
    return {
      posts: data.results ?? [],
      paging: data.page,
      error: null,
    };
  } catch (err) {
    console.error('[ThreadService] Fetch failed:', err);
    return {
      posts: [],
      paging: { page: safePage, limit: 20, total: 0 },
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}