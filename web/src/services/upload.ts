// src/services/upload.ts（精简版）
import { apiFetch } from '@/lib/api/client';

export interface UploadAvatarResult {
  url: string;
}

export async function uploadAvatar(file: File): Promise<UploadAvatarResult> {
  const formData = new FormData();
  formData.append('image', file);

  // skipDataUnwrap 默认为 false → 自动解包返回 data 部分
  return apiFetch<UploadAvatarResult>('/upload', {
    method: 'POST',
    body: formData,
    auth: true,
    // FormData 请求不应被缓存
    cacheStrategy: { cache: 'no-store' },
  });
}