// src/app/threads/[slug]/page.tsx
import { notFound } from 'next/navigation';
import Link from 'next/link';
import type { Metadata } from 'next';
import { getPostTree, getPostFlat, getPostDetail } from '@/services/post-service';
import { buildThreadTree } from '@/lib/utils/thread-tree';
import { getBackContext } from '@/lib/utils/back-context';
import { ViewModeSwitcher } from '@/components/ui/ViewModeSwitcher';
import { PostTree } from './PostTree';
import { PostFlat } from './PostFlat';
import { ReadTracker } from '@/components/features/ReadTracker';

export const dynamic = 'force-dynamic';

interface Props {
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}

// ==================== SEO 元数据 ====================

export async function generateMetadata({ params, searchParams }: Props): Promise<Metadata> {
  try {
    const { slug } = await params;

    let forceRefresh = false;
    try {
      const resolvedSp = await searchParams;
      forceRefresh = resolvedSp?.refresh === '1';
    } catch {
      // ignore
    }

    const detail = await getPostDetail(slug, { noCache: forceRefresh });
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

// ==================== JSON-LD 结构化数据 ====================

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

// ==================== 页面组件 ====================

export default async function ReadPage({ params, searchParams }: Props) {
  let slug: string;
  try {
    const resolved = await params;
    slug = resolved.slug;
  } catch {
    notFound();
  }

  // 仅保留 view 和 refresh 两个功能性 URL 参数
  let resolvedSp: Record<string, string | string[] | undefined> = {};
  try {
    resolvedSp = await searchParams;
  } catch {
    // ignore
  }

  const view = resolvedSp.view === 'flat' ? 'flat' : 'tree';
  const forceRefresh = resolvedSp.refresh === '1';
  const serviceOpts = { noCache: forceRefresh };

  // ✅ 从 Referer 获取返回链接，URL 不再携带 backState 参数
  const { backUrl } = await getBackContext();

  let currentPost = null;
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
      totalReplyCount = Math.max(0, posts.length - 1);
      currentPost = posts[0];
    } else {
      const result = await getPostTree(slug, serviceOpts);
      currentPost = result.currentPost;
      const posts = result.posts ?? [];
      viewData = buildThreadTree(posts);
      totalReplyCount = posts.length > 0 ? posts.length - 1 : 0;
    }
  } catch (error) {
    console.error(`[ReadPage] Failed to fetch currentPost ${slug} (${view}):`, error);
  }

  if (!currentPost) {
    notFound();
  }

  return (
    <>
      <JsonLd post={currentPost} totalReplyCount={totalReplyCount} />

      <div className="detail-top-bar">
        <Link className="back-list-btn" href={backUrl}>← 返回列表</Link>
        <ViewModeSwitcher currentView={view} />
      </div>

      {view === 'flat' ? (
        <PostFlat posts={viewData} totalReplyCount={totalReplyCount} />
      ) : (
        <PostTree
          post={currentPost}
          viewPosts={viewData}
          totalReplyCount={totalReplyCount}
        />
      )}
      <ReadTracker postSlug={slug} nodeSlug={currentPost.node?.slug} />
    </>
  );
}