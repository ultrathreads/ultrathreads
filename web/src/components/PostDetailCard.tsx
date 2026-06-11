// components/PostDetailCard.tsx
'use client';

import { useState, useCallback, useEffect } from 'react';
import Link from 'next/link';
import { toast } from 'sonner';
import type { PostEntity } from '@/types/domain';
import { RelativeTime } from '@/components/RelativeTime';
import Avatar from '@/components/ui/Avatar';
import AuthorLink from '@/components/ui/AuthorLink';
import { likePost, favoritePost } from '@/services/post-service';
import { ApiBusinessError } from '@/lib/api/client';

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
  const [likeCount, setLikeCount] = useState(post.likeCount);
  const [favCount, setFavCount] = useState(post.favoriteCount ?? 0);

  // ✅ props 刷新时同步本地计数
  useEffect(() => setLikeCount(post.likeCount), [post.likeCount]);
  useEffect(() => setFavCount(post.favoriteCount ?? 0), [post.favoriteCount]);

  /**
   * ✅ Sonner 版通用处理器：无需手动管理 loading / timer
   */
  const handleAction = useCallback(
    async (
      actionFn: () => Promise<void>,
      countSetter: React.Dispatch<React.SetStateAction<number>>,
      actionName: string
    ) => {
      // 🚀 乐观更新
      countSetter((prev) => prev + 1);

      // ✅ toast.promise 自动处理 loading → success/error 状态切换
      toast.promise(actionFn(), {
        loading: `${actionName}中...`,
        success: `${actionName}成功`,
        error: (err) => {
          // ❌ 失败回滚
          countSetter((prev) => prev - 1);

          if (err instanceof ApiBusinessError) {
            if (err.code === 401 || err.message === 'AUTH_REQUIRED') {
              return `请先登录后再${actionName}`;
            }
            return err.message || `${actionName}失败`;
          }
          return `${actionName}失败，请稍后重试`;
        },
      });
    },
    []
  );

  const handleLike = () =>
    handleAction(() => likePost(post.id), setLikeCount, '点赞');

  const handleFavorite = () =>
    handleAction(() => favoritePost(post.id), setFavCount, '收藏');

  return (
    <div className="post-detail-card">
      <h1 className="post-detail-title">{post.title}</h1>

      <div className="post-detail-meta">
        <Avatar
          className="detail-author-avatar"
          src={post.user.avatar}
          alt={post.user.nickname}
        />
        <AuthorLink 
          author={post.user.nickname} 
          authorId={post.user.id} 
          className="author-name" 
        />
        <RelativeTime timestamp={post.createTime} />
        {post.node && (
          <Link href={`/?nodeId=${post.node.nodeId}`} className="detail-tag">
            {post.node.name}
          </Link>
        )}

        {post.tags && post.tags.length > 0 && (
          <div className="post-tag">
            {post.tags.map((tag) => (
              <Link 
                key={tag.tagId} 
                href={`/tags/${tag.tagId}`} 
                className="detail-tag"
              >
                {tag.tagName}
              </Link>
            ))}
          </div>
        )}

        <span>阅读 {post.viewCount.toLocaleString()}</span>
        <span>回复 {replyCount}</span>
      </div>

      {post.content ? (
        <div className="post-detail-body" dangerouslySetInnerHTML={{ __html: post.content }} />
      ) : (
        <div className="post-detail-body post-detail-body--empty">暂无内容</div>
      )}

      {/* ✅ 不再需要 <InlineToast />，Sonner 通过 Portal 渲染 */}

      <div className="post-detail-actions">
        <button className="detail-action-btn" onClick={handleLike} type="button">
          👍 点赞 ({likeCount})
        </button>

        <button className="detail-action-btn" onClick={handleFavorite} type="button">
          ⭐ 收藏 ({favCount})
        </button>

        <button
          className={`detail-action-btn ${isEditorOpen ? 'detail-action-btn--active' : ''}`}
          onClick={onReplyClick}
          type="button"
          aria-expanded={isEditorOpen}
          aria-controls="reply-editor"
        >
          {isEditorOpen ? '✖ 不想说了' : '💬 我说两句'}
        </button>
      </div>
    </div>
  );
}