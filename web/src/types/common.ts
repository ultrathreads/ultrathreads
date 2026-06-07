// src/types/common.ts

/**
 * ✅ 统一的帖子/回帖结构（对应后端 SimplePost）
 * 主帖与回帖结构完全一致，通过 parentId 区分层级
 */
export interface PostDetail {
  id: number;
  threadId: number;
  parentId: number;
  type: number;
  user: PostUser;
  node: PostNode;
  tags: Tag[] | null;
  title: string;
  imageList: string[] | null;
  lastCommentUser: PostUser | null;
  lastCommentTime: number;
  viewCount: number;
  commentCount: number;
  likeCount: number;
  createTime: number;
  content?: string;
  toc?: string;
}

/** 主帖与回帖同构别名 */
export type Thread = PostDetail;

export interface PostUser {
  id: number;
  username: string;
  email: string;
  nickname: string;
  avatar: string;
  level: number;
  levelName: string;
  website: string;
  description: string;
  score: number;
  postCount: number;
  commentCount: number;
  passwordSet: boolean;
  status: number;
  createTime: number;
}

export interface PostNode {
  nodeId: number;
  name: string;
  description: string;
  postCount: number;
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
  threads: PostDetail[];
  totalItems: number;
  currentPage: number;
  pageSize: number;
  boards: ForumBoard[];
  tags: Tag[];
}

/** 帖子详情 + 扁平回帖列表 组合响应 */
export interface PostWithThread {
  post: PostDetail;
  replies: PostDetail[];
}
