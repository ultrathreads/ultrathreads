// src/lib/utils/assemble-sideload.ts

import type { ThreadListItem } from '@/services/thread-service';

interface SideloadUser { slug: string; username: string; nickname: string; avatar: string }
interface SideloadNode { slug: string; name: string }

export interface IncludedData {
  users?: SideloadUser[];
  nodes?: SideloadNode[];
}

const FALLBACK_USER: ThreadListItem['user'] = {
  slug: '', username: '未知用户', nickname: '未知用户', avatar: '',
};
const FALLBACK_NODE: ThreadListItem['node'] = {
  slug: '', name: '未知板块',
};

/**
 * 将 sideload 数据注入到列表项中，兼容新旧 API 格式
 * 输入输出均为 ThreadListItem，无需额外的 ViewModel 类型
 */
export function assembleSideload(
  posts: ThreadListItem[],
  included?: IncludedData,
): ThreadListItem[] {
  const userMap = new Map<string, SideloadUser>();
  const nodeMap = new Map<string, SideloadNode>();

  for (const u of included?.users ?? []) userMap.set(u.slug, u);
  for (const n of included?.nodes ?? []) nodeMap.set(n.slug, n);

  return posts.map((post) => {
    // --- User 解析（userSlug 优先 → fallback post.user → 兜底值）---
    let user: NonNullable<ThreadListItem['user']>;
    if (post.userSlug && userMap.has(post.userSlug)) {
      const raw = userMap.get(post.userSlug)!;
      user = {
        slug: raw.slug,
        username: raw.username,
        nickname: raw.nickname,
        avatar: raw.avatar,
      };
    } else if (post.user) {
      user = post.user;
    } else {
      user = FALLBACK_USER;
    }

    // --- Node 解析 ---
    let node: NonNullable<ThreadListItem['node']>;
    if (post.nodeSlug && nodeMap.has(post.nodeSlug)) {
      const raw = nodeMap.get(post.nodeSlug)!;
      node = { slug: raw.slug, name: raw.name };
    } else if (post.node) {
      node = post.node;
    } else {
      node = FALLBACK_NODE;
    }

    return { ...post, user, node };
  });
}