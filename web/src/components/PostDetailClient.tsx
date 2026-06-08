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
  const [replyToId, setReplyToId] = useState<string | number>(post.id);
  const [replyToTitle, setReplyToTitle] = useState<string>(post.title);
  const [shouldAutoFocus, setShouldAutoFocus] = useState(false);

  const handleThreadReplyClick = useCallback((targetId: string | number, targetTitle: string) => {
    setReplyToId(targetId);
    setReplyToTitle(`${post.title}(${post.user.nickname})`);
    setShowEditor(true);
    setShouldAutoFocus(true);
  }, [post.title, post.user.nickname]);

  return (
    <>
      <PostDetailCard
        post={post}
        replyCount={totalReplyCount}
        isEditorOpen={showEditor && replyToId === post.id}
        onReplyClick={() => {
          // 简化切换逻辑，打开时必定触发聚焦
          setShowEditor((prev) => {
            if (!prev) {
              // 从隐藏变为显示 → 重置目标为主帖并触发聚焦
              setReplyToId(post.id);
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
            key={replyToId}
            parentId={replyToId}
            nodeId={post.node.nodeId}
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
              key={reply.id}
              item={reply}
              isRoot
              currentPostId={String(post.id)}
              backState={backState}
              onReplyClick={handleThreadReplyClick}
            />
          ))}
        </ul>
      </div>
    </>
  );
}