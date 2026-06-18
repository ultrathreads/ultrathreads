import type { PostEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';
import { adaptToThreadView } from './thread-adapter';

/** 帖子在树形结构中的包装项，直接继承帖子实体字段 */
interface PostTreeItem extends PostEntity {
  children: PostTreeItem[];
}

interface BuildTreeOptions {
  lastReadAtMap?: Record<string, number>;
}

/**
 * 通用帖子树形构建器
 * - 根帖顺序保持传入顺序（由后端控制）
 * - 所有层级回帖按时间正序排列（最早在前）
 */
export function buildThreadTree(
  posts: PostEntity[],
  options?: BuildTreeOptions,
): ThreadViewItem[] {
  const { lastReadAtMap } = options ?? {};

  // 1. O(n) 建立帖子索引
  const postIndex = new Map<string, PostTreeItem>();
  const roots: PostTreeItem[] = [];

  for (const post of posts) {
    postIndex.set(post.slug, { ...post, children: [] });
  }

  // 2. 挂载子节点（保持原始相对顺序，后续统一排序）
  for (const post of posts) {
    const currentPost = postIndex.get(post.slug)!;
    const isRoot = post.isRoot;

    if (isRoot) {
      roots.push(currentPost);
    } else {
      const parentPost = postIndex.get(post.parentSlug);
      if (parentPost) {
        parentPost.children.push(currentPost);
      } else {
        // 父节点缺失时降级为根节点，避免丢失数据
        roots.push(currentPost);
      }
    }
  }

  // 3. 递归转换为视图模型，每层子节点固定按时间正序
  function toThreadView(treeItem: PostTreeItem): ThreadViewItem {
    treeItem.children.sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime());

    const lookupKey = String(treeItem.nodeSlug);

    const lastReadAt =
    lastReadAtMap !== undefined ? lastReadAtMap[lookupKey] : undefined;

    const base = adaptToThreadView(treeItem, { lastReadAt });

    return {
      ...base,
      replies:
        treeItem.children.length > 0
          ? treeItem.children.map(toThreadView)
          : undefined,
    };
  }

  return roots.map(toThreadView);
}