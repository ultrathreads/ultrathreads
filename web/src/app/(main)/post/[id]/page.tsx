// src/app/read/[id]/page.tsx
import { notFound } from 'next/navigation';
import { getMockPostDetail } from '@/lib/mock-data';
import PostDetailCard from '@/components/PostDetailCard';
import ThreadItem from '@/components/ThreadItem';

interface Props {
  params: Promise<{ id: string }>;
}

export default async function ReadPage({ params }: Props) {
  const { id } = await params;
  const post = getMockPostDetail("3021");

  if (!post) {
    notFound();
  }

  return (
    <>
      {/* 只保留详情页专属内容 */}
      <div className="main-body">
        <div className="detail-back-bar">
          <a className="back-list-btn" href="/">← 返回列表</a>
        </div>
        <PostDetailCard post={post} />
        <div className="thread-tree-container">
          <div className="thread-tree-header">
            <div className="thread-tree-title">💬 回帖讨论 ({post.comments})</div>
            <div className="thread-tree-actions">
              <select className="sort-select" aria-label="回帖排序">
                <option value="oldest">最早回复</option>
                <option value="newest">最新回复</option>
                <option value="hot">最热回复</option>
              </select>
              <button className="collapse-all-btn" title="折叠/展开所有回帖">
                <svg width="12" height="12" viewBox="0 0 12 12" fill="#7f8c8d">
                  <path d="M2 4l4 4 4-4z" />
                </svg>
                <span className="collapse-all-text">折叠回帖</span>
              </button>
            </div>
          </div>
          <ul className="thread">
            <ThreadItem key={post.id} item={post} isRoot />
          </ul>
        </div>
      </div>
    </>
  );
}