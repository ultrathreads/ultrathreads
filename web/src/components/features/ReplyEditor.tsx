'use client';

import { useState, useRef, useEffect, useCallback } from 'react';
import MDEditor from '@uiw/react-md-editor';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { createReply } from '@/services/post-service';

interface ReplyEditorProps {
  parentSlug: string;
  replyToTitle?: string;
  replyToAuthor?: string;
  autoFocus?: boolean;
  onAutoFocusConsumed?: () => void;
  onClose?: () => void;   // ✅ 统一关闭回调（ESC + 点击域外 + 手动关闭）
  onSuccess?: () => void;
}

export default function ReplyEditor({
  parentSlug,
  replyToTitle,
  replyToAuthor,
  autoFocus = false,
  onAutoFocusConsumed,
  onClose,
}: ReplyEditorProps) {
  const [content, setContent] = useState('');
  const router = useRouter();
  const editorWrapperRef = useRef<HTMLDivElement>(null);

  // ✅ 拼接 placeholder
  const placeholder = replyToAuthor
    ? `@${replyToAuthor}${replyToTitle ? `：${replyToTitle}` : ''}...`
    : replyToTitle
      ? `回复：${replyToTitle}...`
      : '支持 Markdown 语法...';

  // ✅ 自动聚焦 & 滚动到可视区域
  useEffect(() => {
    if (!autoFocus) return;
    requestAnimationFrame(() => {
      editorWrapperRef.current?.scrollIntoView({ behavior: 'smooth', block: 'center' });
    });
    const timer = setTimeout(() => onAutoFocusConsumed?.(), 0);
    return () => clearTimeout(timer);
  }, [autoFocus, onAutoFocusConsumed]);

  // ✅ 统一关闭：ESC 键 + 点击域外区域
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
      // 仅当点击目标不在编辑器容器内时触发关闭
      if (editorWrapperRef.current && !editorWrapperRef.current.contains(target)) {
        onClose();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    // ✅ 使用 mousedown 而非 click，避免与按钮点击事件时序冲突
    document.addEventListener('mousedown', handleClickOutside);

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [onClose]);

  const handleSubmit = async () => {
    const trimmed = content.trim();
    if (!trimmed) {
      toast.warning('回复内容不能为空');
      return;
    }

    toast.promise(
      createReply(parentSlug, { content: trimmed }),
      {
        loading: '发布中...',
        success: (result) => {
          setContent('');
          // ✅ 提交成功后也调用 onClose 收起面板
          onClose?.();
          setTimeout(() => {
            router.push(result?.slug ? `/threads/${result.slug}` : '/');
            router.refresh();
          }, 600);
          return '回复发布成功 🎉';
        },
        error: (err) =>
          err instanceof Error ? err.message : '提交失败，请重试',
      }
    );
  };

  return (
    <div className="reply-editor-wrapper" style={{ marginTop: 24 }} ref={editorWrapperRef}>
      <h3 style={{ marginBottom: 12, display: 'flex', alignItems: 'center', gap: 8, flexWrap: 'wrap' }}>
        <span>✏️ 发表回复</span>
      </h3>

      <div data-color-mode="light">
        <MDEditor
          value={content}
          onChange={(val) => setContent(val || '')}
          preview="live"
          height={200}
          visibleDragbar={false}
          textareaProps={{
            autoFocus,
            placeholder,
          }}
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
        {/* ✅ 显式取消按钮，提升可访问性 */}
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
  );
}