// src/components/features/PostForm.tsx
'use client';

import { useState, useRef, FormEvent, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation'; // ✅ 1. 引入 useSearchParams
import MDEditor from '@uiw/react-md-editor';
import { toast } from 'sonner';
import { createRootPost, updateRootPost } from '@/services/post-service';
import { TagInput } from '@/components/ui/TagInput';
import type { NodeEntity } from '@/types/domain';
import { useSiteConfig } from '@/providers/SiteConfigProvider';

interface InitialData {
  slug: string;
  title: string;
  rawContent: string;
  nodeSlug: string;
  tags: string;
}

interface PostFormProps {
  nodes: NodeEntity[];
  initialData?: InitialData | null;
}

function buildPostDetailUrl(slug: string): string {
  const url = new URL(`/threads/${slug}`, window.location.origin);
  url.searchParams.set('refresh', '1');
  return url.pathname + url.search;
}

export function PostForm({ nodes, initialData }: PostFormProps) {
  const router = useRouter();
  const searchParams = useSearchParams(); // ✅ 2. 获取搜索参数实例
  const { recommendTags } = useSiteConfig();

  const isEditMode = Boolean(initialData);

  // ✅ 3. 初始化时优先使用 initialData，其次尝试从 URL 读取 node 参数
  const [title, setTitle] = useState(initialData?.title ?? '');
  const [nodeSlug, setNodeSlug] = useState(() => {
    if (initialData?.nodeSlug) return initialData.nodeSlug;
    // 仅在非编辑模式下读取 URL 参数
    const urlNode = searchParams.get('node');
    // 验证 URL 中的 node 是否在合法节点列表中
    if (urlNode && nodes.some((n) => n.slug === urlNode)) {
      return urlNode;
    }
    return '';
  });

  const [tags, setTags] = useState<string[]>(
    initialData?.tags
      ? initialData.tags.split(',').map((t) => t.trim()).filter(Boolean)
      : []
  );
  const [content, setContent] = useState(initialData?.rawContent ?? '');
  const [attempted, setAttempted] = useState(false);
  const submittingRef = useRef(false);

  // ✅ 4. 修复原代码中 targetSlug 未定义的 Bug
  const targetSlug = initialData?.slug;

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
      tags,
    };

    try {
      await toast.promise(
        isEditMode
          ? updateRootPost(initialData!.slug, payload)
          : createRootPost(payload),
        {
          loading: isEditMode ? '保存中...' : '发布中...',
          success: () => {
            submittingRef.current = false;
            setTimeout(() => {
              // ✅ 使用修复后的 targetSlug
              router.push(targetSlug ? buildPostDetailUrl(targetSlug) : '/');
            }, 600);
            return isEditMode ? '主贴已更新 ✅' : '主贴发布成功 🎉';
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
    <div className="post-form-container">
      <h1 className="post-form-header">
        {isEditMode ? '✏️ 编辑主帖' : '📝 发布新主贴'}
      </h1>
      <form id="postForm" onSubmit={handleSubmit} noValidate>
        {/* ... 表单内容保持不变 ... */}
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
            <label className="form-label">
              标签 <span className="form-hint">（选填，最多3个）</span>
            </label>
            <TagInput
              value={tags}
              onChange={setTags}
              placeholder="输入标签名获取建议，回车添加"
              recommendTags={recommendTags}
              maxTags={3}
            />
          </div>
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
            {isEditMode ? '保存修改' : '发布主贴'}
          </button>
        </div>
      </form>
    </div>
  );
}