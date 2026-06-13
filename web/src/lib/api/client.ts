// src/lib/api/client.ts
import type { ApiResponse } from '@/types/api';

const GO_API_BASE = process.env.GO_API_BASE || 'http://localhost:9527/api';

export interface ApiRequestOptions extends RequestInit {
  auth?: boolean;
  cacheStrategy?: NextFetchRequestConfig;
  skipDataUnwrap?: boolean;
  noCache?: boolean; 
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

// 函数重载：skipDataUnwrap=true 时自动推导返回完整信封
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

        // 没有 Token 时不再 throw，只是不设置 Header
        // Go 后端的 OptionalAuth 会自动降级为游客模式
        if (token) {
          headers.set('Authorization', `Bearer ${token}`);
        }
      } catch (e) {
        // Cookie 基础设施异常也降级为匿名请求，而非阻断页面渲染
        console.error('[apiFetch] Failed to read server cookies:', e);
      }
    } else {
      credentials = 'include';
    }
  }

  let finalCacheStrategy: RequestInit & NextFetchRequestConfig = {};
  
  if (options.noCache) {
    if (typeof window === 'undefined') {
      // 服务端：使用 Next.js 专属的 no-store
      finalCacheStrategy = { cache: 'no-store' };
    } else {
      // 客户端：拼接时间戳击穿浏览器/CDN缓存
      const separator = path.includes('?') ? '&' : '?';
      path = `${path}${separator}_t=${Date.now()}`;
      // 同时设置请求头，尝试击穿中间件
      headers.set('Cache-Control', 'no-cache, no-store, must-revalidate');
      headers.set('Pragma', 'no-cache');
    }
  } else if (auth) {
    // 保留你原有的 auth 缓存策略
    finalCacheStrategy = { cache: 'no-store' as RequestCache };
  } else {
    // 保留你原有的默认缓存策略
    finalCacheStrategy = cacheStrategy;
  }

  const res = await fetch(`${GO_API_BASE}${path}`, {
    ...fetchOptions,
    headers,
    credentials,
    ...finalCacheStrategy,
  });

  // ✅ 先读取原始文本，无论是否为 JSON 都保留用于调试
  // 注意：res.text() 和 res.json() 都是一次性消费，必须用 text + 手动 parse 替代
  let envelope: ApiResponse<unknown> | null = null;
  const contentType = res.headers.get('content-type') ?? '';
  let rawBodyText = '';

  try {
    rawBodyText = await res.text();
  } catch {
    // 读取失败则保持空字符串
  }

  if (contentType.includes('application/json') && rawBodyText) {
    try {
      envelope = JSON.parse(rawBodyText) as ApiResponse<unknown>;
    } catch {
      // JSON 解析失败，envelope 保持 null
    }
  }

  // HTTP 层错误：优先使用信封中的业务消息
  if (!res.ok) {
    const message = envelope?.message || `HTTP ${res.status} ${res.statusText}`;
    const code = envelope?.code ?? res.status;
    throw new ApiBusinessError(message, code, envelope ?? rawBodyText);
  }

  // ✅ 无响应体或非法信封时，直接暴露原始响应内容
  if (!envelope) {
    throw new ApiBusinessError(
      'Invalid API response: missing or non-JSON body',
      -2,
      {
        status: res.status,
        contentType,
        bodyPreview: rawBodyText.slice(0, 500),
      }
    );
  }

  // skipDataUnwrap 模式：直接返回完整信封（由重载保证类型安全）
  if (skipDataUnwrap) {
    return envelope as ApiResponse<T>;
  }

  // 以 success 为唯一权威成功标识
  if (!envelope.success) {
    throw new ApiBusinessError(
      envelope.message || 'Unknown business error',
      envelope.code,
      envelope
    );
  }

  // ✅ 兼容无 data 字段的合法成功响应（如 POST/PUT/DELETE 写操作）
  if (envelope.data === undefined) {
    return null as T;
  }

  return envelope.data as T;
}