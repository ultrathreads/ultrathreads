// src/components/PostDetailCard.tsx
import type { PostEntity } from '@/types/domain';
import { RelativeTime } from './RelativeTime';

interface PostDetailCardProps {
  post: PostEntity;
}

export default function PostDetailCard({ post }: PostDetailCardProps) {
  return (
    <div className="post-detail-card">
      <h1 className="post-detail-title">{post.title}</h1>
      
      <div className="post-detail-meta">
        {/* ✅ 用户信息从嵌套的 user 对象中获取 */}
        <img 
          className="detail-author-avatar" 
          src={post.user.avatar} 
          alt={post.user.nickname} 
        />
        <span className="author-name">{post.user.nickname}</span>
        
        {/* ✅ 使用 RelativeTime 组件避免 SSR Hydration Mismatch */}
        <RelativeTime timestamp={post.createTime} />
        
        {/* ✅ 节点/分类信息从嵌套的 node 对象中获取 */}
        {post.node && (
          <span className="detail-tag">{post.node.name}</span>
        )}
        
        {/* ✅ 字段名对齐后端接口规范 */}
        <span>阅读 {post.viewCount.toLocaleString()}</span>
        <span>回复 {post.commentCount}</span>
      </div>

      {/* ✅ content 字段后端已返回 HTML 字符串，直接渲染 */}
      {post.content ? (
        <div
          className="post-detail-body"
          dangerouslySetInnerHTML={{ __html: post.content }}
        />
      ) : (
        <div className="post-detail-body post-detail-body--empty">
          暂无内容
        </div>
      )}

      <div className="post-detail-actions">
        <button className="detail-action-btn">👍 点赞 ({post.likeCount})</button>
        {/* ⚠️ 注意：当前 API 未返回 favorites 字段，暂用 0 占位或后续扩展接口 */}
        <button className="detail-action-btn">⭐ 收藏 (0)</button>
        <button className="detail-action-btn">💬 回复</button>
        <button className="detail-action-btn">🔗 分享</button>
      </div>
    </div>
  );
}