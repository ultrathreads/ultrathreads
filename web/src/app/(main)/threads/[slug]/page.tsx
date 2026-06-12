// src/app/threads/[slug]/page.tsx
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
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}

function extractBackContext(searchParams: Record<string, string | string[] | undefined>): {
  backUrl: string;
  backState: BackState;
} {
  const nodeSlug = searchParams.nodeSlug;
  const tagSlug = searchParams.tagSlug;
  const page = searchParams.page;

  const backState: BackState = {};
  if (nodeSlug) backState.nodeSlug = String(nodeSlug);
  if (tagSlug) backState.tagSlug = String(tagSlug);
  if (page) backState.page = String(page);

  // 如果没有任何有效状态，直接返回首页
  if (!backState.nodeSlug && !backState.tagSlug && !backState.page) {
    return { backUrl: '/', backState: {} };
  }

  // 根据 tagSlug 或 nodeSlug 动态决定基础路径
  let basePath = '/';
  if (backState.tagSlug) {
    basePath = `/tags/${backState.tagSlug}`;
  } else if (backState.nodeSlug) {
     basePath = `/nodes/${backState.nodeSlug}`;
  }

  // 构建查询参数（仅保留 page 等需要拼接到 URL 上的参数）
  const params = new URLSearchParams();
  if (backState.page) {
    params.set('page', backState.page);
  }

  const queryString = params.toString();
  const backUrl = queryString ? `${basePath}?${queryString}` : basePath;

  return {
    backUrl,
    backState,
  };
}

export default async function ReadPage({ params, searchParams }: Props) {
  // 1. 解析路由参数与搜索参数
  let slug: string;
  try {
    const resolved = await params;
    slug = resolved.slug;
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
      // 先尝试用当前 ID 获取平铺数据
      let result = await getPostFlat(slug);
      let posts = result.posts ?? [];

      // 判断是否为非根帖（根据实际 API 行为调整此条件）
      const isNonRootFlat =
        posts.length === 0 ||
        (posts[0] && !posts[0].isRoot);

      if (isNonRootFlat) {
        const detail = await getPostDetail(slug);

        // threadSlug 缺失时直接抛出，让外层 catch 兜底或触发 notFound
        if (!detail.threadSlug) {
          throw new Error(`Post ${slug} is missing threadSlug, cannot resolve flat view`);
        }

        const threadSlug = String(detail.threadSlug);
        result = await getPostFlat(threadSlug);
        posts = result.posts ?? [];
      }

      viewData = posts;
      totalReplyCount = posts.length > 0 ? posts.length - 1 : 0;
      post = posts[0];
    } else {
      // 树形模式保持原有逻辑不变
      const result = await getPostWithThread(slug);
      post = result.post;
      const replies = result.replies ?? [];
      viewData = buildThreadTree(replies);
      totalReplyCount = replies.length > 0 ? replies.length - 1 : 0;
    }
  } catch (error) {
    console.error(`[ReadPage] Failed to fetch post ${slug} (${view}):`, error);

    // ✅ 兜底：flat 模式报错时，用轻量接口重试一次
    if (view === 'flat') {
      try {
        const detail = await getPostDetail(slug);

        if (!detail.threadSlug) {
          console.error(`[ReadPage] Post ${slug} missing threadSlug, skip flat retry`);
        } else {
          const threadSlug = String(detail.threadSlug);
          const result = await getPostFlat(threadSlug);
          const posts = result.posts ?? [];

          if (posts.length > 0) {
            viewData = posts;
            totalReplyCount = posts.length - 1;
            post = posts[0];
          }
        }
      } catch (retryError) {
        console.error(`[ReadPage] Flat mode retry failed for ${slug}:`, retryError);
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
      <ReadTracker postSlug={slug} />
    </>
  );
}