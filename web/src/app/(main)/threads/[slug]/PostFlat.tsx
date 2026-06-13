// src/app/threads/[slug]/PostFlat.tsx
'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import type { PostEntity } from '@/types/domain';
import PostFlatItem from '@/components/PostFlatItem';
import ReplyEditor from '@/components/features/ReplyEditor';

interface PostFlatProps {
  posts: PostEntity[];
  totalReplyCount: number;
}

/**
 * ✅ 新增：从正文内容中截取纯文本作为回复引用标题
 * 去除 HTML/Markdown 标签、多余空白，并限制最大长度
 */
function getReplyLabelFromContent(content: string | undefined, maxLength = 30): string {
  if (!content) return '原帖内容';

  const plainText = content
    .replace(/<[^>]*>/g, '')           // 去除 HTML 标签
    .replace(/!\[.*?\]\(.*?\)/g, '[图片]') // 将 Markdown 图片替换为占位符
    .replace(/\[([^\]]*)\]\(.*?\)/g, '$1') // 提取 Markdown 链接文本
    .replace(/[#*_~`>]/g, '')          // 去除常见 Markdown 语法符号
    .replace(/\s+/g, ' ')              // 合并连续空白
    .trim();

  if (!plainText) return '原帖内容';
  return plainText.length > maxLength
    ? `${plainText.slice(0, maxLength)}...`
    : plainText;
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

  const closeEditor = useCallback(() => {
    setActiveEditorPostSlug(null);
  }, []);

  useEffect(() => {
    setMounted(true);
  }, []);

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

  // ✅ 修复：平铺模式下从内容截取，而非使用可能为空的 title
  const replyToTitle = activePost
    ? getReplyLabelFromContent(activePost.content)
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