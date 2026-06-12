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
  date: number;            // Unix timestamp (ms)
  lastCommentTime: number;
  isPinned?: boolean;
  nodeName?: string;
  replies?: ThreadViewItem[];
}