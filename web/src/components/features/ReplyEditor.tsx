'use client';

import { useState } from 'react';
import MDEditor from '@uiw/react-md-editor';
import { useRouter } from 'next/navigation';
import { createPost } from '@/services/post-service';

interface ReplyEditorProps {
  parentId: number;
  nodeId: number;
}

/**
 * 从 Markdown 内容中提取纯文本标题
 * - 去除所有 Markdown 语法标记
 * - 取第一行有效文本
 * - 最多截取 20 个字符
 */
function extractTitle(content: string): string {
  const plainText = content
    .replace(/#{1,6}\s*/g, '')       // 去除标题标记 # ## ###
    .replace(/[*_~`]/g, '')          // 去除加粗/斜体/删除线/行内代码
    .replace(/\[([^\]]*)\]\([^)]*\)/g, '$1') // 链接取文本部分
    .replace(/!\[[^\]]*\]\([^)]*\)/g, '')    // 去除图片
    .replace(/>\s*/g, '')            // 去除引用标记
    .replace(/[-+*]\s+/g, '')        // 去除无序列表标记
    .replace(/\d+\.\s+/g, '')        // 去除有序列表标记
    .trim();

  const firstLine = plainText.split('\n').find(line => line.trim().length > 0) || '';
  return firstLine.slice(0, 20);
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

    const title = extractTitle(trimmed);
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