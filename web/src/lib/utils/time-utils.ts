// lib/time-utils.ts
import type { TFunction } from 'i18next';

/**
 * 毫秒时间戳 转为 相对时间文本
 * @param timestamp 毫秒时间戳
 * @param t i18next 翻译函数（由调用方注入）
 */
export function formatTimestamp(
  timestamp: number | string | undefined | null,
  t: TFunction
): string {
  if (!timestamp) return '';

  const ts = typeof timestamp === 'string' ? new Date(timestamp).getTime() : timestamp;
  const now = Date.now();
  const diffMs = now - ts;

  const MS_PER_SEC = 1000;
  const MS_PER_MIN = 60 * MS_PER_SEC;
  const MS_PER_HOUR = 60 * MS_PER_MIN;
  const MS_PER_DAY = 24 * MS_PER_HOUR;

  // 小于 1 分钟
  if (diffMs < MS_PER_MIN) {
    return t('common:justNow');
  }

  // 分钟
  const min = Math.floor(diffMs / MS_PER_MIN);
  if (min < 60) {
    return t('common:minutesAgo', { n: min });
  }

  // 小时
  const hour = Math.floor(diffMs / MS_PER_HOUR);
  if (hour < 24) {
    return t('common:hoursAgo', { n: hour });
  }

  // 天数 < 7 天
  const day = Math.floor(diffMs / MS_PER_DAY);
  if (day < 7) {
    return t('common:daysAgo', { n: day });
  }

  // 超过7天，返回 YYYY-MM-DD（日期格式也可按需放入语言包）
  const date = new Date(ts);
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, '0');
  const d = String(date.getDate()).padStart(2, '0');
  return `${y}-${m}-${d}`;
}