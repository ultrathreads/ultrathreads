'use client';

import { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useTranslation } from '@/lib/i18n/i18n-client';
import { getAllNodes } from '@/services/node-service';
import { getHotTags } from '@/services/tag-service';
import type { NodeEntity, TagEntity } from '@/types/domain';

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

export default function Sidebar() {
  const { t } = useTranslation(['common']);
  const router = useRouter();
  const searchParams = useSearchParams();

  const activeNodeId = searchParams.get('nodeId') ? Number(searchParams.get('nodeId')) : null;
  const activeTag = searchParams.get('tag') || null;

  const [collapsed, setCollapsed] = useState(true);
  const [nodes, setNodes] = useState<NodeEntity[]>([]);
  const [tags, setTags] = useState<TagEntity[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;

    const load = async () => {
      const [{ nodes: fetchedNodes }, { tags: fetchedTags }] = await Promise.all([
        getAllNodes(),
        getHotTags(),
      ]);

      if (!cancelled) {
        setNodes(fetchedNodes);
        setTags(fetchedTags);
        setLoading(false);
      }
    };

    load();
    return () => { cancelled = true; };
  }, []);

  const handleNodeClick = (nodeId: number) => {
    const params = new URLSearchParams(searchParams.toString());

    if (activeNodeId === nodeId) {
      params.delete('nodeId');
    } else {
      params.set('nodeId', String(nodeId));
      params.set('page', '1');
    }

    router.push(`/?${params.toString()}`);
  };

  const handleTagClick = (tagLabel: string) => {
    const params = new URLSearchParams(searchParams.toString());

    if (activeTag === tagLabel) {
      params.delete('tag');
    } else {
      params.set('tag', tagLabel);
      params.set('page', '1');
    }

    router.push(`/?${params.toString()}`);
  };

  return (
    <div className="sidebar-container">
      <div className={`sidebar-content-wrapper ${collapsed ? 'collapsed' : ''}`} id="sidebarContent">
        {/* 导航菜单 */}
        <div className="sidebar-section">
          <div className="sidebar-title">导航菜单</div>
          <ul className="forum-list">
            <li
              className={`forum-item ${activeNodeId === null && activeTag === null ? 'active' : ''}`}
              onClick={() => router.push('/')}
            >
              {t('common:home')}
            </li>
            <li
              className="forum-item"
              onClick={() => router.push('/settings/account')}
            >
              {t('common:mine')}
            </li>
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
                  <span className="forum-count">{node.topicCount}</span>
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="sidebar-section">
          <div className="sidebar-title">热门标签</div>
          {loading ? (
            <div className="forum-list-loading">加载中...</div>
          ) : (
            <div className="tag-cloud">
              {tags.map((tag) => (
                <span
                  key={tag.tagName}
                  className={`tag-item ${activeTag === tag.tagName ? 'active' : ''}`}
                  onClick={() => handleTagClick(tag.tagName)}
                >
                  {tag.tagName}
                </span>
              ))}
            </div>
          )}
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