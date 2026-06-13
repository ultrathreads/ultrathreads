// src/app/threads/[slug]/PostTree.tsx
'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import type { PostEntity } from '@/types/domain';
import type { ThreadViewItem, BackState } from '@/types/view';
import ReplyEditor from '@/components/features/ReplyEditor';
import PostCard from '@/components/features/PostCard';
import ThreadItem from '@/components/features/ThreadItem';
import { extractPostTitle } from '@/lib/utils/post';

interface PostTreeProps {
  post: PostEntity;
  viewPosts: ThreadViewItem[];
  totalReplyCount: number;
  backState: BackState;
}

// ✨ 新增：编辑器状态类型，区分回复和编辑模式
type EditorState = {
  postSlug: string;
  mode: 'reply' | 'edit';
} | null;

export function PostTree({ post, viewPosts, totalReplyCount, backState }: PostTreeProps) {
  // ✨ 替换原有的 activeEditorPostSlug
  const [editorState, setEditorState] = useState<EditorState>(null);
  const [shouldAutoFocus, setShouldAutoFocus] = useState(false);
  const postListRef = useRef<HTMLDivElement>(null);
  const [editorWidth, setEditorWidth] = useState<number | undefined>(undefined);

  // ✨ 统一的编辑器切换方法
  const openEditor = useCallback((postSlug: string, mode: 'reply' | 'edit') => {
    setEditorState((prev) => {
      // 点击同一个按钮时关闭
      if (prev?.postSlug === postSlug && prev?.mode === mode) return null;
      setShouldAutoFocus(true);
      return { postSlug, mode };
    });
  }, []);

  const closeEditor = useCallback(() => setEditorState(null), []);

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

  const activePost = editorState?.postSlug === post.slug
    ? post
    : viewPosts.find((p) => p.slug === editorState?.postSlug);

  const replyToAuthor = activePost?.user?.nickname ?? activePost?.user?.username ?? '匿名用户';

  const replyToTitle = activePost
    ? (extractPostTitle(activePost.content, { maxLength: 30 }) || '原帖内容')
    : '';

  // ✨ 判断当前某个帖子是否处于指定模式的编辑打开状态
  const isEditorOpen = (slug: string, mode: 'reply' | 'edit') =>
    editorState?.postSlug === slug && editorState?.mode === mode;

  return (
    <div ref={postListRef} className="tree-post-container">
      <PostCard
        post={post}
        detailHref={`/threads/${post.slug}`}
        replyCount={totalReplyCount}
        isRoot={post.isRoot}
        onEditClick={() => openEditor(post.slug, 'edit')}
        onReplyClick={() => openEditor(post.slug, 'reply')}
        isReplyEditorOpen={isEditorOpen(post.slug, 'reply')}
        isEditEditorOpen={isEditorOpen(post.slug, 'edit')}
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
                isRoot={post.isRoot}
                currentPostSlug={post.slug}
                backState={backState}
                onReplyClick={(slug) => openEditor(slug, 'reply')}
                // ✨ 透传编辑回调
                onEditClick={(slug) => openEditor(slug, 'edit')}
                // ✨ 分别传递两种模式的打开状态
                isReplyEditorOpen={isEditorOpen(reply.slug, 'reply')}
                isEditEditorOpen={isEditorOpen(reply.slug, 'edit')}
              />
            ))}
          </ul>
        </div>
      )}

      {activePost && editorState && (
        <ReplyEditor
          key={`${activePost.slug}-${editorState.mode}`}
          parentSlug={activePost.slug}
          // ✨ 编辑模式下传入初始内容和模式标识
          postSlug={editorState.mode === 'edit' ? activePost.slug : undefined}
          mode={editorState.mode}
          initialContent={editorState.mode === 'edit' ? activePost.rawContent : undefined}
          replyToTitle={replyToTitle}
          replyToAuthor={replyToAuthor}
          containerWidth={editorWidth}
          autoFocus={shouldAutoFocus}
          onClose={closeEditor}
          onSuccess={closeEditor}
          onAutoFocusConsumed={() => setShouldAutoFocus(false)}
        />
      )}
    </div>
  );
}