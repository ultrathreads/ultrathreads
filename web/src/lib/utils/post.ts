// src/lib/utils/post.ts

interface ExtractTitleOptions {
  /** 最大截取长度，默认 20 */
  maxLength?: number;
}

/**
 * 从 Markdown 内容中提取纯文本标题
 * - 去除所有 Markdown 语法标记
 * - 取第一行有效文本
 * - 按指定长度截取
 * - 兜底返回空字符串
 */
export function extractPostTitle(
  content: string,
  options: ExtractTitleOptions = {}
): string {
  const { maxLength = 20 } = options;

  // 非字符串/空内容直接返回空
  if (!content || typeof content !== 'string') return '';

  let text = content;

  // 1. 优先处理【行首语法】(必须最先执行，多行匹配)
  // 引用 >
  text = text.replace(/^>\s*/gm, '');
  // 无序列表 - + *
  text = text.replace(/^[-+*]\s+/gm, '');
  // 有序列表 1. 2.
  text = text.replace(/^\d+\.\s+/gm, '');
  // 分割线、表格行首简单清理
  text = text.replace(/^[|:-]+\s*/gm, '');

  // 2. 清理标题 # 标记
  text = text.replace(/#{1,6}\s*/g, '');

  // 3. 清理行内格式 * _ ~ `
  text = text.replace(/[*_~`]/g, '');

  // 4. 图片、链接处理
  text = text.replace(/!\[[^\]]*\]\([^)]*\)/g, '');
  text = text.replace(/\[([^\]]*)\]\([^)]*\)/g, '$1');

  // 5. 清理不可见字符、全角空格、普通空格
  text = text.replace(/[\u200B-\u200D\uFEFF\u00A0\s]+/g, ' ').trim();

  // 6. 取第一行非空有效行
  const lines = text.split('\n');
  let firstValidLine = '';
  for (const line of lines) {
    const trimLine = line.trim();
    if (trimLine) {
      firstValidLine = trimLine;
      break;
    }
  }

  // 7. 截取长度
  const title = firstValidLine.slice(0, maxLength).trim();

  // 纯符号/空白则返回空，否则返回标题
  // 正则：纯标点/特殊字符判定
  if (!title || /^[^\u4e00-\u9fa5a-zA-Z0-9]+$/.test(title)) {
    return '';
  }

  return title;
}