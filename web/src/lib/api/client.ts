// src/lib/api/client.ts

const GO_API_BASE = process.env.GO_API_BASE || 'http://localhost:9527/api';

export interface ApiRequestOptions extends RequestInit {
  auth?: boolean;
  cacheStrategy?: NextFetchRequestConfig;
  skipDataUnwrap?: boolean;
  noCache?: boolean;
}

/**
 * 统一 API 异常类
 * status: HTTP 状态码 (400, 403, 500...)
 * code:   后端业务错误码 (如 10001, "USER_NOT_FOUND")
 * message: 人类可读的错误信息
 */
export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
    public code?: number | string,
    public body?: unknown
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * 网络连接异常类（区别于 HTTP 业务错误）
 * 用于 ECONNREFUSED / ENOTFOUND / ETIMEDOUT 等网络层故障
 */
export class ApiConnectionError extends Error {
  constructor(
    path: string,
    public readonly originalError?: unknown
  ) {
    super('接口访问异常，服务暂时不可用，请稍后重试');
    this.name = 'ApiConnectionError';
  }
}

// 成功响应：纯净业务数据
interface SuccessEnvelope<T> {
  data: T;
  meta?: unknown;
}

// 错误响应：包含 code + message
interface ErrorEnvelope {
  code?: number | string;
  message?: string;
  [key: string]: unknown; // 允许 details、field 等扩展字段
}

// 函数重载
export async function apiFetch<T>(
  path: string,
  options: ApiRequestOptions & { skipDataUnwrap: true }
): Promise<SuccessEnvelope<T>>;
export async function apiFetch<T>(
  path: string,
  options?: ApiRequestOptions
): Promise<T>;
export async function apiFetch<T>(
  path: string,
  options: ApiRequestOptions = {}
): Promise<T | SuccessEnvelope<T>> {
  const {
    auth = false,
    cacheStrategy = { next: { revalidate: 60 } },
    skipDataUnwrap = false,
    ...fetchOptions
  } = options;

  // --- Header & Auth ---
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
        if (token) headers.set('Authorization', `Bearer ${token}`);
      } catch (e) {
        console.error('[apiFetch] Failed to read server cookies:', e);
      }
    } else {
      credentials = 'include';
    }
  }

  // --- Cache 策略 ---
  let finalCacheStrategy: RequestInit & NextFetchRequestConfig = {};
  if (options.noCache) {
    if (typeof window === 'undefined') {
      finalCacheStrategy = { cache: 'no-store' };
    } else {
      const separator = path.includes('?') ? '&' : '?';
      path = `${path}${separator}_t=${Date.now()}`;
      headers.set('Cache-Control', 'no-cache, no-store, must-revalidate');
      headers.set('Pragma', 'no-cache');
    }
  } else if (auth) {
    finalCacheStrategy = { cache: 'no-store' as RequestCache };
  } else {
    finalCacheStrategy = cacheStrategy;
  }

  // --- 发起请求 ---
  let res: Response;
  try {
    res = await fetch(`${GO_API_BASE}${path}`, {
      ...fetchOptions,
      headers,
      credentials,
      ...finalCacheStrategy,
    });
  } catch (err) {
    // ✅ 仅捕获网络层异常，不吞掉 AbortError 等其他类型
    if (err instanceof DOMException && err.name === 'AbortError') {
      throw err; // 请求被主动取消，原样抛出
    }
    throw new ApiConnectionError(path, err);
  }

  // 一次性读取原始文本
  let rawBodyText = '';
  try {
    rawBodyText = await res.text();
  } catch { /* ignore */ }

  // 非 2xx：解析错误信封中的 code + message
  if (!res.ok) {
    let errorMessage = `HTTP ${res.status} ${res.statusText}`;
    let errorCode: number | string | undefined;
    let errorBody: unknown = rawBodyText;

    if (rawBodyText) {
      try {
        const parsed = JSON.parse(rawBodyText) as ErrorEnvelope;
        // 👇 核心：从错误体中提取 code 和 message
        errorMessage = parsed.message || errorMessage;
        errorCode = parsed.code;
        errorBody = parsed;
      } catch {
        // 非 JSON 错误体，保留原始文本作为 message
      }
    }

    throw new ApiError(res.status, errorMessage, errorCode, errorBody);
  }

  // ✅ 2xx 成功路径：纯净信封，无 code/message/success
  if (!rawBodyText) {
    if (skipDataUnwrap) {
      return { data: null as T } as SuccessEnvelope<T>;
    }
    return null as T;
  }

  let envelope: SuccessEnvelope<T>;
  try {
    envelope = JSON.parse(rawBodyText) as SuccessEnvelope<T>;
  } catch {
    throw new ApiError(
      -2,
      'Invalid API response: failed to parse JSON body',
      undefined,
      { status: res.status, bodyPreview: rawBodyText.slice(0, 500) }
    );
  }

  if (skipDataUnwrap) {
    return envelope;
  }

  return envelope.data;
}