// src/types/reply.ts (或你项目中 ThreadItem 实际引用的类型文件)
export interface Reply {
  id: number;
  title: string;
  author: string;
  date: number;
  category?: string;
  replies?: Reply[];
}