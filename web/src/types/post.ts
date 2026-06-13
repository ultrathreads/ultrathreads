// types/post.ts

/** 创建根帖请求参数 */
export interface CreateRootPostPayload {
  nodeSlug: string;
  title: string;
  content: string;
  tags?: string[];
  imageList?: string[];
}

export interface UpdateRootPostPayload {
  nodeSlug: string;
  title: string;
  content: string;
  tags?: string[];
}

/** 创建回复请求参数（不再包含 parentSlug，由 URL 路径承载） */
export interface CreateReplyPayload {
  content: string;
  imageList?: string[];
}
export interface UpdateReplyPayload {
  content: string;
}

export interface CreatePostResponse {
  slug: string;
}