// components/NodeHeader.tsx
'use client';

import { useTranslation } from '@/lib/i18n/i18n-client';
import type { NodeEntity } from '@/types/domain';
import { NodeIcon } from '@/components/ui/NodeIcon';

export interface HeaderDisplayData {
  name: string;
  description?: string;
  icon?: string;
}

interface Props {
  node: NodeEntity | null;
  tag?: HeaderDisplayData | null;
}

// ✅ 移除 useAuth、Link、发帖按钮相关逻辑，回归纯展示组件
export default function NodeHeader({ node, tag }: Props) {
  const { t } = useTranslation(['common']);
  const displayData = tag || node;

  if (!displayData) {
    return (
      <div className="board-title-wrapper">
        <span className="board-title-icon">🏠</span>
        <div className="board-title-text">
          <div className="board-title-name">{t('common:all_threads')}</div>
          <div className="board-title-desc">{t('common:all_threads_desc')}</div>
        </div>
      </div>
    );
  }

  const iconToRender = tag ? (tag.icon || '🏷️') : displayData.icon;

  return (
    <div className="board-title-wrapper">
      <NodeIcon icon={iconToRender} className="board-title-icon" />
      <div className="board-title-text">
        <div className="board-title-name">{displayData.name}</div>
        {!tag && displayData.description && (
          <div className="board-title-desc">{displayData.description}</div>
        )}
      </div>
    </div>
  );
}