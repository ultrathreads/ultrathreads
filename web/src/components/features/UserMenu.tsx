'use client';

import Link from 'next/link';
import { useTranslation } from '@/lib/i18n/i18n-client';
import { useAuth } from '@/hooks/use-auth';
import { useClickOutside } from '@/hooks/use-click-outside';
import Avatar from '@/components/ui/Avatar';
import AuthorLink from '@/components/ui/AuthorLink';

export default function UserMenu() {
  const { t } = useTranslation();
  const { user, isLoggedIn, isLoading, error, logout, displayName, avatarUrl, isAdmin, roles } = useAuth();
  const { isOpen, toggle, close, ref } = useClickOutside(false);

  // 数据未就绪或未登录时，不渲染任何内容
  if (isLoading || error || !isLoggedIn || !user) {
    return null;
  }

  // 角色映射：未匹配的 slug 降级显示原始值，空数组显示兜底文案
  const userRoles = user.roles ?? [];
  const roleDisplay = userRoles.length
    ? userRoles.map((role) => t(`roles.${role}`, { defaultValue: role })).join(', ')
    : t('roles._fallback');

  return (
    <div className="user-menu-wrapper" ref={ref}>
      
      <div className="user-menu-trigger" onClick={toggle} style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
        <Avatar 
          className="user-avatar" 
          src={avatarUrl} 
          alt={displayName} 
        />
        <span className="user-name">{displayName}</span>
        <span className={`user-arrow ${isOpen ? 'active' : ''}`}>▼</span>
      </div>

      <div className={`user-dropdown ${isOpen ? 'show' : ''}`} id="userDropdown">
        <div className="dropdown-header">
          <AuthorLink 
            author={`@${user.username}`} 
            authorSlug={user.slug} 
            className="user-username" 
            onClick={close}
          />
          <span className="user-level">{roleDisplay}</span>
        </div>
        <div className="dropdown-divider" />

        {/* ✅ 仅管理员可见的后台入口 */}
        {isAdmin && (
          <>
            <Link href="/admin" className="dropdown-item" onClick={close}>
              🛠️ 后台管理
            </Link>
            <div className="dropdown-divider" />
          </>
        )}

        <Link href="/settings/profile" className="dropdown-item" onClick={close}>
          👤 个人中心
        </Link>

        <Link href="/settings/account" className="dropdown-item" onClick={close}>
          ⚙️ 账号设置
        </Link>

        <div className="dropdown-divider" />

        <div
          className="dropdown-item danger"
          onClick={() => {
            close();
            logout();
          }}
        >
          🚪 退出登录
        </div>
      </div>
    </div>
  );
}