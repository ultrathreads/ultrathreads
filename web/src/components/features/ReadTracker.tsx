// src/app/threads/[id]/ReadTracker.tsx
'use client';

import { useMarkAsRead } from '@/hooks/useMarkAsRead';

export function ReadTracker({ postSlug, nodeSlug }: { postSlug: string, nodeSlug: string }) {
  useMarkAsRead({ postSlug, nodeSlug });
  return null; // 纯逻辑组件，不渲染任何 UI
}