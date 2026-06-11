// src/lib/utils/thread-adapter.ts
import type { ThreadListItem } from '@/services/thread-service';
import type { PostEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';

type ThreadSource = ThreadListItem | PostEntity;

interface AdaptOptions {
  /** 用户对该节点/版块的最后阅读时间戳 (ms) */
  lastReadAt?: number;
}

/**
 * 统一适配器：将帖子传输模型转换为 ThreadViewItem
 * 通过 'type' 字段区分列表项与详情回帖（两者均有 threadId，不能用 threadId 判断）
 */
export function adaptToThreadView(
  source: ThreadSource,
  options: AdaptOptions = {}
): ThreadViewItem {
  // ✅ PostEntity 有 type 字段，ThreadListItem 没有
  const isPostEntity = 'type' in source;

  const post = source as PostEntity;
  const listItem = source as ThreadListItem;

  // ✅ 安全提取用户名，防止 user 对象缺失导致详情页崩溃
  const author = isPostEntity
    ? (post.user?.nickname || post.user?.username)
    : (listItem.user?.nickname || listItem.user?.username);

  const authorId = isPostEntity
    ? (post.user?.id || post.user?.id)
    : (listItem.user?.id || listItem.user?.id);

  const avatar = isPostEntity
      ? (post.user?.avatar || undefined)
      : (listItem.user?.avatar || undefined);

  // ✅ 提取用于判断已读的基准时间（优先用 updateTime，回退到 createTime）
  let contentTimestamp: number | string | undefined;
  if (isPostEntity) {
    contentTimestamp = post.updateTime ?? post.createTime;
  } else {
    // 列表项：如果有最后评论时间且业务上"新回复=未读"，则用它；否则用 createTime
    contentTimestamp = listItem.lastCommentTime || listItem.createTime;
  }

  const contentTime = typeof contentTimestamp === 'number'
    ? contentTimestamp
    : new Date(contentTimestamp ?? 0).getTime();

  // ✅ 核心已读判断
  // 没有 lastReadAt 时默认为未读（安全兜底，避免误标已读）
  const isRead = options.lastReadAt != null && Number.isFinite(contentTime)
    ? contentTime <= options.lastReadAt
    : false;

  return {
    id: source.id,
    threadId: source.threadId,
    parentId: isPostEntity
      ? (post.parentId ?? 0)
      : listItem.parentId,
    title: source.title,
    isPinned: source.isPinned,
    isRead,
    author: author ?? '匿名用户',
    authorId,
    avatar,
    date: Number.isFinite(contentTimestamp) ? contentTimestamp : 0,
    nodeName: isPostEntity
      ? post.node?.name
      : listItem.node?.name,
  };
}