// src/lib/api/client.ts
import { cookies } from 'next/headers';

const GO_API_BASE = process.env.GO_API_BASE || 'http://localhost:8080';

/** 📦 Go 后端统一响应信封 */
interface ApiResponse<T> {
  code: number;
  success?: boolean;
  data: T;
  message: string;
}

export interface ApiRequestOptions extends RequestInit {
  auth?: boolean;
  cacheStrategy?: NextFetchRequestConfig;
  /**
   * ✅ 新增：跳过自动剥壳
   * 设为 true 时返回完整 ApiResponse 信封（含 success/code/message/data）
   * 默认 false，自动提取 data 字段并在业务失败时抛出异常
   */
  skipDataUnwrap?: boolean;
}

/** 自定义业务异常，方便上层区分 HTTP 错误与业务错误 */
export class ApiBusinessError extends Error {
  constructor(message: string, public code: number) {
    super(message);
    this.name = 'ApiBusinessError';
  }
}

export async function apiFetch<T>(
  path: string,
  options: ApiRequestOptions = {}
): Promise<T> {
  // ✅ 关键修复：显式解构 skipDataUnwrap，避免被透传给原生 fetch 导致失效
  const {
    auth = false,
    cacheStrategy = { next: { revalidate: 60 } },
    skipDataUnwrap = false,
    ...fetchOptions
  } = options;

  const headers = new Headers(fetchOptions.headers);
  headers.set('Content-Type', 'application/json');

  if (auth) {
    const cookieStore = await cookies();
    const token = cookieStore.get('access_token')?.value;
    if (!token) throw new ApiBusinessError('AUTH_REQUIRED', -1);
    headers.set('Authorization', `Bearer ${token}`);
  }

  const res = await fetch(`${GO_API_BASE}${path}`, {
    ...fetchOptions,
    headers,
    ...cacheStrategy,
  });

  // ✅ 精确识别 401，便于触发 Token 刷新
  if (res.status === 401) {
    throw new ApiBusinessError('UNAUTHORIZED', 401);
  }

  if (!res.ok) {
    throw new Error(`HTTP Error: ${res.status} ${res.statusText}`);
  }

  // 📦 解析统一响应信封
  const envelope = (await res.json()) as ApiResponse<T>;

  // ✅ skipDataUnwrap = true 时直接返回完整信封，不做任何业务校验与剥壳
  if (skipDataUnwrap) {
    return envelope as unknown as T;
  }

  // ✅ 默认模式：仅以 code === 0 且 success !== false 作为成功标准
  if (envelope.code !== 0 || envelope.success === false) {
    throw new ApiBusinessError(
      envelope.message || 'Unknown business error',
      envelope.code
    );
  }

  // ✅ 只返回纯净的业务数据
  return envelope.data;
}