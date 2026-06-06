// src/lib/api/auth.ts
import { apiFetch } from '@/lib/api/client';

/** Go 后端 /auth/login 响应数据内容 (data 字段内部结构) */
export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  expire_at: string;
}

/** Go 后端 /auth/login 完整响应信封 */
export interface LoginEnvelope {
  code: number;
  message: string;
  success: boolean;
  data: LoginResponse;
}

/** Go 后端 /auth/login/refresh 响应数据内容 (data 字段内部结构) */
export interface RefreshResponse {
  access_token: string;
  refresh_token?: string;
  expires_in: number;
}

/** Go 后端 /auth/login/refresh 完整响应信封 */
export interface RefreshEnvelope {
  code: number;
  message: string;
  success: boolean;
  data: RefreshResponse;
}

/** 登录请求参数 */
export interface LoginParams {
  username: string;
  password: string;
}

/** Go 后端 /user/current 响应数据结构 */
export interface CurrentUser {
  id: number;
  username: string;
  email: string;
  nickname: string;
  avatar: string;
  level: number;
  levelName: string;
  website: string;
  description: string;
  score: number;
  topicCount: number;
  commentCount: number;
  passwordSet: boolean;
  status: number;
  createTime: number; // Unix 毫秒时间戳
}

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
 * ✅ 使用 skipDataUnwrap 阻止 apiFetch 自动剥壳，确保 route.ts 能拿到完整信封
 */
export async function login(params: LoginParams): Promise<LoginEnvelope> {
  return apiFetch<LoginEnvelope>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(params),
    cacheStrategy: { next: { revalidate: 0 } },
    skipDataUnwrap: true,
  });
}

/**
 * 刷新 Token（服务端调用）
 * ✅ 同步修复：刷新接口同样需要读取外层 success 判断是否更新 Cookie
 */
export async function refreshToken(refreshToken: string): Promise<RefreshEnvelope> {
  return apiFetch<RefreshEnvelope>('/auth/login/refresh', {
    method: 'POST',
    body: JSON.stringify({ refresh_token: refreshToken }),
    cacheStrategy: { next: { revalidate: 0 } },
    skipDataUnwrap: true,
  });
}