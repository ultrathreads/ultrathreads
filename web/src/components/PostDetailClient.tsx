// components/PostDetailClient.tsx
'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
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
  const [replyToSlug, setReplyToSlug] = useState<string>(post.slug);
  const [replyToTitle, setReplyToTitle] = useState<string>(post.title);
  const [shouldAutoFocus, setShouldAutoFocus] = useState(false);
  const [mounted, setMounted] = useState(false);

  // 用于测量帖子内容区宽度，使吸底编辑器与内容区等宽
  const contentRef = useRef<HTMLDivElement>(null);
  const [editorWidth, setEditorWidth] = useState<number | undefined>(undefined);

  // SSR 安全守卫
  useEffect(() => {
    setMounted(true);
  }, []);

  // ESC 键关闭编辑器
  useEffect(() => {
    if (!showEditor) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        e.preventDefault();
        e.stopPropagation();
        setShowEditor(false);
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [showEditor]);

  // 宽度同步：监听内容区宽度变化
  useEffect(() => {
    const el = contentRef.current;
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

  const handleThreadReplyClick = useCallback(
    (targetSlug: string, targetTitle: string) => {
      setReplyToSlug(targetSlug);
      setReplyToTitle(`${post.title}(${post.user.nickname})`);
      setShowEditor(true);
      setShouldAutoFocus(true);
    },
    [post.title, post.user.nickname],
  );

  const closeEditor = useCallback(() => {
    setShowEditor(false);
  }, []);

  return (
    <>
      {/* ref 绑定到包裹整个帖子内容的容器，用于宽度测量 */}
      <div ref={contentRef}>
        <PostDetailCard
          post={post}
          replyCount={totalReplyCount}
          isEditorOpen={showEditor && replyToSlug === post.slug}
          onReplyClick={() => {
            setShowEditor((prev) => {
              if (!prev) {
                setReplyToSlug(post.slug);
                setReplyToTitle(post.title);
                setShouldAutoFocus(true);
              }
              return !prev;
            });
          }}
        />

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
      </div>

      {/* 吸底编辑器：与 PostFlat 完全一致的 Portal + 固定定位 */}
      {mounted && showEditor && createPortal(
        <div
          className="fixed-reply-editor"
          style={editorWidth ? { width: editorWidth } : undefined}
          role="dialog"
          aria-modal="true"
          aria-label={`回复 ${replyToSlug === post.slug ? '楼主' : replyToTitle}`}
        >
          <div className="fixed-reply-editor__inner">
            <div className="fixed-reply-editor__header">
              <span className="fixed-reply-editor__label">
                回复{' '}
                <span className="fixed-reply-editor__author">
                  {replyToSlug === post.slug
                    ? post.user?.username ?? '楼主'
                    : replyToTitle}
                </span>
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
                key={replyToSlug}
                parentSlug={replyToSlug}
                replyToTitle={replyToTitle}
                autoFocus={shouldAutoFocus}
                onSuccess={closeEditor}
                onAutoFocusConsumed={() => setShouldAutoFocus(false)}
              />
            </div>
          </div>
        </div>,
        document.body,
      )}
    </>
  );
}