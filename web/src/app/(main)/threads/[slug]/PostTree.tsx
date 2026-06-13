// src/app/threads/[slug]/PostTree.tsx
'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import type { PostEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';
import type { BackState } from '@/components/features/ThreadTree';
import ReplyEditor from '@/components/features/ReplyEditor';
import PostFlatItem from '@/components/PostFlatItem'; // ✨ 根帖直接复用平铺组件，保证视觉一致
import ThreadItem from '@/components/features/ThreadItem';

interface PostTreeProps {
  post: PostEntity;
  viewPosts: ThreadViewItem[];
  totalReplyCount: number;
  backState: BackState;
}

export function PostTree({ post, viewPosts, totalReplyCount, backState }: PostTreeProps) {
  // ✅ 以下状态与 Hooks 与 PostFlat.tsx 完全同构
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
    if (!activeEditorPostSlug) return;
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') { e.preventDefault(); e.stopPropagation(); closeEditor(); }
    };
    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [activeEditorPostSlug, closeEditor]);

  useEffect(() => {
    const el = postListRef.current;
    if (!el) return;
    const updateWidth = () => { const w = el.offsetWidth; if (w > 0) setEditorWidth(w); };
    updateWidth();
    const observer = new ResizeObserver(updateWidth);
    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  // ✅ 统一查找激活帖子 & 生成标签（兼容根帖与子回复）
  const activePost = activeEditorPostSlug === post.slug
    ? post
    : viewPosts.find((p) => p.slug === activeEditorPostSlug);

  const replyLabel = activePost?.slug === post.slug
    ? `楼主 (${activePost.user?.nickname ?? '匿名用户'})`
    : activePost?.user?.nickname ?? '匿名用户';

  return (
    <>
      {/* ✅ 容器 ref 与类名对齐 PostFlat */}
      <div ref={postListRef} className="tree-post-container">
        
        {/* ✨ 根帖：直接复用 PostFlatItem，确保两种视图下根帖像素级一致 */}
        <PostFlatItem
          post={post}
          detailHref={`/threads/${post.slug}`}
          replyCount={totalReplyCount}
          isRoot
          onReplyClick={() => toggleEditor(post.slug)}
          isEditorOpen={activeEditorPostSlug === post.slug}
        />

        {/* ✅ 子回复：保留原有 ThreadItem 组件 */}
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
                  // ✅ 回调签名对齐：只传 slug，由容器统一管理状态
                  onReplyClick={(slug) => toggleEditor(slug)}
                />
              ))}
            </ul>
          </div>
        )}
      </div>

      {/* ✅ Portal 编辑器结构与 PostFlat 完全一致 */}
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
                replyToTitle={activePost.title || post.title}
                replyToAuthor={activePost.user?.nickname ?? activePost.user?.username}
                autoFocus={shouldAutoFocus}
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