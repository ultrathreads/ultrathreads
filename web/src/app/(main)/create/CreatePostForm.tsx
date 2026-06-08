// src/app/create/CreatePostForm.tsx
'use client';

import { useState, useRef, FormEvent } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import MDEditor from '@uiw/react-md-editor';
import { toast } from 'sonner';
import { createPost } from '@/services/post-service';
import { TagInput } from '@/components/ui/TagInput';
import type { NodeEntity } from '@/types/domain';

interface CreatePostFormProps {
  nodes: NodeEntity[];
}

export function CreatePostForm({ nodes }: CreatePostFormProps) {
  const router = useRouter();
  const searchParams = useSearchParams();

  // 从 URL 读取 nodeId 作为初始值，并确保该 nodeId 在 nodes 列表中存在
  const initialNodeId = (() => {
    const id = searchParams.get('nodeId');
    if (id && nodes.some((n) => String(n.nodeId) === id)) {
      return id;
    }
    return '';
  })();

  const [title, setTitle] = useState('');
  const [nodeId, setNodeId] = useState(initialNodeId);
  const [tags, setTags] = useState('');
  const [content, setContent] = useState('');
  const submittingRef = useRef(false);

  // 从哪里来回哪里去，无历史记录时兜底回首页
  const handleCancel = () => {
    if (window.history.length > 1) {
      router.back();
    } else {
      router.push('/');
    }
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!nodeId) { toast.warning('请选择所属板块'); return; }
    if (!content.trim()) { toast.warning('正文内容不能为空'); return; }
    if (submittingRef.current) return;
    submittingRef.current = true;

    try {
      await toast.promise(
        createPost({
          title,
          nodeId: Number(nodeId),
          content,
          tags: tags.split(',').map((t) => t.trim()).filter(Boolean),
        }),
        {
          loading: '发布中...',
          success: (result) => {
            submittingRef.current = false;
            setTimeout(() => {
              router.push(result?.id ? `/post/${result.id}` : '/');
              router.refresh();
            }, 600);
            return '主题发布成功 🎉';
          },
          error: (err) => {
            submittingRef.current = false;
            return err instanceof Error ? err.message : '发布失败，请稍后重试';
          },
        }
      );
    } catch {
      submittingRef.current = false;
    }
  };

  const isFormValid = Boolean(title && content.trim() && nodeId);

  return (
    <form id="createPostForm" onSubmit={handleSubmit}>
      <div className="form-row">
        <div className="form-group">
          <label className="form-label">所属板块<span className="required">*</span></label>
          <select className="form-select" required value={nodeId} onChange={(e) => setNodeId(e.target.value)}>
            <option value="">请选择板块...</option>
            {nodes.map((node) => (
              <option key={node.nodeId} value={node.nodeId}>{node.name}</option>
            ))}
          </select>
        </div>

        <div className="form-group">
          <label className="form-label">标签</label>
          <TagInput
            className="form-input"
            placeholder="输入标签名获取建议，多个用逗号分隔"
            value={tags}
            onChange={setTags}
          />
        </div>
      </div>

      <div className="form-group">
        <label className="form-label">帖子标题<span className="required">*</span></label>
        <input
          type="text"
          className="form-input"
          placeholder="请输入清晰明确的标题（5-100字）"
          maxLength={100}
          required
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
      </div>

      <div className="form-group">
        <label className="form-label">正文内容<span className="required">*</span></label>
        <div data-color-mode="light">
          <MDEditor
            value={content}
            onChange={(val) => setContent(val || '')}
            height={400}
            preview="live"
            visibleDragbar={false}
            textareaProps={{ placeholder: '支持 Markdown 语法，右侧实时预览...' }}
          />
        </div>
      </div>

      <div className="create-post-actions">
        {/* button + onClick，避免硬跳转丢失状态 */}
        <button
          type="button"
          className="btn btn-secondary"
          onClick={handleCancel}
        >
          取消
        </button>
        <button type="submit" className="btn btn-primary" disabled={!isFormValid}>
          发布主题
        </button>
      </div>
    </form>
  );
}