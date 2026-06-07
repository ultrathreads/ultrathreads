// src/lib/utils/post.ts

interface ExtractTitleOptions {
  /** 最大截取长度，默认 20 */
  maxLength?: number;
  /** 兜底标题，默认 '回复' */
  fallback?: string;
}

/**
 * 从 Markdown 内容中提取纯文本标题
 * - 去除所有 Markdown 语法标记
 * - 取第一行有效文本
 * - 按指定长度截取
 * - 兜底返回默认标题
 */
export function extractPostTitle(
  content: string,
  options: ExtractTitleOptions = {}
): string {
  const { maxLength = 20, fallback = '回复' } = options;

  const plainText = content
    .replace(/#{1,6}\s*/g, '')
    .replace(/[*_~`]/g, '')
    .replace(/\[([^\]]*)\]\([^)]*\)/g, '$1')
    .replace(/!\[[^\]]*\]\([^)]*\)/g, '')
    .replace(/>\s*/g, '')
    .replace(/[-+*]\s+/g, '')
    .replace(/\d+\.\s+/g, '')
    .trim();

  const firstLine = plainText.split('\n').find(line => line.trim().length > 0) || '';
  const title = firstLine.slice(0, maxLength).trim();

  // 兜底：如果提取结果为空或纯符号，使用默认标题
  if (!title || /^[\s\W]*$/.test(title)) {
    return fallback;
  }
  return title;
}