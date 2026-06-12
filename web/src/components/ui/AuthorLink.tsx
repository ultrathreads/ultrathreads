'use client';

import Link from 'next/link';

interface AuthorLinkProps {
  author: string;
  authorSlug?: string;
  className?: string;
  onClick?: (e: MouseEvent<HTMLAnchorElement>) => void;
}

export default function AuthorLink({ author, authorSlug, className, onClick }: AuthorLinkProps) {
  // 个人主页路由：优先使用 authorSlug，如果没有则降级使用 author 名称
  const profileUrl = `/users/${authorSlug ?? author}`;

  // 合并事件：外部事件 + 阻止冒泡
  const handleClick = (e: MouseEvent<HTMLAnchorElement>) => {
    // 先执行外部传入的点击事件
    if (onClick) {
      onClick(e);
    }
    // 固定阻止冒泡
    e.stopPropagation();
  };

  return (
    <Link 
      href={profileUrl} 
      className={className}
      onClick={handleClick}
    >
      {author}
    </Link>
  );
}