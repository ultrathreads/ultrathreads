// src/app/threads/[id]/PostTree.tsx
import type { PostEntity } from '@/types/domain';
import type { BackState } from '@/components/features/ThreadTree';
import PostDetailClient from '@/components/PostDetailClient';

interface PostTreeProps {
  post: PostEntity;
  viewPosts: ReturnType<typeof import('@/lib/utils/thread-tree').buildThreadTree>;
  totalReplyCount: number;
  backState: BackState;
}

export function PostTree({ post, viewPosts, totalReplyCount, backState }: PostTreeProps) {
  return (
    <PostDetailClient
      post={post}
      viewPosts={viewPosts}
      totalReplyCount={totalReplyCount}
      backState={backState}
    />
  );
}