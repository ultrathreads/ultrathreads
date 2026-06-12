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

  const initialNodeSlug = (() => {
    const slug = searchParams.get('nodeSlug');
    if (slug && nodes.some((n) => String(n.nodeSlug) === slug)) return slug;
    return '';
  })();

  const [title, setTitle] = useState('');
  const [nodeSlug, setNodeSlug] = useState(initialNodeSlug);
  const [tags, setTags] = useState('');
  const [content, setContent] = useState('');
  const [attempted, setAttempted] = useState(false);
  const submittingRef = useRef(false);

  const handleCancel = () => {
    if (window.history.length > 1) router.back();
    else router.push('/');
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setAttempted(true);

    // ✅ 纯前端校验失败 → 仅激活行内错误，不弹 Toast
    if (!nodeSlug || !title.trim() || !content.trim()) return;

    if (submittingRef.current) return;
    submittingRef.current = true;

    try {
      await toast.promise(
        createPost({
          title: title.trim(),
          nodeSlug: Number(nodeSlug),
          content,
          tags: tags.split(',').map((t) => t.trim()).filter(Boolean),
        }),
        {
          loading: '发布中...',
          success: (result) => {
            submittingRef.current = false;
            setTimeout(() => {
              router.push(result?.slug ? `/threads/${result.slug}` : '/');
              router.refresh();
            }, 600);
            return '主题发布成功 🎉';
          },
          // ✅ Toast 仅处理异步结果：网络异常 / 服务端业务校验拒绝
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

  const errors = {
    nodeSlug: attempted && !nodeSlug,
    title: attempted && !title.trim(),
    content: attempted && !content.trim(),
  };

  const isFormValid = Boolean(nodeSlug && title.trim() && content.trim());

  return (
    <form id="createPostForm" onSubmit={handleSubmit} noValidate>
      <div className="form-row">
        <div className="form-group">
          <label className="form-label">所属板块<span className="required">*</span></label>
          <select
            className={`form-select ${errors.nodeSlug ? 'form-error' : ''}`}
            value={nodeSlug}
            onChange={(e) => { setNodeSlug(e.target.value); if (attempted) setAttempted(true); }}
          >
            <option value="">请选择板块...</option>
            {nodes.map((node) => (
              <option key={node.slug} value={node.slug}>{node.name}</option>
            ))}
          </select>
          {errors.nodeSlug && <p className="form-error-text">请选择所属板块</p>}
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
          className={`form-input ${errors.title ? 'form-error' : ''}`}
          placeholder="请输入清晰明确的标题（5-100字）"
          maxLength={100}
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
        {errors.title && <p className="form-error-text">请输入帖子标题</p>}
      </div>

      <div className="form-group">
        <label className="form-label">正文内容<span className="required">*</span></label>
        <div className={errors.content ? 'md-editor-error' : ''} data-color-mode="light">
          <MDEditor
            value={content}
            onChange={(val) => setContent(val || '')}
            height={400}
            preview="live"
            visibleDragbar={false}
            textareaProps={{ placeholder: '支持 Markdown 语法，右侧实时预览...' }}
          />
        </div>
        {errors.content && <p className="form-error-text">正文内容不能为空</p>}
      </div>

      <div className="create-post-actions">
        <button type="button" className="btn btn-secondary" onClick={handleCancel}>
          取消
        </button>
        <button
          type="submit"
          className="btn btn-primary"
          aria-disabled={!isFormValid}
          style={{ opacity: isFormValid ? 1 : 0.7 }}
        >
          发布主题
        </button>
      </div>
    </form>
  );
}