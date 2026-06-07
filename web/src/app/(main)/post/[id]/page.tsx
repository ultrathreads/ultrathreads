// src/app/post/[id]/page.tsx
import { notFound } from 'next/navigation';
import { getPostWithThread } from '@/services/post-service';
import PostDetailCard from '@/components/PostDetailCard';
import ThreadItem from '@/components/features/ThreadItem';
import type { PostEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';
import { adaptToThreadView } from '@/lib/utils/thread-adapter';

// ✅ 强制动态渲染，避免构建时因缺少静态参数导致 404
export const dynamic = 'force-dynamic';

interface Props {
  params: Promise<{ id: string }>;
}

/** 带子节点的扩展类型，用于内部树形构建 */
interface ThreadNode extends PostEntity {
  children: ThreadNode[];
}

/**
 * 将扁平回帖列表构建为树形结构 (O(n) 复杂度)
 */
function buildThreadTree(replies: PostEntity[]): ThreadNode[] {
  const nodeMap = new Map<number, ThreadNode>();
  const roots: ThreadNode[] = [];

  for (const reply of replies) {
    nodeMap.set(reply.id, { ...reply, children: [] });
  }

  for (const reply of replies) {
    const node = nodeMap.get(reply.id)!;
    // ✅ 兼容 parentId 为 null / undefined / 0 的情况
    const isRoot = !reply.parentId || reply.parentId <= 0;

    if (isRoot) {
      roots.push(node);
    } else {
      const parent = nodeMap.get(reply.parentId);
      if (parent) {
        parent.children.push(node);
      } else {
        console.warn(`[buildThreadTree] Parent ${reply.parentId} not found for reply ${reply.id}`);
        roots.push(node); // 父节点缺失时降级为根节点，避免丢失
      }
    }
  }

  return roots;
}

/**
 * 递归适配器：将树形 ThreadNode 转换为 ThreadViewItem
 * 复用共享 adaptToThreadView 处理单节点字段映射
 */
function adaptTreeNode(node: ThreadNode): ThreadViewItem {
  const base = adaptToThreadView(node);
  return {
    ...base,
    replies: node.children.length > 0
      ? node.children.map(adaptTreeNode)
      : undefined,
  };
}

export default async function ReadPage({ params }: Props) {
  // ✅ 安全解构 params，防止 Next.js 版本差异或路由异常导致崩溃
  let id: string;
  try {
    const resolved = await params;
    id = resolved.id;
  } catch {
    notFound();
  }

  // ✅ 防御性数据获取，避免接口异常直接导致页面 500/404
  let post: PostEntity | null = null;
  let replies: PostEntity[] = [];

  try {
    const result = await getPostWithThread(id);
    post = result.post;
    replies = result.replies ?? [];
  } catch (error) {
    console.error(`[ReadPage] Failed to fetch post ${id}:`, error);
  }

  // ✅ 帖子不存在时显式触发 Next.js 404 页面
  if (!post) {
    notFound();
  }

  // 构建树形结构并适配为视图模型
  const treeNodes = buildThreadTree(replies);
  const adaptedReplies = treeNodes.map(adaptTreeNode);

  return (
    <div className="main-body">
      <div className="detail-back-bar">
        <a className="back-list-btn" href="/">← 返回列表</a>
      </div>

      <PostDetailCard post={post} />

      <div className="thread-tree-container">
        <div className="thread-tree-header">
          <div className="thread-tree-title">
            💬 回帖讨论 ({post.commentCount ?? 0})
          </div>
          <div className="thread-tree-actions">
            <select className="sort-select" aria-label="回帖排序" defaultValue="oldest">
              <option value="oldest">最早回复</option>
              <option value="newest">最新回复</option>
              <option value="hot">最热回复</option>
            </select>
          </div>
        </div>

        <ul className="thread">
          {adaptedReplies.map((reply) => (
            <ThreadItem
              key={reply.id}
              item={reply}
              isRoot
              currentPostId={String(post!.id)} // ✅ 统一转字符串，确保与 ThreadItem 内部比较一致
            />
          ))}
        </ul>
      </div>
    </div>
  );
}