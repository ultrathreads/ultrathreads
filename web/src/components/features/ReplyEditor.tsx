'use client';

import { useState, useRef, useEffect } from 'react';
import { createPortal } from 'react-dom';
import MDEditor from '@uiw/react-md-editor';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { createReply, updateReply } from '@/services/post-service';

interface ReplyEditorProps {
  parentSlug: string;
  postSlug?: string;
  mode?: 'reply' | 'edit';
  initialContent?: string;
  replyToTitle?: string;
  replyToAuthor?: string;
  autoFocus?: boolean;
  containerWidth?: number;
  onAutoFocusConsumed?: () => void;
  onClose?: () => void;
  onSuccess?: () => void;
}

function buildPostDetailUrl(slug: string): string {
  const url = new URL(`/threads/${slug}`, window.location.origin);
  url.searchParams.set('refresh', '1');
  return url.pathname + url.search;
}

export default function ReplyEditor({
  parentSlug,
  postSlug,
  mode = 'reply',
  initialContent,
  replyToTitle,
  replyToAuthor,
  autoFocus = false,
  containerWidth,
  onAutoFocusConsumed,
  onClose,
  onSuccess,
}: ReplyEditorProps) {
  const [content, setContent] = useState(mode === 'edit' ? (initialContent ?? '') : '');
  const [mounted, setMounted] = useState(false);
  const router = useRouter();
  const editorWrapperRef = useRef<HTMLDivElement>(null);
  const submittingRef = useRef(false);

  useEffect(() => setMounted(true), []);

  useEffect(() => {
    if (mode === 'edit') {
      setContent(initialContent ?? '');
    } else {
      setContent('');
    }
  }, [mode, initialContent]);

  const isEdit = mode === 'edit';

  const placeholder = isEdit
    ? '编辑内容...支持 Markdown 语法'
    : replyToAuthor
      ? `@${replyToAuthor}${replyToTitle ? `：${replyToTitle}` : ''}...`
      : replyToTitle
        ? `回复：${replyToTitle}...`
        : '支持 Markdown 语法...';

  useEffect(() => {
    if (!autoFocus) return;
    requestAnimationFrame(() => {
      editorWrapperRef.current?.scrollIntoView({ behavior: 'smooth', block: 'center' });
    });
    const timer = setTimeout(() => onAutoFocusConsumed?.(), 0);
    return () => clearTimeout(timer);
  }, [autoFocus, onAutoFocusConsumed]);

  useEffect(() => {
    if (!onClose) return;
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') { e.preventDefault(); e.stopPropagation(); onClose(); }
    };
    const handleClickOutside = (e: MouseEvent) => {
      if (editorWrapperRef.current && !editorWrapperRef.current.contains(e.target as Node)) onClose();
    };
    document.addEventListener('keydown', handleKeyDown);
    document.addEventListener('click', handleClickOutside);
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      document.removeEventListener('click', handleClickOutside);
    };
  }, [onClose]);

  const handleSubmit = async () => {
    const trimmed = content.trim();
    if (!trimmed) {
      toast.warning(isEdit ? '内容不能为空' : '回复内容不能为空');
      return;
    }
    if (submittingRef.current) return;
    submittingRef.current = true;

    try {
      // ✨ 参照 PostForm：直接将 promise 传入 toast.promise，让 success 自动推导 result 类型
      await toast.promise(
        isEdit
          ? updateReply(postSlug!, { content: trimmed })
          : createReply(parentSlug, { content: trimmed }),
        {
          loading: isEdit ? '保存中...' : '发布中...',
          success: (result) => {
            submittingRef.current = false;
            if (!isEdit) setContent('');
            onClose?.();
            onSuccess?.();

            // 编辑模式：跳转到被编辑内容所属帖子
            // 回复模式：跳转到新创建回复的 slug，兜底使用 parentSlug
            const targetSlug = isEdit
              ? (postSlug ?? parentSlug)
              : ((result as { slug?: string })?.slug ?? parentSlug);

            setTimeout(() => {
              router.push(buildPostDetailUrl(targetSlug));
            }, 600);

            return isEdit ? '修改成功 ✅' : '回复发布成功 🎉';
          },
          error: (err) => {
            submittingRef.current = false;
            return err instanceof Error ? err.message : '操作失败，请重试';
          },
        }
      );
    } catch {
      submittingRef.current = false;
    }
  };

  const headerLabel = isEdit
    ? <span className="fixed-reply-editor__label">✏️ 编辑帖子</span>
    : (
      <span className="fixed-reply-editor__label">
        回复 <span className="fixed-reply-editor__author">{replyToAuthor}</span>
      </span>
    );

  const submitButtonText = isEdit ? '保存修改' : '发布回复';
  const isSubmitDisabled = !content.trim() || submittingRef.current;

  const editorContent = (
    <div
      className="fixed-reply-editor"
      style={containerWidth ? { width: containerWidth } : undefined}
      role="dialog"
      aria-modal="true"
      aria-label={isEdit ? '编辑帖子' : `回复 ${replyToAuthor ?? ''}`}
    >
      <div className="fixed-reply-editor__inner" ref={editorWrapperRef}>
        <div className="fixed-reply-editor__header">
          {headerLabel}
          {onClose && (
            <button onClick={onClose} className="fixed-reply-editor__close" aria-label="关闭">✕</button>
          )}
        </div>

        <div className="fixed-reply-editor__body">
          <div data-color-mode="light">
            <MDEditor
              value={content}
              onChange={(val) => setContent(val || '')}
              preview="live"
              height={300}
              visibleDragbar={false}
              textareaProps={{ autoFocus, placeholder }}
            />
          </div>

          <div style={{ margin: '12px 0', textAlign: 'left', display: 'flex', gap: 8 }}>
            {onClose && (
              <button onClick={onClose} style={{
                padding: '8px 24px', backgroundColor: 'transparent',
                color: '#718096', border: '1px solid #e2e8f0', borderRadius: 6,
                cursor: 'pointer', fontSize: 14,
              }}>
                取消
              </button>
            )}
            <button
              onClick={handleSubmit}
              disabled={isSubmitDisabled}
              style={{
                padding: '8px 24px',
                backgroundColor: isSubmitDisabled ? '#a0aec0' : isEdit ? '#38a169' : '#3182ce',
                color: '#fff', border: 'none', borderRadius: 6,
                cursor: isSubmitDisabled ? 'not-allowed' : 'pointer', fontSize: 14,
              }}
            >
              {submitButtonText}
            </button>
          </div>
        </div>
      </div>
    </div>
  );

  if (!mounted) return null;
  return createPortal(editorContent, document.body);
}