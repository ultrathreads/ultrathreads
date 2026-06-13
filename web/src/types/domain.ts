// src/types/domain.ts

export interface NodeEntity {
  slug: string;
  name: string;
  description: string;
  topicCount: number;
}

export interface PostEntity {
  slug: string;
  threadSlug: string;
  parentSlug: string;
  isRoot: bool;
  type: number;
  user: UserEntity;
  node: NodeEntity;
  tags: TagEntity[] | null;
  title: string;
  imageList: string[] | null;
  lastCommentUser: UserEntity | null;
  lastCommentTime: number;
  viewCount: number;
  commentCount: number;
  likeCount: number;
  createTime: number;
  content?: string;
  rawContent?: string;
  toc?: string;
}

export interface UserEntity {
  slug: string;
  username: string;
  email: string;
  nickname: string;
  avatar: string;
  level: number;
  levelName: string;
  website: string;
  description: string;
  score: number;
  topicCount: number;
  commentCount: number;
  passwordSet: boolean;
  status: number;
  createTime: number;
}

export interface TagEntity {
  slug: string;
  tagName: string;
}

/** 主帖与回帖同构别名 */
export type ThreadEntity = PostEntity;

export interface PostWithThread {
  post: PostEntity;
  replies: PostEntity[];
}

export interface PostWithFlat {
  posts: PostEntity[];
}

export interface NodePageData {
  nodes: NodeEntity[];
  error: string | null;
}

export interface NodeDetailData {
  node: NodeEntity | null;
  error: string | null;
}

export interface TagDetailData {
  tag: TagEntity | null;
  error: string | null;
}