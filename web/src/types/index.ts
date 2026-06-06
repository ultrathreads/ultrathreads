
export interface Thread {
  id: string;
  title: string;
  author: string;
  date: string;
  category: string;
  isRead?: boolean;
  replies: Thread[];
}

export interface ForumBoard {
  name: string;
  icon: string;
  count: number;
}

export interface Tag {
  label: string;
}

export interface PageData {
  threads: Thread[];
  totalItems: number;
  currentPage: number;
  pageSize: number;
  boards: ForumBoard[];
  tags: Tag[];
}