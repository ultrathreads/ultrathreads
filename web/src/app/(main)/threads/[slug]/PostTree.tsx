// src/app/threads/[slug]/PostTree.tsx
'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import type { PostEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';
import type { BackState } from '@/components/features/ThreadTree';
import ReplyEditor from '@/components/features/ReplyEditor';
import PostFlatItem from '@/components/PostFlatItem';
import ThreadItem from '@/components/features/ThreadItem';
import { extractPostTitle } from '@/lib/utils/post';

interface PostTreeProps {
  post: PostEntity;
  viewPosts: ThreadViewItem[];
  totalReplyCount: number;
  backState: BackState;
}

export function PostTree({ post, viewPosts, totalReplyCount, backState }: PostTreeProps) {
  const [activeEditorPostSlug, setActiveEditorPostSlug] = useState<string | null>(null);
  const [shouldAutoFocus, setShouldAutoFocus] = useState(false);
  const [mounted, setMounted] = useState(false);
  const postListRef = useRef<HTMLDivElement>(null);
  const [editorWidth, setEditorWidth] = useState<number | undefined>(undefined);

  const toggleEditor = useCallback((postSlug: string) => {
    setActiveEditorPostSlug((prev) => {
      if (prev !== postSlug) setShouldAutoFocus(true);
      return prev === postSlug ? null : postSlug;
    });
  }, []);

  const closeEditor = useCallback(() => setActiveEditorPostSlug(null), []);

  useEffect(() => setMounted(true), []);

  useEffect(() => {
    const el = postListRef.current;
    if (!el) return;
    const updateWidth = () => {
      const w = el.offsetWidth;
      if (w > 0) setEditorWidth(w);
    };
    updateWidth();
    const observer = new ResizeObserver(updateWidth);
    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  const activePost = activeEditorPostSlug === post.slug
    ? post
    : viewPosts.find((p) => p.slug === activeEditorPostSlug);

  const replyLabel = activePost?.slug === post.slug
    ? `楼主 (${activePost.user?.nickname ?? activePost.user?.username ?? '匿名用户'})`
    : activePost?.user?.nickname ?? activePost?.user?.username ?? '匿名用户';

  const replyToTitle = activePost
    ? (extractPostTitle(activePost.content, { maxLength: 30 }) || '原帖内容')
    : '';

  return (
    <>
      <div ref={postListRef} className="tree-post-container">
        <PostFlatItem
          post={post}
          detailHref={`/threads/${post.slug}`}
          replyCount={totalReplyCount}
          isRoot
          onReplyClick={() => toggleEditor(post.slug)}
          isEditorOpen={activeEditorPostSlug === post.slug}
        />

        {viewPosts.length > 0 && (
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
                  currentPostSlug={post.slug}
                  backState={backState}
                  onReplyClick={(slug) => toggleEditor(slug)}
                />
              ))}
            </ul>
          </div>
        )}
      </div>

      {mounted && activePost && createPortal(
        <div
          className="fixed-reply-editor"
          style={editorWidth ? { width: editorWidth } : undefined}
          role="dialog"
          aria-modal="true"
          aria-label={`回复 ${replyLabel}`}
        >
          <div className="fixed-reply-editor__inner">
            <div className="fixed-reply-editor__header">
              <span className="fixed-reply-editor__label">
                回复 <span className="fixed-reply-editor__author">{replyLabel}</span>
              </span>
              <button onClick={closeEditor} className="fixed-reply-editor__close" aria-label="关闭回复框">✕</button>
            </div>
            <div className="fixed-reply-editor__body">
              <ReplyEditor
                key={activePost.slug}
                parentSlug={activePost.slug}
                replyToTitle={replyToTitle}
                replyToAuthor={activePost.user?.nickname ?? activePost.user?.username}
                autoFocus={shouldAutoFocus}
                onClose={closeEditor}
                onSuccess={closeEditor}
                onAutoFocusConsumed={() => setShouldAutoFocus(false)}
              />
            </div>
          </div>
        </div>,
        document.body
      )}
    </>
  );
}