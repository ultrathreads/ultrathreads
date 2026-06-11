'use client';

import Link from 'next/link';

interface AuthorLinkProps {
  author: string;
  authorId?: string | number;
  className?: string;
}

export default function AuthorLink({ author, authorId, className }: AuthorLinkProps) {
  // 个人主页路由：优先使用 authorId，如果没有则降级使用 author 名称
  const profileUrl = `/users/${authorId ?? author}`;

  return (
    <Link 
      href={profileUrl} 
      className={className}
      onClick={(e) => e.stopPropagation()} // ✅ 防止在列表项中点击作者时，触发整个帖子的跳转
    >
      {author}
    </Link>
  );
}