// components/NodeHeader.tsx
import type { NodeEntity } from '@/types/domain';
import { NodeIcon } from '@/components/NodeIcon';

// 定义一个通用的头部展示接口
export interface HeaderDisplayData {
  name: string;
  description?: string;
  icon?: string;
}

interface Props {
  node: NodeEntity | null;
  tag?: HeaderDisplayData | null;
}

export default function NodeHeader({ node, tag }: Props) {
  // 优先展示 Tag，其次展示 Node，最后展示默认首页状态
  const displayData = tag || node;

  if (!displayData) {
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

  // 标签没有 icon，所以这里做个简单的兜底
  const iconToRender = tag ? (tag.icon || '🏷️') : displayData.icon;

  return (
    <div className="board-title-wrapper">
      <NodeIcon 
        icon={iconToRender}
        className="board-title-icon" 
      />
      <div className="board-title-text">
        <div className="board-title-name">{displayData.name}</div>
        {/* 只有 Node 才有 description，所以仅在非 tag 且有值时渲染 */}
        {!tag && displayData.description && (
          <div className="board-title-desc">{displayData.description}</div>
        )}
      </div>
    </div>
  );
}