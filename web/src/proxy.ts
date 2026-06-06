// src/proxy.ts (或 proxy.ts)
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const DEFAULT_LOCALE = 'zh';
const SUPPORTED_LOCALES = ['zh', 'en'];

export default async function proxy(request: NextRequest) {
  // 1. 读取 Cookie 确定语言
  const cookieLocale = request.cookies.get('NEXT_LOCALE')?.value;
  const locale =
    cookieLocale && SUPPORTED_LOCALES.includes(cookieLocale)
      ? cookieLocale
      : DEFAULT_LOCALE;

  // 2. 克隆请求并注入自定义 Header（作为全链路单一数据源）
  const requestHeaders = new Headers(request.headers);
  requestHeaders.set('x-locale', locale);

  // 3. 使用新 API 重写请求头并继续路由
  return NextResponse.next({
    request: {
      headers: requestHeaders,
    },
  });
}

// Proxy 配置（替代原 middleware 的 config.matcher）
export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico|api).*)'],
};