// src/services/auth.client.ts

'use client';

import { ApiError } from '@/lib/api/client';
import type { ApiResponse } from '@/types/api';
import type { LoginParams, LoginResponse } from '@/types/auth';

/**
 * 客户端登录
 * ✅ 走 Next.js API Route，由 Route Handler 代理并写入 Cookie
 */
export async function loginClient(params: LoginParams): Promise<ApiResponse<LoginResponse>> {
  const res = await fetch('/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(params),
    credentials: 'include',
  });

  if (!res.ok) {
    const data = await res.json().catch(() => null);
    throw new ApiError(
      data?.error || data?.message || `登录失败 (${res.status})`,
      res.status
    );
  }

  return res.json();
}