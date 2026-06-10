// src/hooks/useMarkAsRead.ts
'use client';

import { useEffect, useRef } from 'react';

interface UseMarkAsReadOptions {
  nodeId: string;
  minDuration?: number;   // 默认 3000ms
  scrollThreshold?: number; // 默认 0.3
}

export function useMarkAsRead({
  nodeId,
  minDuration = 3000,
  scrollThreshold = 0.3,
}: UseMarkAsReadOptions) {
  const hasReported = useRef(false);
  const enterTime = useRef(0);
  const timerRef = useRef<ReturnType<typeof setTimeout>>();

  useEffect(() => {
    if (!nodeId || hasReported.current) return;

    const report = () => {
      if (hasReported.current) return;
      hasReported.current = true;

      // ✅ sendBeacon 保证页面跳转/关闭时请求不丢失
      // 与你后端 EventBus 异步 ViewPost 完美匹配
      navigator.sendBeacon(`/api/nodes/${nodeId}/view-post`);

      if (timerRef.current) clearTimeout(timerRef.current);
    };

    // 条件1: 有效停留 + 页面可见
    enterTime.current = Date.now();
    timerRef.current = setTimeout(() => {
      if (document.visibilityState === 'visible') report();
    }, minDuration);

    // 条件2: 滚动深度达标
    const handleScroll = () => {
      if (hasReported.current) return;
      const docHeight = document.documentElement.scrollHeight - window.innerHeight;
      if (docHeight > 0 && window.scrollY / docHeight >= scrollThreshold) {
        report();
      }
    };

    // 补充: 页面隐藏时检查是否已满足时长
    const handleVisibilityChange = () => {
      if (document.visibilityState === 'hidden' && !hasReported.current) {
        if (Date.now() - enterTime.current >= minDuration) report();
      }
    };

    window.addEventListener('scroll', handleScroll, { passive: true });
    document.addEventListener('visibilitychange', handleVisibilityChange);

    return () => {
      window.removeEventListener('scroll', handleScroll);
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      if (timerRef.current) clearTimeout(timerRef.current);
    };
  }, [nodeId, minDuration, scrollThreshold]);
}