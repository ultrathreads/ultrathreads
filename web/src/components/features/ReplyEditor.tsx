// src/app/threads/[id]/ReplyEditor.tsx
'use client';

import { useState, useRef, useEffect } from 'react';
import { createPortal } from 'react-dom';
import MDEditor from '@uiw/react-md-editor';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { createReply } from '@/services/post-service';

interface ReplyEditorProps {
  parentSlug: string;
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
  replyToTitle,
  replyToAuthor,
  autoFocus = false,
  containerWidth,
  onAutoFocusConsumed,
  onClose,
  onSuccess,
}: ReplyEditorProps) {
  const [content, setContent] = useState('');
  const [mounted, setMounted] = useState(false);
  const router = useRouter();
  const editorWrapperRef = useRef<HTMLDivElement>(null);

  useEffect(() => setMounted(true), []);

  const placeholder = replyToAuthor
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
      if (e.key === 'Escape') {
        e.preventDefault();
        e.stopPropagation();
        onClose();
      }
    };

    const handleClickOutside = (e: MouseEvent) => {
      const target = e.target as Node;
      if (editorWrapperRef.current && !editorWrapperRef.current.contains(target)) {
        onClose();
      }
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
      toast.warning('回复内容不能为空');
      return;
    }

    toast.promise(createReply(parentSlug, { content: trimmed }), {
      loading: '发布中...',
      success: (result) => {
        setContent('');
        onClose?.();
        onSuccess?.();
        setTimeout(() => {
          router.push(result?.slug ? `/threads/${result.slug}` : '/');
          router.refresh();
        }, 600);
        return '回复发布成功 🎉';
      },
      error: (err) => (err instanceof Error ? err.message : '提交失败，请重试'),
    });
  };

  const editorContent = (
    <div
      className="fixed-reply-editor"
      style={containerWidth ? { width: containerWidth } : undefined}
      role="dialog"
      aria-modal="true"
      aria-label={`回复 ${replyToAuthor ?? ''}`}
    >
      <div className="fixed-reply-editor__inner" ref={editorWrapperRef}>
        {/* ✅ 直接使用 replyToAuthor */}
        <div className="fixed-reply-editor__header">
          <span className="fixed-reply-editor__label">
            回复 <span className="fixed-reply-editor__author">{replyToAuthor}</span>
          </span>
          {onClose && (
            <button onClick={onClose} className="fixed-reply-editor__close" aria-label="关闭回复框">
              ✕
            </button>
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
                backgroundColor: !content.trim() ? '#a0aec0' : '#3182ce',
                color: '#fff',
                border: 'none',
                borderRadius: 6,
                cursor: !content.trim() ? 'not-allowed' : 'pointer',
                fontSize: 14,
              }}
            >
              发布回复
            </button>
            {onClose && (
              <button
                onClick={onClose}
                style={{
                  padding: '8px 24px',
                  backgroundColor: 'transparent',
                  color: '#718096',
                  border: '1px solid #e2e8f0',
                  borderRadius: 6,
                  cursor: 'pointer',
                  fontSize: 14,
                }}
              >
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