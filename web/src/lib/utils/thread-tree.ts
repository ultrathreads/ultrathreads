import type { PostEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';
import { adaptToThreadView } from './thread-adapter';

interface ThreadNode extends PostEntity {
  children: ThreadNode[];
}

interface BuildThreadTreeOptions {
  /** 用户对该节点/版块的最后阅读时间戳 (ms) */
  lastReadAt?: number;
}

/**
 * 通用帖子树形构建器
 * - 根帖顺序保持传入顺序（由后端控制）
 * - 所有层级回帖按时间正序排列（最早在前）
 */
export function buildThreadTree(
  posts: PostEntity[],
  options: BuildThreadTreeOptions = {} // ✅ 新增可选参数
): ThreadViewItem[] {
  const { lastReadAt } = options;

  // 1. O(n) 建图
  const nodeMap = new Map<number, ThreadNode>();
  const roots: ThreadNode[] = [];

  for (const post of posts) {
    nodeMap.set(post.id, { ...post, children: [] });
  }

  // 2. 挂载子节点（保持原始相对顺序，后续统一排序）
  for (const post of posts) {
    const node = nodeMap.get(post.id)!;
    const isRoot = !post.parentId || post.parentId <= 0;

    if (isRoot) {
      roots.push(node);
    } else {
      const parent = nodeMap.get(post.parentId);
      if (parent) {
        parent.children.push(node);
      } else {
        console.warn(
          `[buildThreadTree] Parent ${post.parentId} not found for post ${post.id}`
        );
        roots.push(node);
      }
    }
  }

  // 3. 递归适配视图，每层子节点固定按时间正序
  function adaptTreeNode(node: ThreadNode): ThreadViewItem {
    node.children.sort(
      (a, b) =>
        new Date(a.createTime).getTime() - new Date(b.createTime).getTime()
    );

    // ✅ 将 lastReadAt 透传给适配器
    const base = adaptToThreadView(node, { lastReadAt });
    return {
      ...base,
      replies:
        node.children.length > 0
          ? node.children.map(adaptTreeNode)
          : undefined,
    };
  }

  return roots.map(adaptTreeNode);
}