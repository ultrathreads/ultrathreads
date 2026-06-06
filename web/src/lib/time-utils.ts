// lib/time-utils.ts

/**
 * 安全地将时间戳转换为本地化字符串
 * 自动兼容 10位(秒) 和 13位(毫秒) 时间戳
 */
export function formatTimestamp(timestamp: number | string | undefined | null): string {
  if (!timestamp) return '未知时间';
  
  const ts = typeof timestamp === 'string' ? Number(timestamp) : timestamp;
  
  // 小于 9999999999 (约2286年) 视为秒级时间戳，自动补全为毫秒
  const milliseconds = ts < 10000000000 ? ts * 1000 : ts;
  
  return new Date(milliseconds).toLocaleString('zh-CN');
}