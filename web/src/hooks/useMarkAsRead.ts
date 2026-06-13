// src/hooks/useMarkAsRead.ts
'use client';

import { useEffect, useRef } from 'react';

interface UseMarkAsReadOptions {
  postSlug: string;
  nodeSlug?: string;
  minDuration?: number;      // 默认 3000ms
  scrollThreshold?: number;  // 默认 0.3
  dedupeInterval?: number;   // 去重时间窗口，默认 5分钟 (300000ms)
}

// 🚀 内存级 LRU 缓存，避免同会话内频繁读取 sessionStorage
const memoryCache = new Map<string, number>();

function getCacheKey(postSlug: string): string {
  return `read:${postSlug}`;
}

function hasRecentlyReported(key: string, interval: number): boolean {
  // 1. 优先查内存
  const memTs = memoryCache.get(key);
  if (memTs && Date.now() - memTs < interval) return true;

  // 2. 回退查 sessionStorage (处理刷新/重新挂载场景)
  try {
    const ssTs = sessionStorage.getItem(key);
    if (ssTs) {
      const ts = parseInt(ssTs, 10);
      if (Date.now() - ts < interval) {
        memoryCache.set(key, ts); // 回填内存
        return true;
      }
    }
  } catch { /* SSR 或隐私模式可能报错，静默忽略 */ }

  return false;
}

function markReported(key: string): void {
  const now = Date.now();
  memoryCache.set(key, now);
  try {
    sessionStorage.setItem(key, String(now));
  } catch { /* 配额满或不可用，仅依赖内存缓存 */ }
}

export function useMarkAsRead({
  postSlug,
  nodeSlug,
  minDuration = 3000,
  scrollThreshold = 0.3,
  dedupeInterval = 5 * 60 * 1000, // 默认5分钟内同一帖子不重复上报
}: UseMarkAsReadOptions) {
  const hasReported = useRef(false);
  const enterTime = useRef(0);
  const timerRef = useRef<ReturnType<typeof setTimeout>>();

  useEffect(() => {
    // 在 effect 初始化时立即检查持久化缓存
    const cacheKey = getCacheKey(postSlug);
    if (hasRecentlyReported(cacheKey, dedupeInterval)) {
      hasReported.current = true; // 直接标记为已上报，跳过所有监听
      return;
    }

    const report = () => {
      if (hasReported.current) return;
      hasReported.current = true;

      // 写入持久化缓存
      markReported(cacheKey);

      navigator.sendBeacon(`/api/posts/${postSlug}/view-post?nodeSlug=${nodeSlug || ''}`);

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
  }, [postSlug, nodeSlug, minDuration, scrollThreshold, dedupeInterval]);
}