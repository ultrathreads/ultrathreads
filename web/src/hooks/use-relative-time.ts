'use client';

import { useTranslation } from '@/lib/i18n/i18n-client';
import { formatTimestamp } from '@/lib/utils/time-utils';

/**
 * 客户端专用的相对时间格式化 Hook
 * 内部自动注入翻译函数，对外暴露无参格式化方法
 */
export function useRelativeTime() {
  const { t } = useTranslation(['common']);

  // 返回一个闭包函数，调用方无需再传 t
  return (timestamp: number | string | undefined | null) => 
    formatTimestamp(timestamp, t);
}