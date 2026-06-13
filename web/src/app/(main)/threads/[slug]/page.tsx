// src/app/threads/[slug]/page.tsx
import { notFound } from 'next/navigation';
import Link from 'next/link';
import type { Metadata } from 'next';
import { getPostWithThread, getPostFlat, getPostDetail } from '@/services/post-service';
import { buildThreadTree } from '@/lib/utils/thread-tree';
import { ViewModeSwitcher } from '@/components/ui/ViewModeSwitcher';
import { PostTree } from './PostTree';
import { PostFlat } from './PostFlat';
import type { BackState } from '@/types/view';
import { ReadTracker } from '@/components/features/ReadTracker';

export const dynamic = 'force-dynamic';

interface Props {
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}

// 动态生成 SEO 元数据
export async function generateMetadata({ params, searchParams }: Props): Promise<Metadata> {
  try {
    const { slug } = await params;
    // 解析 searchParams 并提取刷新信号
    let forceRefresh = false;
    try {
      const resolvedSp = await searchParams;
      forceRefresh = resolvedSp?.refresh === '1';
    } catch {
      // ignore
    }
    
    const serviceOpts = { noCache: forceRefresh };

    // 将 serviceOpts 透传给 getPostDetail
    const detail = await getPostDetail(slug, serviceOpts);

    if (!detail) return {};

    const title = `${detail.title || '无标题'} - 讨论详情`;
    const description =
      detail.summary ||
      detail.content?.replace(/[#*`>\-\[\]]/g, '').slice(0, 160) ||
      '参与社区讨论，查看精彩回复。';

    return {
      title,
      description,
      openGraph: {
        title,
        description,
        type: 'article',
        url: `/threads/${slug}`,
        // 如果有封面图可在此添加: images: [detail.coverUrl],
      },
      twitter: {
        card: 'summary_large_image',
        title,
        description,
      },
      alternates: {
        canonical: `/threads/${slug}`,
      },
    };
  } catch {
    return {};
  }
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

  if (!backState.nodeSlug && !backState.tagSlug && !backState.page) {
    return { backUrl: '/', backState: {} };
  }

  let basePath = '/';
  if (backState.tagSlug) {
    basePath = `/tags/${backState.tagSlug}`;
  } else if (backState.nodeSlug) {
    basePath = `/nodes/${backState.nodeSlug}`;
  }

  const urlParams = new URLSearchParams();
  if (backState.page) urlParams.set('page', backState.page);

  const queryString = urlParams.toString();
  const backUrl = queryString ? `${basePath}?${queryString}` : basePath;

  return { backUrl, backState };
}

// JSON-LD 结构化数据组件（帮助搜索引擎理解帖子内容）
function JsonLd({ post, totalReplyCount }: { post: any; totalReplyCount: number }) {
  const jsonLd = {
    '@context': 'https://schema.org',
    '@type': 'DiscussionForumPosting',
    headline: post.title || '无标题',
    articleBody: post.content?.slice(0, 5000),
    datePublished: post.createdAt,
    dateModified: post.updatedAt,
    author: post.user
      ? {
          '@type': 'Person',
          name: post.user.nickname || post.user.username,
        }
      : undefined,
    interactionStatistic: {
      '@type': 'InteractionCounter',
      interactionType: 'https://schema.org/CommentAction',
      userInteractionCount: totalReplyCount,
    },
  };

  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
    />
  );
}

export default async function ReadPage({ params, searchParams }: Props) {
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

  const view = resolvedSp.view === 'flat' ? 'flat' : 'tree';
  const { backUrl, backState } = extractBackContext(resolvedSp);

  const forceRefresh = resolvedSp.refresh === '1';
  const serviceOpts = { noCache: forceRefresh };

  let post = null;
  let viewData: any = null;
  let totalReplyCount = 0;

  try {
    if (view === 'flat') {
      let result = await getPostFlat(slug);
      let posts = result.posts ?? [];

      const isNonRootFlat = posts.length === 0 || (posts[0] && !posts[0].isRoot);

      if (isNonRootFlat) {
        const detail = await getPostDetail(slug, serviceOpts);
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
      const result = await getPostWithThread(slug, serviceOpts);
      post = result.post;
      const replies = result.replies ?? [];
      viewData = buildThreadTree(replies);
      totalReplyCount = replies.length > 0 ? replies.length - 1 : 0;
    }
  } catch (error) {
    console.error(`[ReadPage] Failed to fetch post ${slug} (${view}):`, error);

    if (view === 'flat') {
      try {
        const detail = await getPostDetail(slug);
        if (detail.threadSlug) {
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

  return (
    <>
      {/* ✅ SEO 结构化数据 */}
      <JsonLd post={post} totalReplyCount={totalReplyCount} />

      <div className="detail-top-bar">
        <Link className="back-list-btn" href={backUrl}>← 返回列表</Link>
        <ViewModeSwitcher currentView={view} />
      </div>

      {view === 'flat' ? (
        <PostFlat posts={viewData} totalReplyCount={totalReplyCount} />
      ) : (
        <PostTree
          post={post}
          viewPosts={viewData}
          totalReplyCount={totalReplyCount}
          backState={backState}
        />
      )}
      <ReadTracker postSlug={slug} nodeSlug={post.node?.slug} />
    </>
  );
}