'use client';
import { useState, useMemo } from 'react';

interface Props {
  totalItems: number;
  pageSize: number;
  currentPage: number;
}

export default function Pagination({ totalItems, pageSize, currentPage: initialPage }: Props) {
  const totalPages = Math.ceil(totalItems / pageSize);
  const [page, setPage] = useState(initialPage);
  const [jumpVal, setJumpVal] = useState(String(initialPage));

  const visiblePages = useMemo(() => {
    const pages: (number | string)[] = [];
    const d = 2;
    const l = Math.max(2, page - d);
    const r = Math.min(totalPages - 1, page + d);
    pages.push(1);
    if (l > 2) pages.push('...');
    for (let i = l; i <= r; i++) pages.push(i);
    if (r < totalPages - 1) pages.push('...');
    if (totalPages > 1) pages.push(totalPages);
    return pages;
  }, [page, totalPages]);

  const goTo = (p: number) => {
    const t = Math.max(1, Math.min(totalPages, p));
    if (t !== page) {
      setPage(t);
      setJumpVal(String(t));
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }
  };

  return (
    <div className="pagination-wrapper">
      <div className="pagination-info" id="paginationInfo">
        共 {totalItems} 条主题 · 第 {page}/{totalPages} 页
      </div>
      <div className="pagination-controls" id="paginationControls">
        <button className={`page-btn ${page === 1 ? 'disabled' : ''}`} onClick={() => goTo(1)}>«</button>
        <button className={`page-btn ${page === 1 ? 'disabled' : ''}`} onClick={() => goTo(page - 1)}>‹</button>
        {visiblePages.map((p, i) =>
          p === '...' ? (
            <span key={`e${i}`} className="page-ellipsis">…</span>
          ) : (
            <button
              key={p}
              className={`page-btn ${p === page ? 'active' : ''}`}
              onClick={() => goTo(p as number)}
            >
              {p}
            </button>
          )
        )}
        <button className={`page-btn ${page === totalPages ? 'disabled' : ''}`} onClick={() => goTo(page + 1)}>›</button>
        <button className={`page-btn ${page === totalPages ? 'disabled' : ''}`} onClick={() => goTo(totalPages)}>»</button>
      </div>
      <div className="pagination-jump">
        跳至
        <input
          className="jump-input"
          type="number"
          min={1}
          max={totalPages}
          value={jumpVal}
          onChange={(e) => setJumpVal(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && goTo(parseInt(jumpVal))}
        />
        <button className="jump-btn" onClick={() => goTo(parseInt(jumpVal))}>GO</button>
      </div>
    </div>
  );
}