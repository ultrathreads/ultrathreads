'use client';

import { useState, useRef, useEffect } from 'react';
import MDEditor from '@uiw/react-md-editor';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
// ✅ 1. 替换导入：createPost → createReply
import { createReply } from '@/services/post-service';

interface ReplyEditorProps {
  parentSlug: string;
  replyToTitle?: string;
  autoFocus?: boolean;
  onAutoFocusConsumed?: () => void;
}

export default function ReplyEditor({
  parentSlug,
  replyToTitle,
  autoFocus = false,
  onAutoFocusConsumed,
}: ReplyEditorProps) {
  const [content, setContent] = useState('');
  const router = useRouter();
  const editorWrapperRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!autoFocus) return;
    requestAnimationFrame(() => {
      editorWrapperRef.current?.scrollIntoView({ behavior: 'smooth', block: 'center' });
    });
    const timer = setTimeout(() => onAutoFocusConsumed?.(), 0);
    return () => clearTimeout(timer);
  }, [autoFocus, onAutoFocusConsumed]);

  const handleSubmit = async () => {
    const trimmed = content.trim();
    if (!trimmed) {
      toast.warning('回复内容不能为空');
      return;
    }

    // ✅ 2. 移除 extractPostTitle 相关逻辑（回复不需要 title）

    toast.promise(
      // ✅ 3. 调用新 API：parentSlug 作为第一个参数，body 只传 content
      createReply(parentSlug, { content: trimmed }),
      {
        loading: '发布中...',
        success: (result) => {
          setContent('');
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
        {replyToTitle && <span className="reply-to-tag">→ {replyToTitle}</span>}
      </h3>

      {/* ✅ 4. 移除隐藏的 parentSlug input（已由组件 props + API 路径承载） */}

      <div data-color-mode="light">
        <MDEditor
          value={content}
          onChange={(val) => setContent(val || '')}
          preview="live"
          height={200}
          visibleDragbar={false}
          textareaProps={{
            autoFocus,
            placeholder: replyToTitle
              ? `回复 @${replyToTitle}...`
              : '支持 Markdown 语法...',
          }}
        />
      </div>

      <div style={{ margin: '12px 0', textAlign: 'left' }}>
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
      </div>
    </div>
  );
}