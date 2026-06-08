// src/app/post/[id]/page.tsx
import { notFound } from 'next/navigation';
import Link from 'next/link';
import { getPostWithThread, getPostFlat } from '@/services/post-service';
import { buildThreadTree } from '@/lib/utils/thread-tree';
import { ViewModeSwitcher } from '@/components/ViewModeSwitcher';
import { PostTree } from './PostTree';
import { PostFlat } from './PostFlat';
import type { BackState } from '@/components/features/ThreadTree';

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
  // 1. 解析路由参数与搜索参数
  let id: string;
  try {
    const resolved = await params;
    id = resolved.id;
  } catch {
    notFound();
  }

  let resolvedSp: Record<string, string | string[] | undefined> = {};
  try {
    resolvedSp = await searchParams;
  } catch {
    // ignore
  }

  // 2. 确定视图模式（默认树形）
  const view = resolvedSp.view === 'flat' ? 'flat' : 'tree';

  // 3. 解析返回上下文
  const { backUrl, backState } = extractBackContext(resolvedSp);

  // 4. 根据视图模式获取对应数据并转换
  let post = null;
  let viewData: any = null;
  let totalReplyCount = 0;

  try {
    if (view === 'flat') {
      const result = await getPostFlat(id);
      const posts = result.posts ?? [];
      viewData = posts;
      totalReplyCount = posts.length -1 ;
      post = posts[0]
    } else {
      const result = await getPostWithThread(id);
      post = result.post;
      const replies = result.replies ?? [];
      viewData = buildThreadTree(replies);
      totalReplyCount = replies.length > 0 ? replies.length - 1 : 0;
    }
  } catch (error) {
    console.error(`[ReadPage] Failed to fetch post ${id} (${view}):`, error);
  }

  if (!post) {
    notFound();
  }

  // 5. 组装页面布局，根据视图模式分发组件
  return (
    <>
      <div className="detail-top-bar">
        <Link className="back-list-btn" href={backUrl}>← 返回列表</Link>
        <ViewModeSwitcher currentView={view} />
      </div>

      {view === 'flat' ? (
        <PostFlat
          posts={viewData}
          totalReplyCount={totalReplyCount}
          backState={backState}
        />
      ) : (
        <PostTree
          post={post}
          viewPosts={viewData}
          totalReplyCount={totalReplyCount}
          backState={backState}
        />
      )}
    </>
  );
}