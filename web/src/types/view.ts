// src/types/view.ts

export interface ThreadViewItem {
  id: number;
  parentId: number;
  threadId: number;
  title: string;
  author: string;
  authorId: number;
  avatar?: string; 
  date: number;            // Unix timestamp (ms)
  lastCommentTime: number;
  isPinned?: boolean;
  nodeName?: string;
  replies?: ThreadViewItem[];
}