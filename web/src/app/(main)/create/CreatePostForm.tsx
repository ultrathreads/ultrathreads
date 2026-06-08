// src/app/create/CreatePostForm.tsx
'use client';

import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import MDEditor from '@uiw/react-md-editor';
import { createPost } from '@/services/post-service';
import type { NodeEntity } from '@/types/domain';

interface CreatePostFormProps {
  nodes: NodeEntity[];
}

export function CreatePostForm({ nodes }: CreatePostFormProps) {
  const router = useRouter();
  
  const [title, setTitle] = useState('');
  // ✅ 使用 node.id 作为 value，而非 name
  const [nodeId, setNodeId] = useState('');
  const [tags, setTags] = useState('');
  const [content, setContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setErrorMsg('');

    if (!nodeId) {
      setErrorMsg('请选择所属板块');
      return;
    }

    if (!content.trim()) {
      setErrorMsg('正文内容不能为空');
      return;
    }
    
    setIsSubmitting(true);

    try {
      const result = await createPost({
        title,
        nodeId: Number(nodeId),
        content,
        tags: tags.split(',').map(t => t.trim()).filter(Boolean),
      });

      // ✅ 跳转到新帖子详情页，若无返回则跳首页
      router.push(result?.id ? `/post/${result.id}` : '/');
      router.refresh(); // ✅ 刷新服务端缓存，使首页列表立即更新
    } catch (err) {
      console.error('[CreatePost] Failed:', err);
      setErrorMsg(err instanceof Error ? err.message : '发布失败，请稍后重试');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form id="createPostForm" onSubmit={handleSubmit}>
      {/* 第一行：板块与标签 */}
      <div className="form-row">
        <div className="form-group">
          <label className="form-label">所属板块<span className="required">*</span></label>
          <select 
            className="form-select" 
            required 
            value={nodeId}
            onChange={(e) => setNodeId(e.target.value)}
          >
            <option value="">请选择板块...</option>
            {/* ✅ 渲染真实节点，value 使用 id */}
            {nodes.map((node) => (
              <option key={node.nodeId} value={node.nodeId}>
                {node.name}
              </option>
            ))}
          </select>
          {/* 无节点时的友好提示 */}
          {nodes.length === 0 && (
            <p style={{ color: '#999', fontSize: '12px', marginTop: 4 }}>
              暂无可用板块，请联系管理员
            </p>
          )}
        </div>
        
        <div className="form-group">
          <label className="form-label">标签</label>
          <input 
            type="text" 
            className="form-input" 
            placeholder="多个标签用逗号分隔，如：Vue3, 性能优化"
            value={tags}
            onChange={(e) => setTags(e.target.value)}
          />
        </div>
      </div>

      {/* 第二行：标题 */}
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

      {/* 第三行：Markdown 编辑器 */}
      <div className="form-group">
        <label className="form-label">正文内容<span className="required">*</span></label>
        <div data-color-mode="light">
          <MDEditor
            autoFocus
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
      </div>

      {/* 底部操作按钮 */}
      <div className="create-post-actions">
        <a href="/" className="btn btn-secondary">取消</a>
        <button 
          type="submit" 
          className="btn btn-primary" 
          disabled={isSubmitting || !title || !content.trim() || !nodeId}
        >
          {isSubmitting ? '发布中...' : '发布主题'}
        </button>
      </div>
    </form>
  );
}