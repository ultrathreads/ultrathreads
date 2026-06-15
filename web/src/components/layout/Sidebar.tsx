'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter, useParams } from 'next/navigation';
import clsx from 'clsx';
import { useTranslation } from '@/lib/i18n/i18n-client';
import { getAllNodes } from '@/services/node-service';
import { getHotTags } from '@/services/tag-service';
import type { NodeEntity, TagEntity } from '@/types/domain';
import { NodeIcon } from '@/components/ui/NodeIcon';
import { SidebarNav } from './SidebarNav';

// ✅ 常量提取，避免魔法字符串散落在代码中
const SIDEBAR_COLLAPSED_KEY = 'sidebar-collapsed';

export default function Sidebar() {
  const { t } = useTranslation(['common']);
  const router = useRouter();
  const params = useParams<{ slug?: string }>();

  // ✅ 初始值从 localStorage 读取，SSR 阶段默认为 true 避免水合不匹配
  const [collapsed, setCollapsed] = useState<boolean>(false);
  const [nodes, setNodes] = useState<NodeEntity[]>([]);
  const [tags, setTags] = useState<TagEntity[]>([]);
  const [loading, setLoading] = useState(true);

  const activeSlug = params?.slug || null;

  // ✅ 挂载后同步 localStorage 状态，仅在客户端执行
  useEffect(() => {
    const stored = localStorage.getItem(SIDEBAR_COLLAPSED_KEY);
    if (stored !== null) {
      setCollapsed(stored === 'true');
    }
  }, []);

  // ✅ 切换时同时写入 localStorage
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
                const isActive = node.slug === activeSlug;
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
                const isActive = tag.slug === activeSlug;
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
        dangerouslySetInnerHTML={{ __html: collapsed ? '&#9654;' : '&#9664;' }}
        aria-label={collapsed ? t('common:expand_sidebar') : t('common:collapse_sidebar')}
      />
    </div>
  );
}