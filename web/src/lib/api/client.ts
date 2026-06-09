// src/lib/api/client.ts
import type { ApiResponse } from '@/types/api';

const GO_API_BASE = process.env.GO_API_BASE || 'http://localhost:9527/api';

export interface ApiRequestOptions extends RequestInit {
  auth?: boolean;
  cacheStrategy?: NextFetchRequestConfig;
  skipDataUnwrap?: boolean;
}

/** 自定义业务异常，方便上层区分 HTTP 错误与业务错误 */
export class ApiBusinessError extends Error {
  constructor(
    message: string,
    public code: number,
    public raw?: unknown
  ) {
    super(message);
    this.name = 'ApiBusinessError';
  }
}

// ✅ 函数重载：skipDataUnwrap=true 时自动推导返回完整信封
export async function apiFetch<T>(
  path: string,
  options: ApiRequestOptions & { skipDataUnwrap: true }
): Promise<ApiResponse<T>>;
export async function apiFetch<T>(
  path: string,
  options?: ApiRequestOptions
): Promise<T>;
export async function apiFetch<T>(
  path: string,
  options: ApiRequestOptions = {}
): Promise<T | ApiResponse<T>> {
  const {
    auth = false,
    cacheStrategy = { next: { revalidate: 60 } },
    skipDataUnwrap = false,
    ...fetchOptions
  } = options;

  const headers = new Headers(fetchOptions.headers);
  if (!(fetchOptions.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }

  let credentials: RequestCredentials | undefined = fetchOptions.credentials;

  if (auth) {
    if (typeof window === 'undefined') {
      try {
        const { cookies } = await import('next/headers');
        const cookieStore = await cookies();
        const token = cookieStore.get('access_token')?.value;
        if (!token) throw new ApiBusinessError('AUTH_REQUIRED', -1);
        headers.set('Authorization', `Bearer ${token}`);
      } catch (e) {
        // 区分"无 Token"和"读取 Cookie 基础设施异常"
        if (e instanceof ApiBusinessError) throw e;
        console.error('[apiFetch] Failed to read server cookies:', e);
        throw new ApiBusinessError('AUTH_REQUIRED', -1);
      }
    } else {
      credentials = 'include';
    }
  }

  const res = await fetch(`${GO_API_BASE}${path}`, {
    ...fetchOptions,
    headers,
    credentials,
    ...cacheStrategy,
  });

  // ✅ 优先尝试解析响应体，避免丢失后端返回的业务错误信息
  let envelope: ApiResponse<unknown> | null = null;
  const contentType = res.headers.get('content-type') ?? '';
  if (contentType.includes('application/json')) {
    try {
      envelope = (await res.json()) as ApiResponse<unknown>;
    } catch {
      // JSON 解析失败，envelope 保持 null
    }
  }

  // ✅ HTTP 层错误：优先使用信封中的业务消息
  if (!res.ok) {
    const message = envelope?.message || `HTTP ${res.status} ${res.statusText}`;
    const code = envelope?.code ?? res.status;
    throw new ApiBusinessError(message, code, envelope);
  }

  // ✅ 无响应体或非法信封
  if (!envelope) {
    throw new ApiBusinessError('Invalid API response: missing JSON body', -2);
  }

  // ✅ skipDataUnwrap 模式：直接返回完整信封（由重载保证类型安全）
  if (skipDataUnwrap) {
    return envelope as ApiResponse<T>;
  }

  // ✅ 以 success 为唯一权威成功标识
  if (!envelope.success) {
    throw new ApiBusinessError(
      envelope.message || 'Unknown business error',
      envelope.code,
      envelope
    );
  }

  return envelope.data as T;
}