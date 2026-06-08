// src/types/domain.ts

export interface NodeEntity {
  nodeId: number;
  name: string;
  description: string;
  topicCount: number;
}

export interface PostEntity {
  id: number;
  threadId: number;
  parentId: number;       // 0 = 主帖, >0 = 回帖
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
  toc?: string;
}

export interface UserEntity {
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
  topicCount: number;
  commentCount: number;
  passwordSet: boolean;
  status: number;
  createTime: number;
}

export interface TagEntity {
  tagId: number;
  tagName: string;
}

/** 主帖与回帖同构别名 */
export type ThreadEntity = PostEntity;

export interface PostWithThread {
  post: PostEntity;
  replies: PostEntity[];
}

export interface NodePageData {
  nodes: NodeEntity[];
  error: string | null;
}

export interface NodeDetailData {
  node: NodeEntity | null;
  error: string | null;
}