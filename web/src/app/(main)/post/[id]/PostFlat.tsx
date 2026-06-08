// src/app/post/[id]/PostFlat.tsx
'use client';

import { useState, useCallback } from 'react';
import type { PostEntity } from '@/types/domain';
import PostFlatItem from '@/components/PostFlatItem';
import ReplyEditor from '@/components/features/ReplyEditor'; // ✅ 直接复用现有编辑器

interface PostFlatProps {
  posts: PostEntity[];
  totalReplyCount: number;
}

export function PostFlat({ posts, totalReplyCount }: PostFlatProps) {
  // 记录当前激活编辑器的帖子ID，null 表示无编辑器打开
  const [activeEditorPostId, setActiveEditorPostId] = useState<number | null>(null);

  // 切换编辑器状态（点击已打开的按钮时自动收起）
  const toggleEditor = useCallback((postId: number) => {
    setActiveEditorPostId((prev) => (prev === postId ? null : postId));
  }, []);

  // 编辑器 autoFocus 消费后的回调，避免重复触发滚动/聚焦
  const handleAutoFocusConsumed = useCallback(() => {
    // 如需在聚焦完成后执行额外逻辑可在此扩展
  }, []);

  return (
    <div className="post-flat-container">
      {posts.length > 0 ? (
        posts.map((post, index) => {
          const isRoot = index === 0;
          const isEditorOpen = activeEditorPostId === post.id;

          return (
            <div key={post.id}>
              <PostFlatItem
                post={post}
                detailHref={`/post/${post.id}`}
                replyCount={isRoot ? (post.commentCount ?? totalReplyCount) : 0}
                isRoot={isRoot}
                // ✅ 仅根帖注入回复交互能力
                {...(isRoot && {
                  onReplyClick: () => toggleEditor(post.id),
                  isEditorOpen,
                })}
              />

              {/* ✅ 编辑器挂载在根帖卡片下方，非根帖不渲染 */}
              {isRoot && isEditorOpen && (
                <ReplyEditor
                  parentId={post.id}
                  nodeId={post.node?.nodeId ?? 0}
                  replyToTitle={post.title}
                  autoFocus
                  onAutoFocusConsumed={handleAutoFocusConsumed}
                />
              )}
            </div>
          );
        })
      ) : (
        <div className="text-center text-gray-400 py-8">暂无回复</div>
      )}

    </div>
  );
}