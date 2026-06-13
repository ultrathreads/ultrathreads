// src/app/edit/[slug]/page.tsx
import { notFound } from 'next/navigation';
import { getAllNodes } from '@/services/node-service';
import { getPostDetail } from '@/services/post-service';
import { PostForm } from '@/components/features/PostForm';
import type { Metadata } from 'next';

interface Props {
  params: Promise<{ slug: string }>;
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  return {
    title: '修改主帖',
    robots: { index: false, follow: false },
  };
}

export default async function EditPostPage({ params }: Props) {
  const { slug } = await params;

  const [nodesResult, postDetail] = await Promise.all([
    getAllNodes(),
    getPostDetail(slug, { noCache: true }).catch((err) => {
      console.error(`获取帖子详情失败 [${slug}]:`, err);
      return null;
    }),
  ]);

  if (!postDetail || !nodesResult?.nodes) {
    notFound();
  }

  const initialData = {
    slug: postDetail.slug,
    title: postDetail.title ?? '',
    rawContent: postDetail.rawContent ?? '',
    nodeSlug: postDetail.node?.slug ?? '',
    tags: Array.isArray(postDetail.tags)
      ? postDetail.tags.map((t: any) => t.tagName || t.name || t).join(', ')
      : '',
  };

  return <PostForm nodes={nodesResult.nodes} initialData={initialData} />;
}