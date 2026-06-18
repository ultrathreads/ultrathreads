// src/app/(main)/users/[slug]/UserProfileCard.tsx
'use client';

import { UserEntity } from '@/types/domain';
import Avatar from '@/components/ui/Avatar';

interface Props {
  user: UserEntity;
}

export default function UserProfileCard({ user }: Props) {
  // 格式化时间为 YYYY/MM/DD 格式，避免时区和地区差异
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}/${month}/${day}`;
  };

  return (
    <>
      {/* 主要信息卡片 */}
      <div className="profile-card text-center">
        {/* 头像 */}
        <div className="profile-avatar-wrapper">
          <Avatar
            className="profile-author-avatar"
            src={user.avatar}
            alt={user.nickname || user.username}
          />
        </div>

        {/* 用户名和昵称 */}
        <div>
          <h2 className="profile-username">
            @{user.username}
            {user.levelName && (
              <span className="user-level">
                {user.levelName}
              </span>
            )}
          </h2>
          
          {user.nickname && (
            <div className="profile-nickname">
              {user.nickname}
            </div>
          )}
        </div>
        
        {/* 数据统计 */}
        <div className="profile-stats-grid">
          <div>
            <div className="stat-number">{user.topicCount}</div>
            <div className="stat-label">帖子</div>
          </div>
          <div>
            <div className="stat-number">{user.score}</div>
            <div className="stat-label">积分</div>
          </div>
          <div>
            <div className="stat-number">{user.followingCount || 0}</div>
            <div className="stat-label">关注</div>
          </div>
        </div>

        {/* 关注按钮 */}
        <button className="btn btn-primary btn-follow">关注 TA</button>

        {/* 注册日期和最后登录时间 */}
        <div className="profile-meta-info">
          {user.createdAt && (
            <span>注册于: {formatDate(user.createdAt)}</span>
          )}
          {user.lastLoginTime && (
            <span>最后登录: {formatDate(user.lastLoginTime)}</span>
          )}
        </div>
      </div>

      {/* 联系信息卡片 */}
      {user.website || user.description ? (
        <div className="profile-card text-left profile-contact-card">
          <h3 className="profile-section-title">联系信息</h3>
          
          {/* 个人主页 */}
          {user.website && (
            <div className="profile-contact-item">
              <div className="profile-label">个人主页</div>
              <div className="profile-link">
                <a href={user.website} target="_blank" rel="noopener noreferrer">
                  {user.website}
                </a>
              </div>
            </div>
          )}
          
          {/* 个人签名 */}
          {user.description && (
            <div className="profile-contact-item">
              <div className="profile-label">个人签名</div>
              <div className="profile-value">
                {user.description}
              </div>
            </div>
          )}
        </div>
      ) : null}
    </>
  );
}