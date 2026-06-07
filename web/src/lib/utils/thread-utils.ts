// src/lib/utils/thread-utils.ts
import type { ThreadViewItem } from '@/types/view';

/** 带子节点的扩展类型，用于内部树形构建 */
type ThreadTreeNode = ThreadViewItem & { replies: ThreadTreeNode[] };

/**
 * 将已适配的扁平 ThreadViewItem 列表转换为嵌套树结构
 * ⚠️ 注意：此函数仅负责构建父子关系，不做任何字段映射
 * 字段映射应在服务端或数据获取层通过 adaptToThreadView 完成
 */
export function buildThreadTree(posts: ThreadViewItem[]): ThreadTreeNode[] {
  const map = new Map<number, ThreadTreeNode>();
  const roots: ThreadTreeNode[] = [];

  // 1. 初始化所有节点（保留原始字段，仅追加空的 replies 数组）
  for (const post of posts) {
    map.set(post.id, { ...post, replies: [] });
  }

  // 2. 根据 parentId 挂载子节点
  for (const post of posts) {
    const node = map.get(post.id)!;
    const isRoot = !post.parentId || post.parentId <= 0;

    if (isRoot) {
      roots.push(node);
    } else {
      const parent = map.get(post.parentId);
      if (parent) {
        parent.replies.push(node);
      } else {
        // 父节点缺失时降级为根节点，避免丢失数据
        console.warn(`[buildThreadTree] Parent ${post.parentId} not found for post ${post.id}`);
        roots.push(node);
      }
    }
  }

  return roots;
}