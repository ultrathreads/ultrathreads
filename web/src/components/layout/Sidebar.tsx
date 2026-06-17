'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { useRouter, useParams, useSearchParams } from 'next/navigation';
import clsx from 'clsx';
import { useTranslation } from '@/lib/i18n/i18n-client';
import { getAllNodes } from '@/services/node-service';
import { getHotTags } from '@/services/tag-service';
import type { NodeEntity, TagEntity } from '@/types/domain';
import { NodeIcon } from '@/components/ui/NodeIcon';
import { SidebarNav } from './SidebarNav';

const SIDEBAR_COLLAPSED_KEY = 'sidebar-collapsed';

export default function Sidebar() {
  const { t } = useTranslation();
  const router = useRouter();
  const params = useParams<{ slug?: string }>();
  const searchParams = useSearchParams();

  const [collapsed, setCollapsed] = useState<boolean>(false);
  const [nodes, setNodes] = useState<NodeEntity[]>([]);
  const [tags, setTags] = useState<TagEntity[]>([]);
  const [loading, setLoading] = useState(true);

  const activeNodeSlug = useMemo(() => {
    return params?.slug || searchParams.get('node') || null;
  }, [params?.slug, searchParams]);

  const activeTagSlug = useMemo(() => {
    return params?.slug || searchParams.get('tag') || null;
  }, [params?.slug, searchParams]);

  useEffect(() => {
    const stored = localStorage.getItem(SIDEBAR_COLLAPSED_KEY);
    if (stored !== null) {
      setCollapsed(stored === 'true');
    }
  }, []);


  const toggleCollapsed = useCallback(() => {
    setCollapsed((prev) => {
      const next = !prev;
      localStorage.setItem(SIDEBAR_COLLAPSED_KEY, String(next));
      return next;
    });
  }, []);

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

  const handleTagClick = (tagSlug: string) => {
    router.push(`/tags/${tagSlug}?page=1`);
  };

  return (
    <div className="sidebar-container">
      <div className={`sidebar-content-wrapper ${collapsed ? 'collapsed' : ''}`} id="sidebarContent">
        {/* 导航菜单 */}
        <div className="sidebar-section">
          <div className="sidebar-title">{t('navigation')}</div>
          <SidebarNav />
        </div>

        {/* 论坛板块 */}
        <div className="sidebar-section">
          <div className="sidebar-title">{t('forum_sections')}</div>
          {loading ? (
            <div className="forum-list-loading">{t('loading')}</div>
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
          <div className="sidebar-title">{t('hot_tags')}</div>
          {loading ? (
            <div className="forum-list-loading">{t('loading')}</div>
          ) : (
            <div className="tag-cloud">
              {tags.map((tag) => {
                const isActive = tag.slug === activeTagSlug;
                return (
                  <span
                    key={tag.slug}
                    className={clsx('tag-item', { active: isActive })}
                    onClick={() => handleTagClick(tag.slug)}
                  >
                    {tag.name}
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
        onClick={toggleCollapsed}
        title={collapsed ? t('expand_sidebar') : t('collapse_sidebar')}
        aria-label={collapsed ? t('expand_sidebar') : t('collapse_sidebar')}
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <rect width="18" height="18" x="3" y="3" rx="2"></rect>
          <path d="M9 3v18"></path>
        </svg>
      </button>
    </div>
  );
}