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
import { assembleSideload } from '@/lib/utils/assemble-sideload';
import type { IncludedData, AssembledPost } from '@/lib/utils/assemble-sideload';

export const dynamic = 'force-dynamic';

interface Props {
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}

// ==================== SEO 元数据 ====================
export async function generateMetadata({ params, searchParams }: Props): Promise<Metadata> {
  try {
    const { slug } = await params;
    const resolvedSp = await searchParams;
    const forceRefresh = resolvedSp?.refresh === '1';

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
      openGraph: { title, description, type: 'article', url: `/threads/${slug}` },
      twitter: { card: 'summary_large_image', title, description },
      alternates: { canonical: `/threads/${slug}` },
    };
  } catch (error) {
    // ✅ 不再静默吞错，便于排查 SEO 问题
    console.warn(`[generateMetadata] Failed for ${params}:`, error);
    return {};
  }
}

// ==================== JSON-LD 结构化数据 ====================
function JsonLd({ post, totalReplyCount }: { post: AssembledPost; totalReplyCount: number }) {
  const jsonLd = {
    '@context': 'https://schema.org',
    '@type': 'DiscussionForumPosting',
    headline: post.title || '无标题',
    articleBody: post.content?.slice(0, 5000),
    // ✅ 兼容 createTime/createdAt 两种命名
    datePublished: post.createdAt ?? post.createTime,
    dateModified: post.updatedAt ?? post.lastCommentTime,
    author: post.user
      ? { '@type': 'Person' as const, name: post.user.nickname || post.user.username }
      : undefined,
    interactionStatistic: {
      '@type': 'InteractionCounter' as const,
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

  let resolvedSp: Record<string, string | string[] | undefined> = {};
  try {
    resolvedSp = await searchParams;
  } catch { /* ignore */ }

  const view = resolvedSp.view === 'flat' ? 'flat' : 'tree';
  const forceRefresh = resolvedSp.refresh === '1';
  const serviceOpts = { noCache: forceRefresh };
  const { backUrl } = await getBackContext();

  let currentPost: AssembledPost | null = null;
  let viewData: any = null;
  let totalReplyCount = 0;
  let includedNodes: IncludedData['nodes'] = [];

  try {
    if (view === 'flat') {
      let rsp = await getPostFlat(slug);
      let posts = assembleSideload(rsp.data ?? [], rsp.included);

      const isNonRootFlat = posts.length === 0 || !posts[0]?.isRoot;
      if (isNonRootFlat) {
        const detail = await getPostDetail(slug, serviceOpts);
        if (!detail.threadSlug) {
          throw new Error(`Post ${slug} is missing threadSlug, cannot resolve flat view`);
        }
        rsp = await getPostFlat(String(detail.threadSlug));
        posts = assembleSideload(rsp.data ?? [], rsp.included);
      }

      viewData = posts;
      totalReplyCount = Math.max(0, posts.length - 1);
      currentPost = posts[0] ?? null;
      includedNodes = rsp.included?.nodes ?? [];
    } else {
      const rsp = await getPostTree(slug, serviceOpts);
      currentPost = rsp.extra as AssembledPost | null;

      const rawData = rsp.data;
      const postsArray = Array.isArray(rawData)
        ? rawData
        : rawData != null
          ? [rawData]       // 单对象 → 包装为数组
          : [];             // null/undefined → 空数组

      const assembledPosts = assembleSideload(postsArray, rsp.included);
      viewData = buildThreadTree(assembledPosts);
      totalReplyCount = assembledPosts.length > 0 ? assembledPosts.length - 1 : 0;
      includedNodes = rsp.included?.nodes ?? [];
    }
  } catch (error) {
    console.error(`[ReadPage] Failed to fetch ${slug} (${view}):`, error);
  }

  if (!currentPost) notFound();

  // 统一补全 node
  if (!currentPost.node) {
    currentPost.node = includedNodes[0] ?? null;
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
        <PostTree post={currentPost} viewPosts={viewData} totalReplyCount={totalReplyCount} />
      )}

      <ReadTracker postSlug={slug} nodeSlug={currentPost.node?.slug} />
    </>
  );
}