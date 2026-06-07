'use client';

import { useState } from 'react';
import MDEditor from '@uiw/react-md-editor';
import { useRouter } from 'next/navigation';
import { createPost } from '@/services/post-service';
import { extractPostTitle } from '@/lib/utils/post';

interface ReplyEditorProps {
  parentId: number;
  nodeId: number;
}

export default function ReplyEditor({ parentId, nodeId }: ReplyEditorProps) {
  const [content, setContent] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const handleSubmit = async () => {
    const trimmed = content.trim();
    if (!trimmed) {
      setError('回复内容不能为空');
      return;
    }

    const title = extractPostTitle(trimmed, { maxLength: 30 });
    if (!title) {
      setError('无法从内容中提取标题，请输入有效文本');
      return;
    }

    setSubmitting(true);
    setError(null);

    try {

      const result = await createPost({
        title,
        nodeId,
        parentId,
        content,
      });

      setContent('');
      // ✅ 跳转到新帖子详情页，若无返回则跳首页
      router.push(result?.id ? `/post/${result.id}` : '/');   
      router.refresh(); // 仅触发当前页面的服务端重新渲染
    } catch (err) {
      setError(err instanceof Error ? err.message : '提交失败，请重试');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="reply-editor-wrapper" style={{ marginTop: 24 }}>
      <h3 style={{ marginBottom: 12 }}>✏️ 发表回复</h3>

      {/* 隐藏属性 */}
      <input type="hidden" name="parentId" value={parentId} />
      <input type="hidden" name="nodeId" value={nodeId} />

      <div data-color-mode="light">
        <MDEditor
          value={content}
          onChange={(val) => setContent(val || '')}
          preview="live"
          height={200}
          visibleDragbar={false}
        />
      </div>

      {error && <p style={{ color: '#e53e3e', fontSize: 14, marginTop: 8 }}>{error}</p>}

      <div style={{ margin: '12px 0', textAlign: 'right' }}>
        <button
          onClick={handleSubmit}
          disabled={submitting || !content.trim()}
          style={{
            padding: '8px 24px',
            backgroundColor: submitting ? '#a0aec0' : '#3182ce',
            color: '#fff',
            border: 'none',
            borderRadius: 6,
            cursor: submitting ? 'not-allowed' : 'pointer',
            fontSize: 14,
          }}
        >
          {submitting ? '提交中...' : '发布回复'}
        </button>
      </div>
    </div>
  );
}