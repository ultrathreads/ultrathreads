'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import type { PostEntity } from '@/types/domain';
import PostCard from '@/components/features/PostCard';
import ReplyEditor from '@/components/features/ReplyEditor';
import { extractPostTitle } from '@/lib/utils/post';

interface PostFlatProps {
  posts: PostEntity[];
  totalReplyCount: number;
}

// ✨ 编辑器状态类型，区分回复和编辑模式
type EditorState = {
  postSlug: string;
  mode: 'reply' | 'edit';
} | null;

export function PostFlat({ posts, totalReplyCount }: PostFlatProps) {
  // ✨ 替换原有的 activeEditorPostSlug
  const [editorState, setEditorState] = useState<EditorState>(null);
  const [shouldAutoFocus, setShouldAutoFocus] = useState(false);
  const postListRef = useRef<HTMLDivElement>(null);
  const [editorWidth, setEditorWidth] = useState<number | undefined>(undefined);

  // ✨ 统一的编辑器切换方法
  const openEditor = useCallback((postSlug: string, mode: 'reply' | 'edit') => {
    setEditorState((prev) => {
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
      const width = el.offsetWidth;
      if (width > 0) setEditorWidth(width);
    };

    updateWidth();
    const observer = new ResizeObserver(updateWidth);
    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  const rootPost = posts.find((p) => p.isRoot) ?? posts[0];
  const activePost =
    editorState?.postSlug === rootPost?.slug
      ? rootPost
      : posts.find((p) => p.slug === editorState?.postSlug) ?? null;

  const replyToAuthor =
    activePost?.user?.nickname ?? activePost?.user?.username ?? '匿名用户';

  const replyToTitle = activePost
    ? extractPostTitle(activePost.content, { maxLength: 30 }) || '原帖内容'
    : '';

  // ✨ 判断当前某个帖子是否处于指定模式的编辑打开状态
  const isEditorOpen = (slug: string, mode: 'reply' | 'edit') =>
    editorState?.postSlug === slug && editorState?.mode === mode;

  return (
    <div ref={postListRef} className="post-list-container">
      {posts.length > 0 ? (
        posts.map((post) => (
          <div key={post.slug}>
            <PostCard
              post={post}
              detailHref={`/threads/${post.slug}`}
              replyCount={post.isRoot ? (post.commentCount ?? totalReplyCount) : 0}
              isRoot={post.isRoot}
              // ✨ 分别传递回复和编辑回调
              onReplyClick={() => openEditor(post.slug, 'reply')}
              onEditClick={() => openEditor(post.slug, 'edit')}
              // ✨ 分别传递两种模式的打开状态
              isReplyEditorOpen={isEditorOpen(post.slug, 'reply')}
              isEditEditorOpen={isEditorOpen(post.slug, 'edit')}
            />
          </div>
        ))
      ) : (
        <div className="empty-tip">暂无回复</div>
      )}

      {activePost && editorState && (
        <ReplyEditor
          // ✨ key 加入 mode，确保模式切换时编辑器重新挂载
          key={`${activePost.slug}-${editorState.mode}`}
          parentSlug={activePost.slug}
          postSlug={editorState.mode === 'edit' ? activePost.slug : undefined}
          // ✨ 传入模式和初始内容
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