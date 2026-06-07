'use client';

import { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { TagEntity } from '@/types/domain';
import { useTranslation } from '@/lib/i18n/i18n-client';
import { getAllNodes, type ForumNode } from '@/services/node-service';

// ==================== 工具函数 ====================
const getIconByName = (name: string): string => {
  const iconMap: Record<string, string> = {
    '公告': '📢',
    '问答': '🙋‍♂️',
    '教程': '📚',
    '分享': '💡',
    '技术交流': '💻',
    '生活日常': '☕',
    '开源项目': '🚀',
  };
  return iconMap[name] || '📁';
};

interface Props {
  tags?: TagEntity[];
}

export default function Sidebar({ tags }: Props) {
  const { t } = useTranslation(['common']);
  const router = useRouter();
  const searchParams = useSearchParams();

  // ✅ 从 URL 读取当前激活的节点，保证刷新页面或前进后退时高亮状态正确
  const activeNodeId = searchParams.get('nodeId') ? Number(searchParams.get('nodeId')) : null;

  const [collapsed, setCollapsed] = useState(true);
  const [nodes, setNodes] = useState<ForumNode[]>([]);
  const [loading, setLoading] = useState(true);

  const mockTags: TagEntity[] = [
    { label: 'React' },
    { label: 'Next.js' },
    { label: 'TypeScript' },
    { label: 'TailwindCSS' },
    { label: 'Node.js' },
  ];

  // ---------- 初始化：通过 Service 获取板块列表 ----------
  useEffect(() => {
    let cancelled = false;

    const loadNodes = async () => {
      const { nodes: fetchedNodes } = await getAllNodes();
      if (!cancelled) {
        setNodes(fetchedNodes);
        setLoading(false);
      }
    };

    loadNodes();
    return () => { cancelled = true; };
  }, []);

  // ✅ 核心修改：点击时改变 URL，触发 RSC 重渲染，移除 onNodeSelect 回调
  const handleNodeClick = (nodeId: number) => {
    const params = new URLSearchParams(searchParams.toString());
    
    if (activeNodeId === nodeId) {
      // 再次点击同一节点时取消选中，回到全部帖子
      params.delete('nodeId');
    } else {
      params.set('nodeId', String(nodeId));
      // 切换节点时重置到第一页，避免新节点下没有对应页码的数据
      params.set('page', '1');
    }

    router.push(`/?${params.toString()}`);
  };

  const safeTags = Array.isArray(tags) ? tags : mockTags;

  return (
    <div className="sidebar-container">
      <div className={`sidebar-content-wrapper ${collapsed ? 'collapsed' : ''}`} id="sidebarContent">
        <div className="sidebar-section">
          <div className="sidebar-title">导航菜单</div>
          <ul className="forum-list">
            {/* ✅ 首页也改为 URL 驱动，清空 nodeId */}
            <li 
              className={`forum-item cursor-pointer ${activeNodeId === null ? 'active' : ''}`}
              onClick={() => router.push('/')}
            >
              {t('common:home')}
            </li>
            <li className="forum-item">{t('common:mine')}</li>
          </ul>
        </div>
        
        {/* 论坛板块 */}
        <div className="sidebar-section">
          <div className="sidebar-title">论坛板块</div>
          {loading ? (
            <div className="forum-list-loading">加载中...</div>
          ) : (
            <ul className="forum-list">
              {nodes.map((node) => (
                <li
                  key={node.nodeId}
                  className={`forum-item cursor-pointer ${
                    node.nodeId === activeNodeId ? 'active' : ''
                  }`}
                  onClick={() => handleNodeClick(node.nodeId)}
                >
                  <span>{getIconByName(node.name)}</span>
                  <span>{node.name}</span>
                  <span className="forum-count">{node.postCount}</span>
                </li>
              ))}
            </ul>
          )}
        </div>
        
        <div className="sidebar-section">
          <div className="sidebar-title">热门标签</div>
          <div className="tag-cloud">
            {safeTags.map((tag) => (
              <span key={tag.label} className="tag-item">{tag.label}</span>
            ))}
          </div>
        </div>
      </div>
      
      <div className="sidebar-track" />
      <button
        id="toggle-btn"
        onClick={() => setCollapsed(!collapsed)}
        dangerouslySetInnerHTML={{ __html: collapsed ? '&#9654;' : '&#9664;' }}
      />
    </div>
  );
}