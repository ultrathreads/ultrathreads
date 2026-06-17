"use client";

import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

// 纯函数保持不变
function formatTimestamp(timestamp: number | string | undefined | null, t: any): string {
  if (!timestamp) return '';
  const ts = typeof timestamp === 'string' ? new Date(timestamp).getTime() : timestamp;
  const diffMs = Date.now() - ts;
  
  if (diffMs < 60000) return t('justNow');
  const min = Math.floor(diffMs / 60000);
  if (min < 60) return t('minutesAgo', { n: min });
  const hour = Math.floor(diffMs / 3600000);
  if (hour < 24) return t('hoursAgo', { n: hour });
  const day = Math.floor(diffMs / 86400000);
  if (day < 7) return t('daysAgo', { n: day });
  
  const d = new Date(ts);
  return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`;
}

interface RelativeTimeProps {
  timestamp: number | string | undefined | null;
  className?: string;
}

export function RelativeTime({ timestamp, className }: RelativeTimeProps) {
  const { t } = useTranslation();
  
  const [text, setText] = useState('');
  
  useEffect(() => {
    setText(formatTimestamp(timestamp, t));
    const id = setInterval(() => setText(formatTimestamp(timestamp, t)), 60_000);
    return () => clearInterval(id);
  }, [timestamp, t]);

  const iso = timestamp ? new Date(timestamp).toISOString() : undefined;

  // 如果 text 为空，可以渲染一个不可见的占位符防止布局偏移(CLS)
  // 或者直接渲染机器可读的 dateTime 属性，对 SEO 依然友好
  return (
    <time className={className} dateTime={iso}>
      {text || '\u00A0'} 
    </time>
  );
}