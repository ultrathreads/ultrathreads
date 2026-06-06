
export type FeConfigsType = {
  siteTitle?: string;
  siteDescription:? string;
}

export interface Reply {
  id: number;
  title: string;
  author: string;       // 注意：这里目前是 user_id 的字符串形式
  date: string;         // 注意：这里是格式化后的本地时间字符串
  category?: string;
  replies: Reply[];     // 递归嵌套结构
}

/** 主帖类型（当前业务场景下与 Reply 结构一致） */
export type Thread = Reply;

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