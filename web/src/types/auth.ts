// src/types/auth.ts

export interface RegisterFormData {
  username: string;
  email: string;
  password: string;
}

/** 登录请求参数 */
export interface LoginParams {
  username: string;
  password: string;
}

/** /auth/login 响应数据 */
export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  expire_at: string;
}

/** /auth/login/refresh 响应数据 */
export interface RefreshResponse {
  access_token: string;
  refresh_token?: string;
  expires_in: number;
}

/** /user/current 响应数据 */
export interface CurrentUser {
  id: number;
  username: string;
  email: string;
  nickname: string;
  avatar: string;
  website: string;
  description: string;
  score: number;
  topicCount: number;
  commentCount: number;
  passwordSet: boolean;
  status: number;
  createdAt: string;

  roles?: string[];
  permissions?: string[];
}