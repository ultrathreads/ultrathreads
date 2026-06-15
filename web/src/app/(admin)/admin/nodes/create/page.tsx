import NodeForm from '../NodeForm';

export default function CreatePage() {
  return (
    <div className="admin-body">
      <div className="admin-page-header">
        <h1 className="admin-page-title">新建节点</h1>
      </div>
      <div style={{ maxWidth: 720 }}>
        <NodeForm />
      </div>
    </div>
  );
}