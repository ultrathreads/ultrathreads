// src/app/post/[id]/page.tsx
import { notFound } from 'next/navigation';
import Link from 'next/link';
import { getPostWithThread } from '@/services/post-service';
import PostDetailCard from '@/components/PostDetailCard';
import ThreadItem from '@/components/features/ThreadItem';
import type { PostEntity } from '@/types/domain';
import type { ThreadViewItem } from '@/types/view';
import type { BackState } from '@/components/features/ThreadTree';
import { adaptToThreadView } from '@/lib/utils/thread-adapter';

export const dynamic = 'force-dynamic';

interface Props {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}

interface ThreadNode extends PostEntity {
  children: ThreadNode[];
}

function buildThreadTree(replies: PostEntity[]): ThreadNode[] {
  const nodeMap = new Map<number, ThreadNode>();
  const roots: ThreadNode[] = [];

  for (const reply of replies) {
    nodeMap.set(reply.id, { ...reply, children: [] });
  }

  for (const reply of replies) {
    const node = nodeMap.get(reply.id)!;
    const isRoot = !reply.parentId || reply.parentId <= 0;

    if (isRoot) {
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

function adaptTreeNode(node: ThreadNode): ThreadViewItem {
  const base = adaptToThreadView(node);
  return {
    ...base,
    replies: node.children.length > 0 ? node.children.map(adaptTreeNode) : undefined,
  };
}

/**
 * 从 searchParams 中提取回溯参数
 * 同时用于构建返回 URL 和传递给子组件的 backState
 */
function extractBackContext(searchParams: Record<string, string | string[] | undefined>): {
  backUrl: string;
  backState: BackState;
} {
  const nodeId = searchParams.nodeId;
  const page = searchParams.page;

  const backState: BackState = {};
  if (nodeId) backState.nodeId = String(nodeId);
  if (page) backState.page = String(page);

  // 无参数时返回干净首页
  if (!backState.nodeId && !backState.page) {
    return { backUrl: '/', backState: {} };
  }

  const params = new URLSearchParams();
  if (backState.nodeId) params.set('nodeId', backState.nodeId);
  if (backState.page) params.set('page', backState.page);

  return {
    backUrl: `/?${params.toString()}`,
    backState,
  };
}

export default async function ReadPage({ params, searchParams }: Props) {
  let id: string;
  try {
    const resolved = await params;
    id = resolved.id;
  } catch {
    notFound();
  }

  let backUrl = '/';
  let backState: BackState = {};
  try {
    const sp = await searchParams;
    const ctx = extractBackContext(sp);
    backUrl = ctx.backUrl;
    backState = ctx.backState;
  } catch {
    // 解析失败静默降级
  }

  let post: PostEntity | null = null;
  let replies: PostEntity[] = [];

  try {
    const result = await getPostWithThread(id);
    post = result.post;
    replies = result.replies ?? [];
  } catch (error) {
    console.error(`[ReadPage] Failed to fetch post ${id}:`, error);
  }

  if (!post) {
    notFound();
  }

  const treeNodes = buildThreadTree(replies);
  const adaptedReplies = treeNodes.map(adaptTreeNode);
  const totalReplyCount = replies.length - 1;

  return (
    <div className="main-body">
      <div className="detail-back-bar">
        <Link className="back-list-btn" href={backUrl}>
          ← 返回列表
        </Link>
      </div>

      <PostDetailCard post={post} replyCount={totalReplyCount} />

      <div className="thread-tree-container">
        <div className="thread-tree-header">
          <div className="thread-tree-title">💬 回帖讨论 ({totalReplyCount})</div>
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
              currentPostId={String(post!.id)}
              backState={backState}
            />
          ))}
        </ul>
      </div>
    </div>
  );
}