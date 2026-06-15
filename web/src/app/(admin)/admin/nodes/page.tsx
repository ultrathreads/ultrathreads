import Link from 'next/link';
import { getNodes } from './actions';
import SortableNodeTable from './SortableNodeTable'; // 替换原来的内联表格

export default async function NodesPage() {
  const { nodes, pageInfo } = await getNodes();

  return (
    <div className="admin-body">
      <div className="admin-page-header">
        <div>
          <h1 className="admin-page-title">节点管理</h1>
          <p className="admin-page-desc">
            共 {pageInfo.total} 个节点 · 第 {pageInfo.page}/{Math.ceil(pageInfo.total / pageInfo.pageSize) || 1} 页
          </p>
        </div>
        <Link href="/admin/nodes/create" className="admin-btn admin-btn-primary">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M12 5v14M5 12h14"/>
          </svg>
          新建节点
        </Link>
      </div>

      <div className="admin-card" style={{ padding: 0 }}>
        {/* ✅ 替换为可拖拽的客户端表格组件 */}
        <SortableNodeTable nodes={nodes} />
      </div>
    </div>
  );
}