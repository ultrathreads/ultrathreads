// src/services/auth.server.ts

import { apiFetch } from '@/lib/api/client';
import { ApiError } from '@/lib/api/client'; 
import type { ApiResponse } from '@/types/api';
import type {
  LoginParams,
  LoginResponse,
  RefreshResponse,
  CurrentUser,
} from '@/types/auth';

/**
 * 获取当前登录用户信息（服务端/客户端均可调用）
 * ✅ 未登录或 Token 无效时返回 null，而非抛出 401 异常
 */
export async function getCurrentUser(): Promise<CurrentUser | null> {
  try {
    return await apiFetch<CurrentUser>('/user/current', {
      method: 'GET',
      auth: true,
      cacheStrategy: { next: { revalidate: 60 } },
    });
  } catch (err) {
    // 401 = 未登录 / Token 过期 / Cookie 为空 → 安全返回 null
    if (err instanceof ApiError && err.status === 401) {
      return null;
    }
    // 其他错误（500、网络故障等）继续向上抛出，不吞掉真正的异常
    throw err;
  }
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