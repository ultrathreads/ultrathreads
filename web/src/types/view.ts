// src/types/view.ts

export interface ThreadViewItem {
  id: number;
  parentId: number;
  threadId: number;
  title: string;
  author: string;
  date: number;            // Unix timestamp (ms)
  lastCommentTime: number;
  nodeName?: string;
  replies?: ThreadViewItem[];
}