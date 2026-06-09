// src/lib/auth/session.ts
import { cookies } from 'next/headers';
import { getCurrentUser, type CurrentUser } from '@/services/auth.server';

/**
 * 服务端专用：安全地获取当前会话用户
 * - 直接读取 HttpOnly Cookie，不经过 API 路由
 * - 未登录或 Token 失效时返回 null，永不抛错
 */
export async function getServerSession(): Promise<CurrentUser | null> {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get('access_token')?.value;

    if (!token) return null;

    // 复用已有的用户获取逻辑（可传入 token 参数）
    const user = await getCurrentUser(token);
    return user ?? null;
  } catch {
    // Token 过期、格式错误等所有异常统一视为未登录
    return null;
  }
}