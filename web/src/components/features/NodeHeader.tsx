// components/NodeHeader.tsx
import type { NodeEntity } from '@/types/domain';
import { NodeIcon } from '@/components/NodeIcon';

interface Props {
  node: NodeEntity | null;
}

export default function NodeHeader({ node }: Props) {
  if (!node) {
    return (
      <div className="board-title-wrapper">
        <span className="board-title-icon">🏠</span>
        <div className="board-title-text">
          <div className="board-title-name">全部帖子</div>
          <div className="board-title-desc">浏览论坛所有最新讨论</div>
        </div>
      </div>
    );
  }

  return (
    <div className="board-title-wrapper">
      <NodeIcon 
        icon={node.icon} 
        className="board-title-icon" 
      />
      <div className="board-title-text">
        <div className="board-title-name">{node.name}</div>
        <div className="board-title-desc">{node.description}</div>
      </div>
    </div>
  );
}