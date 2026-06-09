// src/types/api.ts

import type { PostEntity } from './domain';

export interface ApiResponse<T> {
  code: number;
  message: string;
  success: boolean;
  data: T;
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

export interface PaginationMeta {
  /** 总记录数 */
  totalItems: number;
  /** 当前页码（从 1 开始） */
  currentPage: number;
  /** 每页条数 */
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