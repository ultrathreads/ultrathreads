import { notFound } from 'next/navigation';
import Link from 'next/link';
import { getPostWithThread } from '@/services/post-service';
import type { PostEntity } from '@/types/domain';
import type { BackState } from '@/components/features/ThreadTree';
import { buildThreadTree } from '@/lib/utils/thread-tree';
import PostDetailClient from '@/components/PostDetailClient';

export const dynamic = 'force-dynamic';

interface Props {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}

function extractBackContext(searchParams: Record<string, string | string[] | undefined>): {
  backUrl: string;
  backState: BackState;
} {
  const nodeId = searchParams.nodeId;
  const page = searchParams.page;

  const backState: BackState = {};
  if (nodeId) backState.nodeId = String(nodeId);
  if (page) backState.page = String(page);

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

  const viewPosts = buildThreadTree(replies);
  const totalReplyCount = replies.length - 1;

  return (
    <div className="main-body">
      <div className="detail-back-bar">
        <Link className="back-list-btn" href={backUrl}>
          ← 返回列表
        </Link>
      </div>

      <PostDetailClient
        post={post}
        viewPosts={viewPosts}
        totalReplyCount={totalReplyCount}
        backState={backState}
      />
    </div>
  );
}