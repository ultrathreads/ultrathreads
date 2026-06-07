// src/types/adapters.ts

import type { PostEntity } from './domain';
import type { ThreadViewItem } from './view';

/**
 * 将后端扁平回帖列表转换为前端树形视图结构
 * @param rootPost  当前主帖
 * @param flatReplies 后端返回的扁平回帖数组
 */
export function buildThreadTree(
  rootPost: PostEntity,
  flatReplies: PostEntity[]
): ThreadViewItem {
  const map = new Map<number, ThreadViewItem>();

  // 先创建所有节点的视图对象
  const toView = (p: PostEntity): ThreadViewItem => ({
    id: p.id,
    title: p.title,
    author: p.user.nickname,
    date: p.createTime,
    category: p.node?.name,
    replies: [],
  });

  const rootView = toView(rootPost);
  map.set(rootPost.id, rootView);

  for (const reply of flatReplies) {
    map.set(reply.id, toView(reply));
  }

  // 挂载父子关系
  for (const reply of flatReplies) {
    const node = map.get(reply.id)!;
    const parent = map.get(reply.parentId);
    if (parent) {
      parent.replies!.push(node);
    } else {
      // 孤儿节点兜底挂到根节点下
      rootView.replies!.push(node);
    }
  }

  return rootView;
}

/**
 * 将首页帖子列表批量转换为视图模型
 */
export function adaptThreadsToView(posts: PostEntity[]): ThreadViewItem[] {
  return posts.map((p) => ({
    id: p.id,
    title: p.title,
    author: p.user.nickname,
    date: p.createTime,
    category: p.node?.name,
  }));
}