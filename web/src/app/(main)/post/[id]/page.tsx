// src/app/post/[id]/page.tsx
import { getPostWithThread } from '@/services/post-service';
import PostDetailCard from '@/components/PostDetailCard';
import ThreadItem from '@/components/features/ThreadItem';
import type { PostDetail } from '@/types/common';
import type { Reply } from '@/types/reply';

interface Props {
  params: Promise<{ id: string }>;
}

/** 带子节点的扩展类型，用于内部树形构建 */
interface ThreadNode extends PostDetail {
  children: ThreadNode[];
}

/** 
 * 将扁平回帖列表构建为树形结构 (O(n) 复杂度)
 */
function buildThreadTree(replies: PostDetail[]): ThreadNode[] {
  const nodeMap = new Map<number, ThreadNode>();
  const roots: ThreadNode[] = [];

  for (const reply of replies) {
    nodeMap.set(reply.id, { ...reply, children: [] });
  }

  for (const reply of replies) {
    const node = nodeMap.get(reply.id)!;
    if (reply.parentId === 0) {
      roots.push(node);
    } else {
      const parent = nodeMap.get(reply.parentId);
      if (parent) {
        parent.children.push(node);
      } else {
        console.warn(`[buildThreadTree] Parent ${reply.parentId} not found for reply ${reply.id}`);
        roots.push(node);
      }
    }
  }
  return roots;
}

/**
 * ✅ 核心适配器：将新的 ThreadNode 转换为 ThreadItem 期望的 Reply 格式
 * 隔离新旧数据结构，保证 ThreadItem 及其复用方零修改
 */
function adaptToReply(node: ThreadNode): Reply {
  return {
    id: node.id,
    title: node.title,
    author: node.user.nickname,
    date: node.createTime,
    category: node.node?.name,
    // 递归转换子节点
    replies: node.children.length > 0 
      ? node.children.map(adaptToReply) 
      : undefined,
  };
}

export default async function ReadPage({ params }: Props) {
  const { id } = await params;
  const { post, replies } = await getPostWithThread(id);
  
  // 1. 构建树形结构
  const treeNodes = buildThreadTree(replies);
  
  // 2. ✅ 转换为 ThreadItem 兼容的数据格式
  const adaptedReplies = treeNodes.map(adaptToReply);

  return (
    <div className="main-body">
      <div className="detail-back-bar">
        <a className="back-list-btn" href="/">← 返回列表</a>
      </div>

      <PostDetailCard post={post} />

      <div className="thread-tree-container">
        <div className="thread-tree-header">
          <div className="thread-tree-title">💬 回帖讨论 ({post.commentCount})</div>
          <div className="thread-tree-actions">
            <select className="sort-select" aria-label="回帖排序" defaultValue="oldest">
              <option value="oldest">最早回复</option>
              <option value="newest">最新回复</option>
              <option value="hot">最热回复</option>
            </select>
          </div>
        </div>

        <ul className="thread">
          {/* ✅ 传入转换后的数据，isRoot 保持与原组件契约一致 */}
          {adaptedReplies.map((reply) => (
            <ThreadItem 
              key={reply.id} 
              item={reply} 
              isRoot 
              currentPostId={post.id}
            />
          ))}
        </ul>
      </div>
    </div>
  );
}