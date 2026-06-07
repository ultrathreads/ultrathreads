import type { PostEntity } from '@/types/domain';
import { RelativeTime } from '@/components/RelativeTime';

interface PostDetailCardProps {
  post: PostEntity;
  replyCount?: number;
  onReplyClick?: () => void;
  isEditorOpen?: boolean;
}

export default function PostDetailCard({
  post,
  replyCount,
  onReplyClick,
  isEditorOpen = false,
}: PostDetailCardProps) {
  return (
    <div className="post-detail-card">
      <h1 className="post-detail-title">{post.title}</h1>

      <div className="post-detail-meta">
        <img className="detail-author-avatar" src={post.user.avatar} alt={post.user.nickname} />
        <span className="author-name">{post.user.nickname}</span>
        <RelativeTime timestamp={post.createTime} />
        {post.node && <span className="detail-tag">{post.node.name}</span>}
        <span>阅读 {post.viewCount.toLocaleString()}</span>
        <span>回复 {replyCount}</span>
      </div>

      {post.content ? (
        <div className="post-detail-body" dangerouslySetInnerHTML={{ __html: post.content }} />
      ) : (
        <div className="post-detail-body post-detail-body--empty">暂无内容</div>
      )}

      <div className="post-detail-actions">
        <button className="detail-action-btn">👍 点赞 ({post.likeCount})</button>
        <button className="detail-action-btn">⭐ 收藏 (0)</button>

        {/* ✅ 按钮样式 & 文案联动 */}
        <button
          className={`detail-action-btn ${isEditorOpen ? 'detail-action-btn--active' : ''}`}
          onClick={onReplyClick}
          type="button"
          aria-expanded={isEditorOpen}
          aria-controls="reply-editor"
        >
          {isEditorOpen ? '✖ 取消' : '💬 回复'}
        </button>
      </div>
    </div>
  );
}