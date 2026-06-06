// src/lib/api/client.ts

const GO_API_BASE = process.env.GO_API_BASE || 'http://localhost:9527/api';

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
  const {
    auth = false,
    cacheStrategy = { next: { revalidate: 60 } },
    skipDataUnwrap = false,
    ...fetchOptions
  } = options;

  const headers = new Headers(fetchOptions.headers);
  headers.set('Content-Type', 'application/json');

  // ✅ 核心修复：区分服务端与客户端的鉴权方式
  let credentials: RequestCredentials | undefined = fetchOptions.credentials;

  if (auth) {
    if (typeof window === 'undefined') {
      // 🖥️ 服务端 SSR：读取 Cookie 并手动注入 Authorization Header
      try {
        const { cookies } = await import('next/headers');
        const cookieStore = await cookies();
        const token = cookieStore.get('access_token')?.value;
        if (!token) throw new ApiBusinessError('AUTH_REQUIRED', -1);
        headers.set('Authorization', `Bearer ${token}`);
      } catch (e) {
        // 防止动态导入本身失败导致未捕获异常
        console.warn('[apiFetch] Server-side cookie read failed:', e);
        throw new ApiBusinessError('AUTH_REQUIRED', -1);
      }
    } else {
      // 🌐 客户端浏览器：依赖浏览器自动携带 Cookie
      credentials = 'include';
    }
  }

  const res = await fetch(`${GO_API_BASE}${path}`, {
    ...fetchOptions,
    headers,
    credentials,
    ...cacheStrategy,
  });

  // ✅ 精确识别 401，便于触发 Token 刷新
  if (res.status === 401) {
    throw new ApiBusinessError('UNAUTHORIZED', 401);
  }

  if (!res.ok) {
    throw new Error(`HTTP Error: ${res.status} ${res.statusText}`);
  }

  // 📦 拆解统一响应信封
  const envelope = (await res.json()) as ApiResponse<T>;

  // ✅ skipDataUnwrap = true 时直接返回完整信封
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