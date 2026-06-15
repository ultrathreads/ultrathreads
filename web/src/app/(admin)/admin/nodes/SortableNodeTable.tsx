'use client';

import { useState, useRef } from 'react';
import Link from 'next/link';
import DeleteButton from './DeleteButton';
import { updateNodeSort } from './actions'; // 需要新增此 action

interface NodeItem {
  id: number;
  name: string;
  icon: string;
  sortNo: number;
  topicCount: number;
  status: number;
}

interface Props {
  nodes: NodeItem[];
}

export default function SortableNodeTable({ nodes: initialNodes }: Props) {
  const [nodes, setNodes] = useState(initialNodes);
  const [draggingId, setDraggingId] = useState<number | null>(null);
  const [saving, setSaving] = useState(false);
  const dragOverId = useRef<number | null>(null);

  const handleDragStart = (id: number) => {
    setDraggingId(id);
  };

  const handleDragOver = (e: React.DragEvent, id: number) => {
    e.preventDefault(); // 必须阻止默认行为才能触发 drop
    dragOverId.current = id;
  };

  const handleDrop = async () => {
    if (draggingId === null || dragOverId.current === null || draggingId === dragOverId.current) {
      setDraggingId(null);
      return;
    }

    const fromIndex = nodes.findIndex((n) => n.id === draggingId);
    const toIndex = nodes.findIndex((n) => n.id === dragOverId.current);

    // 乐观更新：先调整本地顺序
    const reordered = [...nodes];
    const [moved] = reordered.splice(fromIndex, 1);
    reordered.splice(toIndex, 0, moved);

    // 重新计算 sortNo（按新索引生成）
    const updatedNodes = reordered.map((node, index) => ({
      ...node,
      sortNo: (index + 1) * 10, // 以10为步长，方便后续插入
    }));

    setNodes(updatedNodes);
    setDraggingId(null);
    setSaving(true);

    try {
      // 批量提交新的排序
      await updateNodeSort(updatedNodes.map((n) => ({ id: n.id, sortNo: n.sortNo })));
    } catch (err) {
      console.error('排序保存失败:', err);
      // 失败时回滚
      setNodes(initialNodes);
    } finally {
      setSaving(false);
    }
  };

  const handleDragEnd = () => {
    setDraggingId(null);
    dragOverId.current = null;
  };

  return (
    <div className="admin-table-wrapper" style={{ border: 'none', borderRadius: '8px' }}>
      {saving && (
        <div style={{ padding: '8px 16px', fontSize: 13, color: '#666', borderBottom: '1px solid #eee' }}>
          ⏳ 正在保存排序...
        </div>
      )}
      <table className="admin-table">
        <thead>
          <tr>
            <th style={{ width: 40 }}></th>
            <th>ID</th>
            <th>图标</th>
            <th>名称</th>
            <th>排序</th>
            <th>话题数</th>
            <th>状态</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          {nodes.length === 0 ? (
            <tr>
              <td colSpan={8}>
                <div className="admin-empty">
                  <div className="admin-empty-title">暂无节点</div>
                  <div className="admin-empty-desc">点击右上角按钮创建第一个节点</div>
                </div>
              </td>
            </tr>
          ) : (
            nodes.map((node) => (
              <tr
                key={node.id}
                draggable
                onDragStart={() => handleDragStart(node.id)}
                onDragOver={(e) => handleDragOver(e, node.id)}
                onDrop={handleDrop}
                onDragEnd={handleDragEnd}
                style={{
                  cursor: 'grab',
                  opacity: draggingId === node.id ? 0.4 : 1,
                  transition: 'opacity 0.15s',
                }}
              >
                {/* 拖拽手柄列 */}
                <td style={{ textAlign: 'center', color: '#bbb', userSelect: 'none' }}>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                    <circle cx="9" cy="6" r="1.5" />
                    <circle cx="15" cy="6" r="1.5" />
                    <circle cx="9" cy="12" r="1.5" />
                    <circle cx="15" cy="12" r="1.5" />
                    <circle cx="9" cy="18" r="1.5" />
                    <circle cx="15" cy="18" r="1.5" />
                  </svg>
                </td>
                <td>{node.id}</td>
                <td>
                  {node.icon?.startsWith('<svg') ? (
                    <span
                      dangerouslySetInnerHTML={{ __html: node.icon }}
                      style={{ display: 'inline-block', width: 24, height: 24 }}
                    />
                  ) : (
                    <span style={{ fontSize: 20 }}>{node.icon || '—'}</span>
                  )}
                </td>
                <td style={{ fontWeight: 500 }}>{node.name}</td>
                <td>{node.sortNo}</td>
                <td>{node.topicCount}</td>
                <td>
                  <span className={`admin-badge ${node.status === 0 ? 'admin-badge-success' : 'admin-badge-neutral'}`}>
                    {node.status === 0 ? '启用' : '禁用'}
                  </span>
                </td>
                <td>
                  <div className="admin-table-actions">
                    <Link href={`/admin/nodes/${node.id}/edit`} className="admin-btn admin-btn-secondary admin-btn-sm">
                      编辑
                    </Link>
                    <DeleteButton id={node.id} name={node.name} />
                  </div>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}