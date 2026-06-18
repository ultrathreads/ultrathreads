// src/types/view.ts

export interface ThreadViewItem {
  id: number;
  slug: string;
  parentSlug: string;
  threadSlug: string;
  title: string;
  author: string;
  authorSlug: string;
  avatar?: string; 
  date: string;            // ISO 8601 时间字符串
  lastRepliedAt: string;
  isPinned?: boolean;
  node?: { slug: string; name: string; };
  tags?: { slug: string; name: string }[];
  replies?: ThreadViewItem[];
}

/** 从列表页透传的回溯状态 */
export interface BackState {
  nodeSlug?: string;
  tagSlug?: string;
  page?: string;
}
