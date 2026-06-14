// app/nodes/[slug]/page.tsx
import { getServerTranslation } from '@/lib/i18n/i18n-server';
import { getThreadPageData } from '@/services/thread-service';
import { getNodeDetail } from '@/services/node-service';
import ThreadListView from '@/components/features/ThreadListView';
import type { Metadata } from 'next';
import { notFound } from 'next/navigation';

// 定义路由参数接口
interface Props {
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ page?: string }>;
}

// 动态生成 Metadata
export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const { node } = await getNodeDetail(slug);
  return { title: node ? `${node.name} - 帖子列表` : '帖子列表' };
}

export default async function NodePage({ params, searchParams }: Props) {
  const t = await getServerTranslation(['common', 'home']);
  const { slug } = await params;
  const { page: pageStr } = await searchParams;
  
  const page = Math.max(1, parseInt(pageStr || '1', 10));
  
  if (!slug || typeof slug !== 'string' || slug.trim() === '') {
    notFound();
  }
  // 并发获取数据
  const [threadResult, nodeResult] = await Promise.all([
    getThreadPageData(page, slug),
    getNodeDetail(slug),
  ]);


  return (
    <ThreadListView 
      t={t} 
      threadResult={threadResult} 
      nodeResult={nodeResult} 
      page={page} 
      safeNodeSlug={slug} 
    />
  );
}