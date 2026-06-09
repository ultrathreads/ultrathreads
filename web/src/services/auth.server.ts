// src/services/auth.server.ts

import { apiFetch } from '@/lib/api/client';
import type { ApiResponse } from '@/types/api';
import type {
  LoginParams,
  LoginResponse,
  RefreshResponse,
  CurrentUser,
} from '@/types/auth';

/**
 * 获取当前登录用户信息（服务端/客户端均可调用）
 * ✅ 必须带 auth: true 自动注入 Cookie 中的 Token
 */
export async function getCurrentUser(): Promise<CurrentUser> {
  return apiFetch<CurrentUser>('/user/current', {
    method: 'GET',
    auth: true,
    cacheStrategy: { next: { revalidate: 60 } },
  });
}

/**
 * 登录（服务端调用）
 * ✅ skipDataUnwrap 保留完整信封供 Route Handler 写入 Cookie
 */
export async function login(params: LoginParams): Promise<ApiResponse<LoginResponse>> {
  return apiFetch<ApiResponse<LoginResponse>>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(params),
    cacheStrategy: { next: { revalidate: 0 } },
    skipDataUnwrap: true,
  });
}

/**
 * 刷新 Token（服务端调用）
 * ✅ skipDataUnwrap 保留完整信封供 Route Handler 判断是否更新 Cookie
 */
export async function refreshToken(refreshToken: string): Promise<ApiResponse<RefreshResponse>> {
  return apiFetch<ApiResponse<RefreshResponse>>('/auth/login/refresh', {
    method: 'POST',
    body: JSON.stringify({ refresh_token: refreshToken }),
    cacheStrategy: { next: { revalidate: 0 } },
    skipDataUnwrap: true,
  });
}