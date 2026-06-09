// src/lib/utils/thread-adapter.ts
import type { ThreadListItem } from '@/services/thread-service';
import type { PostEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';

type ThreadSource = ThreadListItem | PostEntity;

/**
 * 统一适配器：将帖子传输模型转换为 ThreadViewItem
 * 通过 'type' 字段区分列表项与详情回帖（两者均有 threadId，不能用 threadId 判断）
 */
export function adaptToThreadView(source: ThreadSource): ThreadViewItem {
  // ✅ PostEntity 有 type 字段，ThreadListItem 没有
  const isPostEntity = 'type' in source;

  const post = source as PostEntity;
  const listItem = source as ThreadListItem;

  // ✅ 安全提取用户名，防止 user 对象缺失导致详情页崩溃
  const author = isPostEntity
    ? (post.user?.nickname || post.user?.username)
    : (listItem.user?.nickname || listItem.user?.username);

  const avatar = isPostEntity
      ? (post.user?.avatar || undefined)
      : (listItem.user?.avatar || undefined);


  // ✅ 统一提取时间戳，并确保输出严格为 number (ms)
  let rawTimestamp: number | string | undefined;
  if (isPostEntity) {
    rawTimestamp = post.createTime;
  } else {
    rawTimestamp = listItem.lastCommentTime || listItem.createTime;
  }

  // 防御性转换：兼容后端可能返回字符串时间戳或非法值的情况
  const timestamp = typeof rawTimestamp === 'number'
    ? rawTimestamp
    : new Date(rawTimestamp ?? 0).getTime();

  return {
    id: source.id,
    threadId: source.threadId, // ✅ 补全必填字段
    parentId: isPostEntity
      ? (post.parentId ?? 0)   // ✅ 兜底为 0，避免树构建时误判为根节点
      : listItem.parentId,
    title: source.title,
    author: author ?? '匿名用户', // ✅ 最终兜底，保证一定有值
    avatar,
    date: Number.isFinite(timestamp) ? timestamp : 0, // ✅ 防止 NaN 污染视图层
    nodeName: isPostEntity
      ? post.node?.name
      : listItem.node?.name,
  };
}