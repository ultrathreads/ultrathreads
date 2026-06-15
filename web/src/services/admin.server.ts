// src/services/admin.server.ts

import { apiFetch } from '@/lib/api/client';
import type { ApiResponse } from '@/types/api';
import type { SystemInfo } from '@/types/admin';

/**
 * 获取后台仪表盘系统信息（仅服务端调用）
 */
export async function getSystemInfo(): Promise<SystemInfo> {
  const res = await apiFetch<ApiResponse<SystemInfo>>('/admin/dashboard/systeminfo', {
    method: 'GET',
    auth: true,
    skipDataUnwrap: true,
    cacheStrategy: { next: { revalidate: 30 } },
  });
  return res.data;
}