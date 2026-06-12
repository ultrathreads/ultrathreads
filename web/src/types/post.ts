// types/post.ts

/** 创建根帖请求参数 */
export interface CreateRootPostPayload {
  nodeSlug: string;       // ✅ 根帖必填
  title: string;
  content: string;
  tags?: string[];
  imageList?: string[];
}

/** 创建回复请求参数（不再包含 parentSlug，由 URL 路径承载） */
export interface CreateReplyPayload {
  content: string;
  imageList?: string[];
}

// 响应类型保持不变，或按需拆分
export type CreatePostResponse = SimplePost;