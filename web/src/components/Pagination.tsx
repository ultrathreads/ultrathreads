'use client';

import { useMemo } from 'react';

interface PaginationProps {
  totalItems: number;
  pageSize: number;
  currentPage: number;
  onPageChange: (page: number) => void;
}

export default function Pagination({
  totalItems,
  pageSize,
  currentPage,
  onPageChange,
}: PaginationProps) {
  const totalPages = Math.ceil(totalItems / pageSize);

  const visiblePages = useMemo(() => {
    const pages: (number | string)[] = [];
    if (totalPages <= 0) return [1];

    const delta = 2;
    const left = Math.max(2, currentPage - delta);
    const right = Math.min(totalPages - 1, currentPage + delta);

    pages.push(1);
    if (left > 2) pages.push('...');
    for (let i = left; i <= right; i++) pages.push(i);
    if (right < totalPages - 1) pages.push('...');
    if (totalPages > 1) pages.push(totalPages);

    return pages;
  }, [currentPage, totalPages]);

  const goTo = (p: number) => {
    const target = Math.max(1, Math.min(totalPages, p));
    if (target !== currentPage && !Number.isNaN(target)) {
      onPageChange(target);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }
  };

  // 防止 totalPages 为 0 或 NaN 时 UI 异常
  const safeTotalPages = Math.max(1, totalPages);

  return (
    <div className="pagination-wrapper">
      <div className="pagination-info" id="paginationInfo">
        共 {totalItems} 条主题 · 第 {currentPage}/{safeTotalPages} 页
      </div>

      <div className="pagination-controls" id="paginationControls">
        <button
          className={`page-btn ${currentPage === 1 ? 'disabled' : ''}`}
          onClick={() => goTo(1)}
          disabled={currentPage === 1}
        >
          «
        </button>
        <button
          className={`page-btn ${currentPage === 1 ? 'disabled' : ''}`}
          onClick={() => goTo(currentPage - 1)}
          disabled={currentPage === 1}
        >
          ‹
        </button>

        {visiblePages.map((p, i) =>
          p === '...' ? (
            <span key={`ellipsis-${i}`} className="page-ellipsis">
              …
            </span>
          ) : (
            <button
              key={p}
              className={`page-btn ${p === currentPage ? 'active' : ''}`}
              onClick={() => goTo(p as number)}
            >
              {p}
            </button>
          ),
        )}

        <button
          className={`page-btn ${currentPage === safeTotalPages ? 'disabled' : ''}`}
          onClick={() => goTo(currentPage + 1)}
          disabled={currentPage === safeTotalPages}
        >
          ›
        </button>
        <button
          className={`page-btn ${currentPage === safeTotalPages ? 'disabled' : ''}`}
          onClick={() => goTo(safeTotalPages)}
          disabled={currentPage === safeTotalPages}
        >
          »
        </button>
      </div>

      <div className="pagination-jump">
        跳至
        {/* 👇 key={currentPage} 确保切换页码时 input 自动重置为当前页 */}
        <input
          key={currentPage}
          className="jump-input"
          type="number"
          min={1}
          max={safeTotalPages}
          defaultValue={String(currentPage)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              goTo(parseInt(e.currentTarget.value, 10));
            }
          }}
        />
        <button
          className="jump-btn"
          onClick={(e) => {
            const input = e.currentTarget.previousElementSibling as HTMLInputElement;
            goTo(parseInt(input?.value, 10));
          }}
        >
          GO
        </button>
      </div>
    </div>
  );
}