// src/app/post/[id]/ReadTracker.tsx
'use client';

import { useMarkAsRead } from '@/hooks/useMarkAsRead';

export function ReadTracker({ nodeId, postId }: { nodeId: string, postId: string }) {
  useMarkAsRead({ nodeId, postId });
  return null; // 纯逻辑组件，不渲染任何 UI
}