// src/lib/utils/thread-adapter.ts
import type { ThreadListItem } from '@/services/thread-service';
import type { PostEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';

type ThreadSource = ThreadListItem | PostEntity;

/**
 * 统一适配器：将帖子传输模型转换为 ThreadViewItem
 * 通过 'threadId' 字段区分列表项与详情回帖
 */
export function adaptToThreadView(source: ThreadSource): ThreadViewItem {
  const isListItem = 'threadId' in source;

  // ✅ 安全提取用户名，防止 user 对象缺失导致详情页崩溃
  const author = isListItem
    ? (source as ThreadListItem).user?.nickname || (source as ThreadListItem).user?.username
    : (source as PostEntity).user?.nickname || (source as PostEntity).user?.username;

  return {
    id: source.id,
    parentId: isListItem
      ? (source as ThreadListItem).parentId
      : ((source as PostEntity).parentId ?? 0), // ✅ 兜底为 0，避免树构建时误判
    title: source.title,
    author: author ?? '匿名用户', // ✅ 最终兜底
    date: isListItem
      ? (source as ThreadListItem).lastCommentTime || (source as ThreadListItem).createTime
      : (source as PostEntity).createTime,
    nodeName: isListItem
      ? (source as ThreadListItem).node?.name
      : (source as PostEntity).node?.name,
  };
}