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
  // 根据模式决定初始内容
  const [content, setContent] = useState(mode === 'edit' ? (initialContent ?? '') : '');
  const [mounted, setMounted] = useState(false);
  const router = useRouter();
  const editorWrapperRef = useRef<HTMLDivElement>(null);

  useEffect(() => setMounted(true), []);

  // 当 mode/initialContent 变化时同步内容（防御性处理）
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

  // ... autoFocus 和 clickOutside/Escape 的 useEffect 保持不变 ...
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

  // 统一提交逻辑：根据 mode 分发
  const handleSubmit = async () => {
    const trimmed = content.trim();
    if (!trimmed) {
      toast.warning(isEdit ? '内容不能为空' : '回复内容不能为空');
      return;
    }

    const actionFn = isEdit
      ? () => updateReply(postSlug!, { content: trimmed })   // ✨ 用 postSlug
      : () => createReply(parentSlug, { content: trimmed }); // ✨ 用 parentSlug

    const loadingText = isEdit ? '保存中...' : '发布中...';
    const successText = isEdit ? '修改成功 ✅' : '回复发布成功 🎉';

    toast.promise(actionFn(), {
      loading: loadingText,
      success: () => {
        if (!isEdit) setContent('');
        onClose?.();
        onSuccess?.();

        // 编辑成功后：URL 追加 refresh=1 并强制刷新
        if (isEdit) {
          const url = new URL(window.location.href);
          url.searchParams.set('refresh', '1');
          router.replace(url.pathname + url.search);
        } else {
          // 回复模式保持原有逻辑
          router.refresh();
        }

        return successText;
      },
      error: (err) => (err instanceof Error ? err.message : '操作失败，请重试'),
    });
  };

  // ✨ UI 文案根据模式切换
  const headerLabel = isEdit
    ? <span className="fixed-reply-editor__label">✏️ 编辑帖子</span>
    : (
      <span className="fixed-reply-editor__label">
        回复 <span className="fixed-reply-editor__author">{replyToAuthor}</span>
      </span>
    );

  const submitButtonText = isEdit ? '保存修改' : '发布回复';

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
              height={200}
              visibleDragbar={false}
              textareaProps={{ autoFocus, placeholder }}
            />
          </div>

          <div style={{ margin: '12px 0', textAlign: 'left', display: 'flex', gap: 8 }}>
            <button
              onClick={handleSubmit}
              disabled={!content.trim()}
              style={{
                padding: '8px 24px',
                backgroundColor: !content.trim() ? '#a0aec0' : isEdit ? '#38a169' : '#3182ce',
                color: '#fff', border: 'none', borderRadius: 6,
                cursor: !content.trim() ? 'not-allowed' : 'pointer', fontSize: 14,
              }}
            >
              {submitButtonText}
            </button>
            {onClose && (
              <button onClick={onClose} style={{
                padding: '8px 24px', backgroundColor: 'transparent',
                color: '#718096', border: '1px solid #e2e8f0', borderRadius: 6,
                cursor: 'pointer', fontSize: 14,
              }}>
                取消
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );

  if (!mounted) return null;
  return createPortal(editorContent, document.body);
}