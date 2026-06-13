// src/app/threads/[slug]/PostFlat.tsx
'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import type { PostEntity } from '@/types/domain';
import PostFlatItem from '@/components/PostFlatItem';
import ReplyEditor from '@/components/features/ReplyEditor';
import { extractPostTitle } from '@/lib/utils/post';

interface PostFlatProps {
  posts: PostEntity[];
  totalReplyCount: number;
}

export function PostFlat({ posts, totalReplyCount }: PostFlatProps) {
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
      const width = el.offsetWidth;
      if (width > 0) setEditorWidth(width);
    };

    updateWidth();
    const observer = new ResizeObserver(updateWidth);
    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  const rootPost = posts.find((p) => p.isRoot) ?? posts[0];

  const activePost = activeEditorPostSlug === rootPost?.slug
    ? rootPost
    : posts.find((p) => p.slug === activeEditorPostSlug) ?? null;

  const replyLabel = activePost?.slug === rootPost?.slug
    ? `楼主 (${activePost.user?.nickname ?? activePost.user?.username ?? '匿名用户'})`
    : activePost?.user?.nickname ?? activePost?.user?.username ?? '匿名用户';

  const replyToTitle = activePost
    ? (extractPostTitle(activePost.content, { maxLength: 30 }) || '原帖内容')
    : '';

  return (
    <>
      <div ref={postListRef} className="post-list-container">
        {posts.length > 0 ? (
          posts.map((post) => {
            const isRoot = post.isRoot;
            const isEditorOpen = activeEditorPostSlug === post.slug;
            return (
              <div key={post.slug}>
                <PostFlatItem
                  post={post}
                  detailHref={`/threads/${post.slug}`}
                  replyCount={isRoot ? (post.commentCount ?? totalReplyCount) : 0}
                  isRoot={isRoot}
                  onReplyClick={() => toggleEditor(post.slug)}
                  isEditorOpen={isEditorOpen}
                />
              </div>
            );
          })
        ) : (
          <div className="empty-tip">暂无回复</div>
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
              <button
                onClick={closeEditor}
                className="fixed-reply-editor__close"
                aria-label="关闭回复框"
              >
                ✕
              </button>
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