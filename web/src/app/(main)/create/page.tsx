// src/app/create/page.tsx
import { getAllNodes } from '@/services/node-service';
import { PostForm } from '@/components/features/PostForm';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: '发布新主题',
};

export default async function CreatePostPage() {
  const { nodes, error: nodesError } = await getAllNodes();

  if (nodesError) {
    console.error('加载板块失败:', nodesError);
  }

  return (
    <div className="main-body">
      <div className="post-form-container">
        <h1 className="post-form-header">✏️ 发布新主题</h1>
        <PostForm nodes={nodes} />
      </div>
    </div>
  );
}