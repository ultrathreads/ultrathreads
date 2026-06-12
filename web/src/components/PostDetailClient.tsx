// components/PostDetailClient.tsx
'use client';

import { useState, useCallback } from 'react';
import PostDetailCard from '@/components/PostDetailCard';
import ReplyEditor from '@/components/features/ReplyEditor';
import ThreadItem from '@/components/features/ThreadItem';
import type { PostEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';
import type { BackState } from '@/components/features/ThreadTree';

interface PostDetailClientProps {
  post: PostEntity;
  viewPosts: ThreadViewItem[];
  totalReplyCount: number;
  backState: BackState;
}

export default function PostDetailClient({
  post,
  viewPosts,
  totalReplyCount,
  backState,
}: PostDetailClientProps) {
  // 默认隐藏编辑器
  const [showEditor, setShowEditor] = useState(false);
  const [replyToSlug, setReplyToSlug] = useState<string>(post.slug);
  const [replyToTitle, setReplyToTitle] = useState<string>(post.slug);
  const [shouldAutoFocus, setShouldAutoFocus] = useState(false);

  const handleThreadReplyClick = useCallback((targetSlug: string, targetTitle: string) => {
    setReplyToSlug(targetSlug);
    setReplyToTitle(`${post.title}(${post.user.nickname})`);
    setShowEditor(true);
    setShouldAutoFocus(true);
  }, [post.title, post.user.nickname]);

  return (
    <>
      <PostDetailCard
        post={post}
        replyCount={totalReplyCount}
        isEditorOpen={showEditor && replyToSlug === post.slug}
        onReplyClick={() => {
          // 简化切换逻辑，打开时必定触发聚焦
          setShowEditor((prev) => {
            if (!prev) {
              // 从隐藏变为显示 → 重置目标为主帖并触发聚焦
              setReplyToSlug(post.slug);
              setReplyToTitle(post.title);
              setShouldAutoFocus(true);
            }
            return !prev;
          });
        }}
      />

      {showEditor && (
        <div id="reply-editor">
          <ReplyEditor
            key={replyToSlug}
            parentSlug={replyToSlug}
            replyToTitle={replyToTitle}
            autoFocus={shouldAutoFocus}
            onAutoFocusConsumed={() => setShouldAutoFocus(false)}
          />
        </div>
      )}

      <div className="thread-tree-container">
        <div className="thread-tree-header">
          <div className="thread-tree-title">💬 回帖讨论 ({totalReplyCount})</div>
          <div className="thread-tree-actions">
            <select className="sort-select" aria-label="回帖排序" defaultValue="oldest">
              <option value="oldest">最早回复</option>
              <option value="newest">最新回复</option>
              <option value="hot">最热回复</option>
            </select>
          </div>
        </div>

        <ul className="thread">
          {viewPosts.map((reply) => (
            <ThreadItem
              key={reply.slug}
              item={reply}
              isRoot
              currentPostSlug={String(post.slug)}
              backState={backState}
              onReplyClick={handleThreadReplyClick}
            />
          ))}
        </ul>
      </div>
    </>
  );
}