'use client';
import { useState } from 'react';
import { ForumBoard, Tag } from '@/types';
import { useTranslation } from '@/lib/i18n-client';

interface Props {
  boards?: ForumBoard[];
  tags?: Tag[];
}

export default function Sidebar({ boards, tags }: Props) {
  const { t } = useTranslation(['common']);
  const [collapsed, setCollapsed] = useState(true);

  // 1. 定义模拟数据（Mock Data）
  const mockBoards: ForumBoard[] = [
    { name: '技术交流', icon: '💻', count: 128 },
    { name: '生活日常', icon: '☕', count: 256 },
    { name: '开源项目', icon: '🚀', count: 64 },
    { name: '问答求助', icon: '🙋‍♂️', count: 89 },
  ];

  const mockTags: Tag[] = [
    { label: 'React' },
    { label: 'Next.js' },
    { label: 'TypeScript' },
    { label: 'TailwindCSS' },
    { label: 'Node.js' },
  ];

  // 2. 使用传入的 props，如果没有传入则使用模拟数据兜底，并确保它一定是数组
  const safeBoards = Array.isArray(boards) ? boards : mockBoards;
  const safeTags = Array.isArray(tags) ? tags : mockTags;

  return (
    <div className="sidebar-container">
      <div className={`sidebar-content-wrapper ${collapsed ? 'collapsed' : ''}`} id="sidebarContent">
        <div className="sidebar-section">
          <div className="sidebar-title">导航菜单</div>
          <ul className="forum-list">
            <li className="forum-item">{t('common:home')}</li>
            <li className="forum-item">{t('common:mine')}</li>
          </ul>
        </div>
        
        <div className="sidebar-section">
          <div className="sidebar-title">论坛板块</div>
          <ul className="forum-list">
            {/* 3. 使用安全的变量进行 map 渲染，彻底告别 undefined 报错 */}
            {safeBoards.map((b) => (
              <li key={b.name} className="forum-item">
                <span>{b.icon}</span>
                <span>{b.name}</span>
                <span className="forum-count">{b.count}</span>
              </li>
            ))}
          </ul>
        </div>
        
        <div className="sidebar-section">
          <div className="sidebar-title">热门标签</div>
          <div className="tag-cloud">
            {safeTags.map((t) => (
              <span key={t.label} className="tag-item">{t.label}</span>
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