// src/app/edit/[slug]/page.tsx
import { notFound } from 'next/navigation';
import { getAllNodes } from '@/services/node-service';
import { getPostDetail } from '@/services/post-service';
import { PostForm } from '@/components/features/PostForm';

interface Props {
  params: Promise<{ slug: string }>;
}

export async function generateMetadata({ params }: Props) {
  const { slug } = await params;
  return {
    title: '编辑主题',
    robots: { index: false, follow: false },
  };
}

export default async function EditPostPage({ params }: Props) {
  const { slug } = await params;

  const [nodesResult, postDetail] = await Promise.all([
    getAllNodes(),
    getPostDetail(slug,  { noCache: true }).catch(() => null),
  ]);

  if (!postDetail) notFound();

  const initialData = {
    slug: postDetail.slug,
    title: postDetail.title ?? '',
    rawContent: postDetail.rawContent ?? '',
    nodeSlug: postDetail.node?.slug ?? '',
    tags: Array.isArray(postDetail.tags)
      ? postDetail.tags.map((t: any) => t.tagName || t.name || t).join(', ')
      : '',
  };

  return (
    <div className="main-body">
      <div className="post-form-container">
        <h1 className="post-form-header">✏️ 编辑主题</h1>
        <PostForm nodes={nodesResult.nodes} initialData={initialData} />
      </div>
    </div>
  );
}