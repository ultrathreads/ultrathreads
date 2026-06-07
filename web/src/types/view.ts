// src/types/view.ts

export interface ThreadViewItem {
  id: number;
  title: string;
  author: string;
  date: number;            // Unix timestamp (ms)
  category?: string;
  replies?: ThreadViewItem[];
}