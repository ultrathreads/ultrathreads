//src/components/BackToTop.tsx
'use client';

import { useState, useEffect, useRef, useCallback } from 'react';

const SCROLL_THRESHOLD = 300;

export default function BackToTop() {
  const [visible, setVisible] = useState(false);
  const containerRef = useRef<HTMLElement | null>(null);

  useEffect(() => {
    const el = document.getElementById('mainContent');
    if (!el) return;

    containerRef.current = el;

    const handler = () => {
      const isOverThreshold = el.scrollTop > SCROLL_THRESHOLD;
      setVisible((prev) => (prev !== isOverThreshold ? isOverThreshold : prev));
    };

    // 初始化时检查一次，防止 SSR 水合后状态不一致
    handler();

    el.addEventListener('scroll', handler, { passive: true });
    return () => {
      el.removeEventListener('scroll', handler);
      containerRef.current = null;
    };
  }, []);

  const scrollToTop = useCallback(() => {
    containerRef.current?.scrollTo({ top: 0, behavior: 'smooth' });
  }, []);

  return (
    <button
      id="back-to-top"
      className={visible ? 'visible' : ''}
      onClick={scrollToTop}
      aria-label="返回顶部"
      title="返回顶部"
      tabIndex={visible ? 0 : -1}
      aria-hidden={!visible}
    >
      ↑
    </button>
  );
}