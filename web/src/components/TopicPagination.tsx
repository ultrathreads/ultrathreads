'use client';

import { useRouter, usePathname } from 'next/navigation';
import Pagination from '@/components/Pagination';

interface TopicPaginationProps {
  totalItems: number;
  pageSize: number;
  currentPage: number;
}

export default function TopicPagination({
  totalItems,
  pageSize,
  currentPage,
}: TopicPaginationProps) {
  const router = useRouter();
  const pathname = usePathname();

  const handlePageChange = (page: number) => {
    // 保留当前 URL 的其他查询参数，仅更新 page
    const params = new URLSearchParams(window.location.search);
    if (page === 1) {
      params.delete('page');
    } else {
      params.set('page', String(page));
    }
    router.push(`${pathname}?${params.toString()}`);
  };

  return (
    <Pagination
      totalItems={totalItems}
      pageSize={pageSize}
      currentPage={currentPage}
      onPageChange={handlePageChange}
    />
  );
}