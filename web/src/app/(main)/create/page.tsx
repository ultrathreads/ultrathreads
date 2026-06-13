// src/app/create/page.tsx
import { getAllNodes } from '@/services/node-service';
import { PostForm } from '@/components/features/PostForm';
import EmptyTip from '@/components/ui/EmptyTip';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: '发布新主贴',
};

export default async function CreatePostPage() {
  const { nodes, error: nodesError } = await getAllNodes();

  if (nodesError) {
    console.error('加载板块失败:', nodesError);
    return <EmptyTip text="板块数据加载失败，请稍后刷新重试" variant="error" />;
  }

  return <PostForm nodes={nodes} />;
}