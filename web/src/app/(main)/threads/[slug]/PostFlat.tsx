// src/app/threads/[slug]/PostFlat.tsx
'use client';

import { useState, useCallback } from 'react';
import type { PostEntity } from '@/types/domain';
import PostFlatItem from '@/components/PostFlatItem';
import ReplyEditor from '@/components/features/ReplyEditor';

interface PostFlatProps {
  posts: PostEntity[];
  totalReplyCount: number;
}

export function PostFlat({ posts, totalReplyCount }: PostFlatProps) {
  const [activeEditorPostSlug, setActiveEditorPostSlug] = useState<number | null>(null);

  // 切换编辑器状态（点击已打开的按钮时自动收起）
  const toggleEditor = useCallback((postSlug: string) => {
    setActiveEditorPostSlug((prev) => (prev === postSlug ? null : postSlug));
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
          const isEditorOpen = activeEditorPostSlug === post.slug;

          return (
            <div key={post.slug}>
              <PostFlatItem
                post={post}
                detailHref={`/threads/${post.slug}`}
                replyCount={isRoot ? (post.commentCount ?? totalReplyCount) : 0}
                isRoot={isRoot}
                // ✅ 仅根帖注入回复交互能力
                {...(true && {
                  onReplyClick: () => toggleEditor(post.slug),
                  isEditorOpen,
                })}
              />

              {/* ✅ 编辑器挂载在根帖卡片下方，非根帖不渲染 */}
              {isRoot && isEditorOpen && (
                <ReplyEditor
                  parentSlug={post.slug}
                  nodeSlug={post.node?.nodeSlug ?? 0}
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