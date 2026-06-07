// src/lib/utils/thread-adapter.ts
import type { PostDetail, Reply } from '@/types/common';

/**
 * 将真实 API 的 PostDetail / CommentNode 转换为 ThreadItem 所需的 Reply 结构
 * @param post - 主帖或评论节点
 * @param activeId - 当前 URL 选中的 ID，用于标记高亮
 * @param children - 已转换的子节点数组
 */
export function toReply(
  post: PostDetail & { children?: (PostDetail & { children?: any[] })[] },
  activeId: number,
): Reply {
  const isActive = post.id === activeId;

  // ✅ 递归转换子节点，保持树形结构
  const replies: Reply[] = (post.children || []).map((child) =>
    toReply(child, activeId)
  );

  return {
    id: post.id,
    title: post.title || '',           // 评论可能无标题，兜底空串
    author: String(post.user.id),      // ⚠️ 按现有约定，author 存 user_id 字符串
    date: new Date(post.createTime).toLocaleString('zh-CN'), // 转为本地时间字符串
    category: post.node?.name,
    replies,
    // 💡 扩展字段：通过类型断言或扩展 Reply 接口传递高亮状态
    // 若不想污染 Reply 类型，可在 ThreadItem 中通过 props 单独传入 activeId
    _isActive: isActive,
  } as Reply & { _isActive: boolean };
}