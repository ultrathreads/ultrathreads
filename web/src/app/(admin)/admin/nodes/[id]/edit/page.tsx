import { getNode } from '../../actions';
import NodeForm from '../../NodeForm';

export const dynamic = 'force-dynamic';

interface EditPageProps {
  params: Promise<{ id: string }>;
}

export default async function EditPage({ params }: EditPageProps) {
  const { id } = await params;
  const nodeId = Number(id);

  // 防御性校验：防止无效 ID 触发后端请求
  if (!Number.isFinite(nodeId) || nodeId <= 0) {
    return (
      <div className="admin-body">
        <h1 className="admin-page-title">无效的节点ID</h1>
      </div>
    );
  }

  const node = await getNode(nodeId);

  return (
    <div className="admin-body">
      <div className="admin-page-header">
        <h1 className="admin-page-title">编辑节点: {node.name}</h1>
      </div>
      <div style={{ maxWidth: 720 }}>
        <NodeForm initialData={node} />
      </div>
    </div>
  );
}