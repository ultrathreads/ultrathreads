'use client';

import { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import clsx from 'clsx';
import { useTranslation } from '@/lib/i18n/i18n-client';
import { getAllNodes } from '@/services/node-service';
import { getHotTags } from '@/services/tag-service';
import type { NodeEntity, TagEntity } from '@/types/domain';
import { NodeIcon } from '@/components/NodeIcon';
import { SidebarNav } from './SidebarNav';

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
      try {
        const [{ nodes: fetchedNodes }, { tags: fetchedTags }] = await Promise.all([
          getAllNodes(),
          getHotTags(),
        ]);

        if (!cancelled) {
          setNodes(fetchedNodes);
          setTags(fetchedTags);
        }
      } catch (error) {
        console.error('Failed to load sidebar data:', error);
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
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
        {/* ✅ 导航菜单 */}
        <div className="sidebar-section">
          <div className="sidebar-title">{t('common:navigation')}</div>
          <SidebarNav />
        </div>

        {/* 论坛板块 */}
        <div className="sidebar-section">
          <div className="sidebar-title">{t('common:forum_sections')}</div>
          {loading ? (
            <div className="forum-list-loading">{t('common:loading')}</div>
          ) : (
            <ul className="forum-list">
              {nodes.map((node) => {
                // ✅ 修复：将 isActive 的计算移入 map 回调内部
                const isActive = node.nodeId === activeNodeId;
                
                return (
                  <li
                    key={node.nodeId}
                    className={clsx('forum-item cursor-pointer', { active: isActive })}
                    onClick={() => handleNodeClick(node.nodeId)}
                  >
                    {/* ✅ 修复：使用正确的 isActive 变量 + 添加兜底值 */}
                    <NodeIcon 
                      icon={node.icon} 
                      className={isActive ? 'text-blue-600 dark:text-blue-400' : ''} 
                    />
                    <span className="truncate">{node.name}</span>
                    <span className="forum-count">{node.topicCount}</span>
                  </li>
                );
              })}
            </ul>
          )}
        </div>

        {/* 热门标签 */}
        <div className="sidebar-section">
          <div className="sidebar-title">{t('common:hot_tags')}</div>
          {loading ? (
            <div className="forum-list-loading">{t('common:loading')}</div>
          ) : (
            <div className="tag-cloud">
              {tags.map((tag) => (
                <span
                  key={tag.tagName}
                  className={clsx('tag-item', { active: activeTag === tag.tagName })}
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
        aria-label={collapsed ? t('common:expand_sidebar') : t('common:collapse_sidebar')}
      />
    </div>
  );
}