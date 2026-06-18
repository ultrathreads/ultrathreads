// src/components/features/MyFavoritesList.tsx
'use client';

import Link from 'next/link';
import { useTransition, useOptimistic, useState } from 'react';
import { toast } from 'sonner'; 
import type { FavoriteItem } from '@/services/favorite-service';
import { RelativeTime } from '@/components/ui/RelativeTime';

interface Props {
  initialFavorites: FavoriteItem[];
  // ✅ 接收 Server Action
  onDeleteFavoriteAction: (formData: FormData) => Promise<void>;
}

export default function MyFavoritesList({ initialFavorites, onDeleteFavoriteAction }: Props) {
  const [isPending, startTransition] = useTransition();
  const [deletingId, setDeletingId] = useState<number | null>(null);
  
  // ✅ 乐观更新：点击删除时，UI 立即从列表中移除该项，无需等待网络请求
  const [optimisticFavorites, removeOptimistic] = useOptimistic(
    initialFavorites,
    (state, favoriteIdToRemove: number) => 
      state.filter((item) => item.favoriteId !== favoriteIdToRemove)
  );

  const handleDelete = (item: FavoriteItem) => {
    startTransition(async () => {
      setDeletingId(item.favoriteId);
      removeOptimistic(item.favoriteId); // 乐观更新：先移除 UI
      
      const formData = new FormData();
      formData.append('entityType', item.entityType);
      formData.append('entityId', String(item.entityId));
      
      try {
        // 调用 Server Action，如果服务端抛出错误，这里会进入 catch
        await onDeleteFavoriteAction(formData);
        
        // 执行到这里，说明服务端成功处理了请求
        // toast.success('已取消收藏');
      } catch (error) {
        // 捕获服务端抛出的错误
        toast.error(error instanceof Error ? error.message : '操作失败，请重试');
        // 注意：由于使用了乐观更新，如果失败，你可能还需要把该项加回列表
      } finally {
        // 无论成功还是失败，都重置按钮的 loading 状态
        setDeletingId(null);
      }
    });
  };

  const pageTitle = '我的书签';

  return (
    <div className="thread-tree-container">
      <div className="thread-tree-header">
        <h3>{pageTitle}</h3>
      </div>
      <ul className="thread">
        {optimisticFavorites.map((item) => (
          <li key={item.slug}>
            <div className="entry">
              <svg className="icon-favorite-svg" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="#f39c12">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
              </svg>

              <Link className="subject" href={item.url}>
                {item.title || '(无标题)'}
              </Link>

              <span className="metadata">
                {item.user && (
                  <span className="author-name" title={item.user.nickname}>
                    {item.user.nickname}
                  </span>
                )}
                <span className="tail">
                  <RelativeTime timestamp={item.createdAt} />
                </span>
                
                {/* 删除按钮 */}
                <button
                  onClick={() => handleDelete(item)}
                  disabled={isPending}
                  className="favorite-delete-btn detail-action-btn"
                  title="取消收藏"
                >
                  ✕
                </button>
              </span>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}