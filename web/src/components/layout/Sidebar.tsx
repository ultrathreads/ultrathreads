'use client';

import { useState, useEffect } from 'react';
import { useRouter, useSearchParams, useParams, usePathname } from 'next/navigation';
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

  const params = useParams<{ slug?: string; tagId?: string }>(); 
  const pathname = usePathname(); 

  const [collapsed, setCollapsed] = useState(true);
  const [nodes, setNodes] = useState<NodeEntity[]>([]);
  const [tags, setTags] = useState<TagEntity[]>([]);
  const [loading, setLoading] = useState(true);

  const activeNodeSlug = params?.slug || null;
  const activeTagId = params?.tagId ? Number(params.tagId) : null;

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

  const handleNodeClick = (nodeSlug: string) => {
    router.push(`/nodes/${nodeSlug}?page=1`);
  };

  const handleTagClick = (tagId: number) => {
     router.push(`/tags/${tagId}?page=1`);
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
                const isActive = node.slug === activeNodeSlug;
                
                return (
                  <li
                    key={node.slug}
                    className={clsx('forum-item', { active: isActive })}
                    onClick={() => handleNodeClick(node.slug)}
                  >
                    <NodeIcon icon={node.icon} />
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
              {tags.map((tag) => {
                const isActive = tag.tagId === activeTagId;

                return (
                  <span
                    key={tag.tagId}
                    className={clsx('tag-item', { active: isActive })}
                    onClick={() => handleTagClick(tag.tagId)}
                  >
                    {tag.tagName}
                  </span>
                );
              })}
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