// src/types/api.ts

import type { PostEntity } from './domain';

/** 通用 API 响应信封 */
export interface ApiResponse<T> {
  code: number;
  data: T;
  message: string;
}

/** 帖子详情 + 扁平回帖列表 组合响应 */
export interface PostWithRepliesResponse {
  post: PostEntity;
  replies: PostEntity[];   // 后端返回的是扁平列表，树化在前端完成
}

/** 分页列表响应 */
export interface PageResponse<T> {
  items: T[];
  totalItems: number;
  currentPage: number;
  pageSize: number;
}

/** 论坛首页聚合数据 */
export interface ForumBoardEntity {
  name: string;
  icon: string;
  count: number;
}

export interface HomePageData {
  threads: PostEntity[];
  boards: ForumBoardEntity[];
  tags: { label: string }[];
  pagination: {
    totalItems: number;
    currentPage: number;
    pageSize: number;
  };
}