// lib/thread-utils.ts
import type { SimplePost } from '@/lib/api/posts';
import type { Thread, Reply } from '@/types';
import { formatTimestamp } from '@/lib/time-utils';

/**
 * 将后端扁平 SimplePost 列表转换为 ThreadTree 所需的嵌套树结构
 * - 主帖：parent_id === 0
 * - 回帖：按 parent_id 挂载到对应父节点下
 * - 字段映射：SimplePost → Thread / Reply
 */
export function buildThreadTree(posts: SimplePost[]): Thread[] {
  const map = new Map<number, Reply & { replies: Reply[] }>();
  const roots: Thread[] = [];

  for (const post of posts) {
    map.set(post.id, {
      id: post.id,
      title: post.title,
      author: post.user.nickname || post.user.username || '匿名用户',
      // ✅ 使用驼峰，且保留 * 1000 的时间戳转换
      date: formatTimestamp(post.createTime), 
      category: undefined,
      replies: [],
    });
  }

  for (const post of posts) {
    const node = map.get(post.id)!;
    // ✅ 使用驼峰，并保留容错判断
    const isRoot = !post.parentId || post.parentId <= 0; 
    
    if (isRoot) {
      roots.push(node as Thread);
    } else {
      const parent = map.get(post.parentId); // ✅ 使用驼峰
      if (parent) {
        parent.replies.push(node);
      }
    }
  }

  return roots;
}