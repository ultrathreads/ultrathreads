'use client';

import { useState, useRef, useEffect } from 'react';
import MDEditor from '@uiw/react-md-editor';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { createPost } from '@/services/post-service';
import { extractPostTitle } from '@/lib/utils/post';

interface ReplyEditorProps {
  parentId: number;
  nodeId: number;
  replyToTitle?: string;
  autoFocus?: boolean;
  onAutoFocusConsumed?: () => void;
}

export default function ReplyEditor({
  parentId,
  nodeId,
  replyToTitle,
  autoFocus = false,
  onAutoFocusConsumed,
}: ReplyEditorProps) {
  const [content, setContent] = useState('');
  const router = useRouter();
  const editorWrapperRef = useRef<HTMLDivElement>(null);

  // ✅ 仅负责滚动定位，聚焦交给 MDEditor 自身处理
  useEffect(() => {
    if (!autoFocus) return;

    requestAnimationFrame(() => {
      editorWrapperRef.current?.scrollIntoView({ behavior: 'smooth', block: 'center' });
    });

    // ✅ MDEditor 通过 textareaProps.autoFocus 自行处理聚焦后，通知父组件重置
    // 延迟一帧确保 onAutoFocusConsumed 不会在当前渲染周期内触发状态变更警告
    const timer = setTimeout(() => {
      onAutoFocusConsumed?.();
    }, 0);

    return () => clearTimeout(timer);
  }, [autoFocus, onAutoFocusConsumed]);

  const handleSubmit = async () => {
    const trimmed = content.trim();
    if (!trimmed) { toast.warning('回复内容不能为空'); return; }

    const title = extractPostTitle(trimmed, { maxLength: 30 });
    if (!title) { toast.warning('无法从内容中提取标题，请输入有效文本'); return; }

    toast.promise(
      createPost({ title, nodeId, parentId, content }),
      {
        loading: '发布中...',
        success: (result) => {
          setContent('');
          setTimeout(() => {
            router.push(result?.id ? `/post/${result.id}` : '/');
            router.refresh();
          }, 600);
          return '回复发布成功 🎉';
        },
        error: (err) => err instanceof Error ? err.message : '提交失败，请重试',
      }
    );
  };

  return (
    <div className="reply-editor-wrapper" style={{ marginTop: 24 }} ref={editorWrapperRef}>
      <h3 style={{ marginBottom: 12, display: 'flex', alignItems: 'center', gap: 8, flexWrap: 'wrap' }}>
        <span>✏️ 发表回复</span>
        {replyToTitle && <span className="reply-to-tag">→ {replyToTitle}</span>}
      </h3>

      <input type="hidden" name="parentId" value={parentId} />
      <input type="hidden" name="nodeId" value={nodeId} />

      <div data-color-mode="light">
        {/* ✅ 核心修复：通过 textareaProps 透传 autoFocus，由 MDEditor 内部在正确时机执行 */}
        <MDEditor
          value={content}
          onChange={(val) => setContent(val || '')}
          preview="live"
          height={200}
          visibleDragbar={false}
          textareaProps={{
            autoFocus,
            placeholder: replyToTitle ? `回复 @${replyToTitle}...` : '支持 Markdown 语法...',
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