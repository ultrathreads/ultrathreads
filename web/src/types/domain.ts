// src/types/domain.ts

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
  postCount: number;
  commentCount: number;
  passwordSet: boolean;
  status: number;
  createTime: number;
}

export interface NodeEntity {
  nodeId: number;
  name: string;
  description: string;
  postCount: number;
}

export interface TagEntity {
  label: string;
}

/** 主帖与回帖同构别名 */
export type ThreadEntity = PostEntity;