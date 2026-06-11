// src/app/(main)/users/[slug]/layout.tsx
import { notFound } from 'next/navigation';
import { getUserBySlug } from '@/services/user-service';
import UserProfileCard from './UserProfileCard';

interface Props {
  children: React.ReactNode;
  params: Promise<{ slug: string }>;
}

export default async function UserLayout({ children, params }: Props) {
  const { slug } = await params;
  const user = await getUserBySlug(slug);

  if (!user) notFound();

  return (
    // 复用全局设置页的容器布局，完美适配左右分栏
    <div className="profile-container">
      {/* 左侧：帖子列表内容区 */}
      <main className="profile-content">
        {children}
      </main>

      {/* 右侧：个人信息卡片，使用 sticky 悬浮 */}
      <aside className="profile-sidebar">
        <div className="sticky top-24">
          <UserProfileCard user={user} />
        </div>
      </aside>
    </div>
  );
}