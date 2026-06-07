// src/app/create/page.tsx
import { getAllNodes } from '@/services/node-service';
import { CreatePostForm } from './CreatePostForm';

export default async function CreatePostPage() {
  // ✅ 在服务端获取真实节点数据
  const { nodes, error } = await getAllNodes();

  if (error) {
    console.error('加载板块失败:', error);
  }

  return (
    <div className="main-body">
      <div className="create-post-container">
        <h1 className="create-post-header">✏️ 发布新主题</h1>
        {/* ✅ 将真实数据作为 props 传入客户端组件 */}
        <CreatePostForm nodes={nodes} />
      </div>
    </div>
  );
}