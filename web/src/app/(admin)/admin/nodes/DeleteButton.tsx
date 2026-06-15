'use client';
import { useState } from 'react';
import { deleteNode } from './actions';

export default function DeleteButton({ id, name }: { id: number; name: string }) {
  const [loading, setLoading] = useState(false);

  const handleDelete = async () => {
    if (!confirm(`确定要删除节点「${name}」吗？此操作不可恢复。`)) return;
    setLoading(true);
    const res = await deleteNode(id);
    setLoading(false);
    if (!res.success) alert(res.message);
  };

  return (
    <button 
      onClick={handleDelete} 
      disabled={loading} 
      className="admin-btn admin-btn-danger admin-btn-sm"
    >
      {loading ? '删除中...' : '删除'}
    </button>
  );
}