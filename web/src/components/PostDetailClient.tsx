'use client';

import { useState } from 'react';
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
  const [showEditor, setShowEditor] = useState(false);

  return (
    <>
      {/* ✅ 将 showEditor 传递给 PostDetailCard 驱动按钮样式 */}
      <PostDetailCard
        post={post}
        replyCount={totalReplyCount}
        isEditorOpen={showEditor}
        onReplyClick={() => setShowEditor((prev) => !prev)}
      />

      {showEditor && (
        <div id="reply-editor">
          <ReplyEditor parentId={post.id} nodeId={post.node.nodeId} />
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
            />
          ))}
        </ul>
      </div>
    </>
  );
}