// src/app/post/[id]/page.tsx
import { notFound } from 'next/navigation';
import Link from 'next/link';
import { getPostWithThread, getPostFlat, getPostDetail } from '@/services/post-service';
import { buildThreadTree } from '@/lib/utils/thread-tree';
import { ViewModeSwitcher } from '@/components/ViewModeSwitcher';
import { PostTree } from './PostTree';
import { PostFlat } from './PostFlat';
import type { BackState } from '@/components/features/ThreadTree';
import { ReadTracker } from './ReadTracker';

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
  const tagId = searchParams.tagId;
  const page = searchParams.page;

  const backState: BackState = {};
  if (nodeId) backState.nodeId = String(nodeId);
  if (tagId) backState.tagId = String(tagId);
  if (page) backState.page = String(page);

  if (!backState.nodeId && !backState.tagId && !backState.page) {
    return { backUrl: '/', backState: {} };
  }

  const params = new URLSearchParams();
  if (backState.nodeId) params.set('nodeId', backState.nodeId);
  if (backState.tagId) params.set('tagId', backState.tagId);
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
      // ✅ 先尝试用当前 ID 获取平铺数据
      let result = await getPostFlat(id);
      let posts = result.posts ?? [];

      // ✅ 判断是否为非根帖（根据实际 API 行为调整此条件）
      const isNonRootFlat =
        posts.length === 0 ||
        (posts[0] && String(posts[0].id) !== String(id));

      if (isNonRootFlat) {
        // ✅ 轻量级获取元数据，仅取 threadId
        const detail = await getPostDetail(id);

        // ✅ threadId 缺失时直接抛出，让外层 catch 兜底或触发 notFound
        if (!detail.threadId) {
          throw new Error(`Post ${id} is missing threadId, cannot resolve flat view`);
        }

        const threadId = String(detail.threadId);
        result = await getPostFlat(threadId);
        posts = result.posts ?? [];
      }

      viewData = posts;
      totalReplyCount = posts.length > 0 ? posts.length - 1 : 0;
      post = posts[0];
    } else {
      // 树形模式保持原有逻辑不变
      const result = await getPostWithThread(id);
      post = result.post;
      const replies = result.replies ?? [];
      viewData = buildThreadTree(replies);
      totalReplyCount = replies.length > 0 ? replies.length - 1 : 0;
    }
  } catch (error) {
    console.error(`[ReadPage] Failed to fetch post ${id} (${view}):`, error);

    // ✅ 兜底：flat 模式报错时，用轻量接口重试一次
    if (view === 'flat') {
      try {
        const detail = await getPostDetail(id);

        if (!detail.threadId) {
          console.error(`[ReadPage] Post ${id} missing threadId, skip flat retry`);
        } else {
          const threadId = String(detail.threadId);
          const result = await getPostFlat(threadId);
          const posts = result.posts ?? [];

          if (posts.length > 0) {
            viewData = posts;
            totalReplyCount = posts.length - 1;
            post = posts[0];
          }
        }
      } catch (retryError) {
        console.error(`[ReadPage] Flat mode retry failed for ${id}:`, retryError);
      }
    }
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
        />
      ) : (
        <PostTree
          post={post}
          viewPosts={viewData}
          totalReplyCount={totalReplyCount}
          backState={backState}
        />
      )}
      <ReadTracker nodeId={String(post.node?.nodeId ?? '')} postId={id} />
    </>
  );
}