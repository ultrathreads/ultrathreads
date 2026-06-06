// src/lib/api/client.ts
import { cookies } from 'next/headers';

const GO_API_BASE = process.env.GO_API_BASE || 'http://localhost:8080';

/** 📦 Go 后端统一响应信封 */
interface ApiResponse<T> {
  code: number;
  data: T;
  message: string;
}

export interface ApiRequestOptions extends RequestInit {
  auth?: boolean;
  cacheStrategy?: NextFetchRequestConfig;
}

/** 自定义业务异常，方便上层区分 HTTP 错误与业务错误 */
export class ApiBusinessError extends Error {
  constructor(public code: number, message: string) {
    super(message);
    this.name = 'ApiBusinessError';
  }
}

export async function apiFetch<T>(
  path: string,
  options: ApiRequestOptions = {}
): Promise<T> {
  const { auth = false, cacheStrategy = { next: { revalidate: 60 } }, ...fetchOptions } = options;

  const headers = new Headers(fetchOptions.headers);
  headers.set('Content-Type', 'application/json');

  if (auth) {
    const cookieStore = await cookies();
    const token = cookieStore.get('access_token')?.value;
    if (!token) throw new ApiBusinessError(-1, 'AUTH_REQUIRED');
    headers.set('Authorization', `Bearer ${token}`);
  }

  const res = await fetch(`${GO_API_BASE}${path}`, {
    ...fetchOptions,
    headers,
    ...cacheStrategy,
  });

  // ❌ 1. HTTP 层错误（网络不通、502、404等）
  if (!res.ok) {
    throw new Error(`HTTP Error: ${res.status} ${res.statusText}`);
  }

  // 📦 2. 拆解统一响应信封
  const envelope = (await res.json()) as ApiResponse<T>;

  // ❌ 3. 业务层错误（code !== 0）
  if (envelope.code !== 0) {
    throw new ApiBusinessError(envelope.code, envelope.message || 'Unknown business error');
  }

  // ✅ 4. 只返回纯净的业务数据
  return envelope.data;
}