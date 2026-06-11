'use client';

import { useState, useCallback, useEffect } from 'react';
import Link from 'next/link';
import { toast } from 'sonner';
import type { PostEntity } from '@/types/domain';
import { RelativeTime } from '@/components/RelativeTime';
import Avatar from '@/components/ui/Avatar';
import { likePost, favoritePost } from '@/services/post-service';
import { ApiBusinessError } from '@/lib/api/client';
import AuthorLink from '@/components/ui/AuthorLink';

interface PostFlatItemProps {
  post: PostEntity;
  detailHref: string;
  replyCount?: number;
  onReplyClick?: () => void;
  isEditorOpen?: boolean;
  isRoot?: boolean;
}

export default function PostFlatItem({
  post,
  detailHref,
  replyCount = 0,
  onReplyClick,
  isEditorOpen = false,
  isRoot = false,
}: PostFlatItemProps) {
  const [likeCount, setLikeCount] = useState(post.likeCount);
  const [favCount, setFavCount] = useState(post.favoriteCount ?? 0);
  const [actionLoading, setActionLoading] = useState<'like' | 'favorite' | null>(null);

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
        if (err instanceof ApiBusinessError) {
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
    handleAction(() => likePost(post.id), setLikeCount, '点赞', 'like');
  const handleFavorite = () =>
    handleAction(() => favoritePost(post.id), setFavCount, '收藏', 'favorite');

  return (
    // 注入原生 id + scrollMarginTop 防顶部导航遮挡
    <div
      id={`post-${post.id}`}
      className="post-detail-card"
      style={{ marginBottom: '12px', scrollMarginTop: '100px' }}
    >
      {/* ✅ 仅根帖显示标题 */}
      {isRoot && <h2 className="post-detail-title">{post.title}</h2>}

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

        {/* 仅根帖显示节点标签、阅读数和回复数 */}
        {isRoot && (
          <>
            {post.node && (
              <Link href={`/?nodeId=${post.node.nodeId}`} className="detail-node">
                {post.node.name}
              </Link>
            )}

            {post.tags && post.tags.length > 0 && (
              <>
                {post.tags.map((tag) => (
                  <Link 
                    key={tag.tagId} 
                    href={`/tags/${tag.tagId}`} 
                    className="detail-tag"
                  >
                    #{tag.tagName}
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

        {onReplyClick && (
          <button
            className={`detail-action-btn ${isEditorOpen ? 'detail-action-btn--active' : ''}`}
            onClick={onReplyClick}
            type="button"
            aria-expanded={isEditorOpen}
            aria-controls="reply-editor"
          >
            {isEditorOpen ? '✖ 不想说了' : '💬 我说两句'}
          </button>
        )}
      </div>
    </div>
  );
}