// src/components/features/PostCard.tsx
'use client';

import { useState, useCallback, useEffect } from 'react';
import Link from 'next/link';
import { toast } from 'sonner';
import type { PostEntity } from '@/types/domain';
import { RelativeTime } from '@/components/ui/RelativeTime';
import Avatar from '@/components/ui/Avatar';
import { likePost, favoritePost } from '@/services/post-service';
import { ApiError } from '@/lib/api/client';
import AuthorLink from '@/components/ui/AuthorLink';
import { useAuth } from '@/hooks/use-auth';

// ✨ 接口名同步更新
interface PostCardProps {
  post: PostEntity;
  detailHref: string;
  replyCount?: number;
  onReplyClick?: () => void;
  onEditClick?: () => void;
  isReplyEditorOpen?: boolean;
  isEditEditorOpen?: boolean;
  isRoot?: boolean;
}

// ✨ 组件名 + 默认导出同步更新
export default function PostCard({
  post,
  detailHref,
  replyCount = 0,
  onReplyClick,
  onEditClick,
  isReplyEditorOpen = false,
  isEditEditorOpen = false,
  isRoot = false,
}: PostCardProps) {
  // ... 其余逻辑完全不变，保持原样即可
  const [likeCount, setLikeCount] = useState(post.likeCount);
  const [favCount, setFavCount] = useState(post.favoriteCount ?? 0);
  const [actionLoading, setActionLoading] = useState<'like' | 'favorite' | null>(null);

  const { user } = useAuth();
  const canEdit = user?.slug === post.user.slug;

  useEffect(() => setLikeCount(post.likeCount), [post.likeCount]);
  useEffect(() => setFavCount(post.favoriteCount ?? 0), [post.favoriteCount]);

  const handleAction = useCallback(
    async (
      actionFn: () => Promise<void>,
      countSetter: React.Dispatch<React.SetStateAction<number>>,
      actionName: string,
      loadingKey: 'like' | 'favorite'
    ) => {
      if (actionLoading) return;

      let prevCount = 0;
      countSetter((prev) => {
        prevCount = prev;
        return prev + 1;
      });
      setActionLoading(loadingKey);

      try {
        await toast.promise(actionFn(), {
          loading: `${actionName}中...`,
          success: `${actionName}成功`,
        });
      } catch (err) {
        countSetter(prevCount);
        if (err instanceof ApiError) {
          if (err.code === 401 || err.message === 'AUTH_REQUIRED') {
            toast.error(`请先登录后再${actionName}`);
          } else {
            toast.error(err.message || `${actionName}失败`);
          }
        } else {
          toast.error(`${actionName}失败，请稍后重试`);
        }
      } finally {
        setActionLoading(null);
      }
    },
    [actionLoading]
  );

  const handleLike = () =>
    handleAction(() => likePost(post.slug), setLikeCount, '点赞', 'like');
  const handleFavorite = () =>
    handleAction(() => favoritePost(post.slug), setFavCount, '收藏', 'favorite');

  return (
    <div
      id={`post-${post.slug}`}
      className="post-detail-card"
      style={{ marginBottom: '12px', scrollMarginTop: '100px' }}
    >
      {isRoot && (
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '12px' }}>
          <h2 className="post-detail-title" style={{ margin: 0 }}>{post.title}</h2>
          {canEdit && (
            <Link
              href={`/edit/${post.slug}`}
              className="detail-action-btn"
              style={{ height: '32px', fontSize: '0.8125rem', textDecoration: 'none' }}
            >
              ✏️ 编辑
            </Link>
          )}
        </div>
      )}

      <div className="post-detail-meta">
        <Avatar
          className="detail-author-avatar"
          src={post.user.avatar}
          alt={post.user.nickname}
        />
        <AuthorLink
          author={post.user.nickname}
          authorSlug={post.user.slug}
          className="author-name"
        />
        <RelativeTime timestamp={post.createTime} />

        {isRoot && (
          <>
            {post.node && (
              <Link href={`/nodes/${post.node.nodeSlug}`} className="detail-node">
                {post.node.name}
              </Link>
            )}

            {post.tags && post.tags.length > 0 && (
              <>
                {post.tags.map((tag) => (
                  <Link
                    key={tag.slug}
                    href={`/tags/${tag.slug}`}
                    className="detail-tag"
                  >
                    #{tag.name}
                  </Link>
                ))}
              </>
            )}
            <span>阅读 {post.viewCount.toLocaleString()}</span>
            <span>回复 {replyCount}</span>
          </>
        )}
      </div>

      {post.content ? (
        <div
          className="post-detail-body"
          dangerouslySetInnerHTML={{ __html: post.content }}
        />
      ) : (
        <div className="post-detail-body post-detail-body--empty">暂无内容</div>
      )}

      <div className="post-detail-actions">
        <button
          className="detail-action-btn"
          onClick={handleLike}
          disabled={actionLoading !== null}
          type="button"
        >
          👍 点赞 ({likeCount})
        </button>

        <button
          className="detail-action-btn"
          onClick={handleFavorite}
          disabled={actionLoading !== null}
          type="button"
        >
          ⭐ 收藏 ({favCount})
        </button>

        {!isRoot && canEdit && onEditClick && (
          <button
            className={`detail-action-btn ${isEditEditorOpen ? 'detail-action-btn--active' : ''}`}
            onClick={onEditClick}
            type="button"
            aria-expanded={isEditEditorOpen}
            aria-controls="reply-editor"
          >
            {isEditEditorOpen ? '✖ 取消编辑' : '✏️ 编辑'}
          </button>
        )}

        {onReplyClick && (
          <button
            className={`detail-action-btn ${isReplyEditorOpen ? 'detail-action-btn--active' : ''}`}
            onClick={onReplyClick}
            type="button"
            aria-expanded={isReplyEditorOpen}
            aria-controls="reply-editor"
          >
            {isReplyEditorOpen ? '✖ 不想说了' : '💬 我说两句'}
          </button>
        )}
      </div>
    </div>
  );
}