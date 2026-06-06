// src/components/PostDetailCard.tsx
import { PostDetailData } from '@/lib/mock-data';

interface PostDetailCardProps {
  post: PostDetailData;
}

export default function PostDetailCard({ post }: PostDetailCardProps) {
  return (
    <div className="post-detail-card">
      <h1 className="post-detail-title">{post.title}</h1>
      <div className="post-detail-meta">
        <img className="detail-author-avatar" src={post.authorAvatar} alt={post.author} />
        <span className="author-name">{post.author}</span>
        <time dateTime={post.date}>{post.date}</time>
        <span className="detail-tag">{post.tag}</span>
        <span>阅读 {post.views.toLocaleString()}</span>
        <span>回复 {post.comments}</span>
      </div>
      {/* 使用 dangerouslySetInnerHTML 来渲染模拟的 HTML 内容 */}
      <div className="post-detail-body" dangerouslySetInnerHTML={{ __html: post.content }} />
      <div className="post-detail-actions">
        <button className="detail-action-btn">👍 点赞 ({post.likes})</button>
        <button className="detail-action-btn">⭐ 收藏 ({post.favorites})</button>
        <button className="detail-action-btn">💬 回复</button>
        <button className="detail-action-btn">🔗 分享</button>
      </div>
    </div>
  );
}