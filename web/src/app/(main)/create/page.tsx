// src/app/create/page.tsx
'use client';

import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { mockBoards } from '@/lib/mock-data';

export default function CreatePostPage() {
  const router = useRouter();
  
  // 表单状态
  const [title, setTitle] = useState('');
  const [category, setCategory] = useState('');
  const [tags, setTags] = useState('');
  const [content, setContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  // 处理提交
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    // 模拟网络请求
    await new Promise((resolve) => setTimeout(resolve, 1200));

    console.log('发布数据:', { title, category, tags, content });
    
    setIsSubmitting(false);
    alert('🎉 帖子发布成功！');
    router.push('/');
  };

  return (
    <div className="main-body">
      <div className="create-post-container">
        <h1 className="create-post-header">✏️ 发布新主题</h1>
        
        <form id="createPostForm" onSubmit={handleSubmit}>
          {/* 第一行：板块与标签 */}
          <div className="form-row">
            <div className="form-group">
              <label className="form-label">所属板块<span className="required">*</span></label>
              <select 
                className="form-select" 
                required 
                value={category}
                onChange={(e) => setCategory(e.target.value)}
              >
                <option value="">请选择板块...</option>
                {mockBoards.map((board) => (
                  <option key={board.name} value={board.name}>
                    {board.icon} {board.name}
                  </option>
                ))}
              </select>
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

          {/* 第三行：正文内容 */}
          <div className="form-group">
            <label className="form-label">正文内容<span className="required">*</span></label>
            <div className="editor-wrapper">
              {/* 简易工具栏（仅作展示） */}
              <div className="editor-toolbar">
                <button type="button" className="toolbar-btn" title="加粗"><b>B</b></button>
                <button type="button" className="toolbar-btn" title="斜体"><i>I</i></button>
                <button type="button" className="toolbar-btn" title="引用">❝</button>
                <span className="toolbar-divider"></span>
                <button type="button" className="toolbar-btn" title="代码块">&lt;/&gt;</button>
                <button type="button" className="toolbar-btn" title="链接">🔗</button>
                <button type="button" className="toolbar-btn" title="图片">🖼️</button>
              </div>
              <textarea 
                className="form-textarea" 
                placeholder="支持 Markdown 语法..." 
                required
                value={content}
                onChange={(e) => setContent(e.target.value)}
              />
            </div>
          </div>

          {/* 底部操作按钮 */}
          <div className="create-post-actions">
            <a href="/" className="btn btn-secondary">取消</a>
            <button 
              type="submit" 
              className="btn btn-primary" 
              disabled={isSubmitting || !title || !content}
            >
              {isSubmitting ? '发布中...' : '发布主题'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}