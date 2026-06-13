// src/components/features/PostForm.tsx
'use client';

import { useState, useRef, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import MDEditor from '@uiw/react-md-editor';
import { toast } from 'sonner';
import { createRootPost, updateRootPost } from '@/services/post-service';
import { TagInput } from '@/components/ui/TagInput';
import type { NodeEntity } from '@/types/domain';

interface InitialData {
  slug: string;
  title: string;
  rawContent: string;
  nodeSlug: string;
  tags: string;
}

// ✅ 接口名同步更新
interface PostFormProps {
  nodes: NodeEntity[];
  initialData?: InitialData | null;
}

// ✅ 组件名更新为 PostForm
export function PostForm({ nodes, initialData }: PostFormProps) {
  const router = useRouter();
  const isEditMode = Boolean(initialData);

  // ✅ 使用 initialData 初始化状态，避免客户端闪烁
  const [title, setTitle] = useState(initialData?.title ?? '');
  const [nodeSlug, setNodeSlug] = useState(initialData?.nodeSlug ?? '');
  const [tags, setTags] = useState(initialData?.tags ?? '');
  const [content, setContent] = useState(initialData?.rawContent ?? '');
  const [attempted, setAttempted] = useState(false);
  const submittingRef = useRef(false);

  const handleCancel = () => {
    if (isEditMode && initialData) {
      router.push(`/threads/${initialData.slug}`);
    } else if (window.history.length > 1) {
      router.back();
    } else {
      router.push('/');
    }
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setAttempted(true);

    if (!nodeSlug || !title.trim() || !content.trim()) return;
    if (submittingRef.current) return;
    submittingRef.current = true;

    const payload = {
      title: title.trim(),
      nodeSlug,
      content,
      tags: tags.split(',').map((t) => t.trim()).filter(Boolean),
    };

    try {
      await toast.promise(
        isEditMode
          ? updateRootPost(initialData!.slug, payload)
          : createRootPost(payload),
        {
          loading: isEditMode ? '保存中...' : '发布中...',
          success: (result) => {
            submittingRef.current = false;
            setTimeout(() => {
              const targetSlug = isEditMode ? initialData!.slug : result?.slug;
              router.push(targetSlug ? `/threads/${targetSlug}` : '/');
              router.refresh();
            }, 600);
            return isEditMode ? '主题已更新 ✅' : '主题发布成功 🎉';
          },
          error: (err) => {
            submittingRef.current = false;
            return err instanceof Error ? err.message : '操作失败，请稍后重试';
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
    // ✅ form id 同步去除 create 前缀
    <form id="postForm" onSubmit={handleSubmit} noValidate>
      <div className="form-row">
        <div className="form-group">
          <label className="form-label">
            所属板块<span className="required">*</span>
          </label>
          <select
            className={`form-select ${errors.nodeSlug ? 'form-error' : ''}`}
            value={nodeSlug}
            onChange={(e) => {
              setNodeSlug(e.target.value);
              if (attempted) setAttempted(true);
            }}
            disabled={isEditMode}
          >
            <option value="">请选择板块...</option>
            {nodes.map((node) => (
              <option key={node.slug} value={node.slug}>
                {node.name}
              </option>
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
        <label className="form-label">
          帖子标题<span className="required">*</span>
        </label>
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
        <label className="form-label">
          正文内容<span className="required">*</span>
        </label>
        <div
          className={errors.content ? 'md-editor-error' : ''}
          data-color-mode="light"
        >
          <MDEditor
            value={content}
            onChange={(val) => setContent(val || '')}
            height={400}
            preview="live"
            visibleDragbar={false}
            textareaProps={{
              placeholder: '支持 Markdown 语法，右侧实时预览...',
            }}
          />
        </div>
        {errors.content && <p className="form-error-text">正文内容不能为空</p>}
      </div>

      {/* ✅ CSS 类名也建议后续同步改为 post-form-actions */}
      <div className="post-form-actions">
        <button
          type="button"
          className="btn btn-secondary"
          onClick={handleCancel}
        >
          取消
        </button>
        <button
          type="submit"
          className="btn btn-primary"
          aria-disabled={!isFormValid}
          style={{ opacity: isFormValid ? 1 : 0.7 }}
        >
          {isEditMode ? '保存修改' : '发布主题'}
        </button>
      </div>
    </form>
  );
}