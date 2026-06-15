'use client';

import { useState } from 'react';
import Link from 'next/link'; 
import { useRouter } from 'next/navigation';
import { createNode, updateNode, type NodePayload } from './actions';

interface Props {
  initialData?: { id: number; name: string; description: string; icon: string; sortNo: number; status: number };
}

export default function NodeForm({ initialData }: Props) {
  const router = useRouter();
  const isEdit = !!initialData;
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    const fd = new FormData(e.currentTarget);
    const payload: NodePayload = {
      name: fd.get('name') as string,
      description: (fd.get('description') as string) || '',
      icon: (fd.get('icon') as string) || '',
      sortNo: Number(fd.get('sortNo')) || 0,
      status: Number(fd.get('status')),
    };

    const result = isEdit 
      ? await updateNode(initialData!.id, payload) 
      : await createNode(payload);

    setLoading(false);
    if (!result.success) {
      setError(result.message);
    } else {
      router.push('/admin/nodes');
    }
  };

  return (
    <div className="admin-card">
      <div className="admin-card-header">
        <h2 className="admin-card-title">{isEdit ? '编辑节点' : '新建节点'}</h2>
      </div>

      {error && (
        <div className="admin-toast error" style={{ marginBottom: 20, position: 'relative', right: 'auto', top: 'auto', minWidth: 'auto' }}>
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        {/* 名称 */}
        <div className="admin-form-group">
          <label className="admin-form-label">
            节点名称 <span className="required">*</span>
          </label>
          <input 
            name="name" 
            required 
            maxLength={32}
            defaultValue={initialData?.name} 
            className="admin-form-input" 
            placeholder="请输入节点名称，最多32个字符"
          />
        </div>

        {/* 图标 */}
        <div className="admin-form-group">
          <label className="admin-form-label">节点图标</label>
          <input 
            name="icon" 
            defaultValue={initialData?.icon} 
            className="admin-form-input" 
            placeholder="输入图片 URL 或字体图标类名"
          />
          <p className="admin-form-hint">支持图片链接或 CSS 类名，留空则不显示图标</p>
        </div>

        {/* 描述 */}
        <div className="admin-form-group">
          <label className="admin-form-label">描述</label>
          <textarea 
            name="description" 
            rows={3} 
            defaultValue={initialData?.description} 
            className="admin-form-textarea" 
            placeholder="节点的简要介绍"
          />
        </div>

        {/* 排序与状态双列布局 */}
        <div className="admin-form-row">
          <div className="admin-form-group">
            <label className="admin-form-label">排序号</label>
            <input 
              name="sortNo" 
              type="number" 
              defaultValue={initialData?.sortNo ?? 0} 
              className="admin-form-input" 
            />
            <p className="admin-form-hint">数字越小越靠前</p>
          </div>
          <div className="admin-form-group">
            <label className="admin-form-label">状态</label>
            <select 
              name="status" 
              defaultValue={initialData?.status ?? 0} 
              className="admin-form-select"
            >
              <option value={0}>启用</option>
              <option value={1}>禁用</option>
            </select>
          </div>
        </div>

        {/* 提交按钮 */}
        <div style={{ paddingTop: 8 }}>
          <button 
            type="submit" 
            disabled={loading} 
            className="admin-btn admin-btn-primary"
          >
            {loading && <span className="admin-spinner" style={{ width: 14, height: 14, borderWidth: 2, marginRight: 6 }} />}
            {isEdit ? '保存修改' : '创建节点'}
          </button>
          <Link 
            href="/admin/nodes" 
            className="admin-btn admin-btn-secondary" 
            style={{ marginLeft: 12 }}
          >
            取消
          </Link>
        </div>
      </form>
    </div>
  );
}

// 注意：需要在文件顶部补充 import Link from 'next/link';