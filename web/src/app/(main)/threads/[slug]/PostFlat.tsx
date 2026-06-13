// src/app/threads/[slug]/PostFlat.tsx
'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import type { PostEntity } from '@/types/domain';
import PostFlatItem from '@/components/PostFlatItem';
import ReplyEditor from '@/components/features/ReplyEditor';
// ✅ 已移除: import styles from './FixedReplyEditor.module.css';

interface PostFlatProps {
  posts: PostEntity[];
  totalReplyCount: number;
}

export function PostFlat({ posts, totalReplyCount }: PostFlatProps) {
  const [activeEditorPostSlug, setActiveEditorPostSlug] = useState<string | null>(null);
  const [mounted, setMounted] = useState(false);
  const postListRef = useRef<HTMLDivElement>(null);
  const [editorWidth, setEditorWidth] = useState<number | undefined>(undefined);

  const toggleEditor = useCallback((postSlug: string) => {
    setActiveEditorPostSlug((prev) => (prev === postSlug ? null : postSlug));
  }, []);

  const closeEditor = useCallback(() => {
    setActiveEditorPostSlug(null);
  }, []);

  // SSR 安全守卫
  useEffect(() => {
    setMounted(true);
  }, []);

  // ESC 键关闭
  useEffect(() => {
    if (!activeEditorPostSlug) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        e.preventDefault();
        e.stopPropagation();
        closeEditor();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [activeEditorPostSlug, closeEditor]);

  // 宽度同步（仅增强样式，不阻塞渲染）
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

  const activePost = activeEditorPostSlug
    ? posts.find((p) => p.slug === activeEditorPostSlug)
    : null;

  return (
    <>
      {/* ✅ 帖子列表：保持原始类名不变，ref 仅用于测量 */}
      <div ref={postListRef} className="your-original-post-list-class-name">
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

      {/* ✅ 吸底编辑器：使用全局类名替代 CSS Module */}
      {mounted && activePost && createPortal(
        <div
          className="fixed-reply-editor"
          style={editorWidth ? { width: editorWidth } : undefined}
          role="dialog"
          aria-modal="true"
          aria-label={`回复 ${activePost.author?.name ?? '匿名用户'}`}
        >
          <div className="fixed-reply-editor__inner">
            <div className="fixed-reply-editor__header">
              <span className="fixed-reply-editor__label">
                回复 <span className="fixed-reply-editor__author">{activePost.author?.name ?? '匿名用户'}</span>
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
                parentSlug={activePost.slug}
                replyToTitle={activePost.title}
                autoFocus
                onSuccess={closeEditor}
              />
            </div>
          </div>
        </div>,
        document.body
      )}
    </>
  );
}