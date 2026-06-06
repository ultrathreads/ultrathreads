// src/app/api/auth/login/route.ts
import { NextRequest, NextResponse } from 'next/server';
import { login, type LoginParams, type LoginEnvelope } from '@/services/auth-service';
import { ApiBusinessError } from '@/lib/api/client';

export async function POST(req: NextRequest) {
  try {
    const body = (await req.json()) as LoginParams;

    // 🔍 诊断日志 1：确认请求参数
    console.log('[BFF Login] Request params:', { username: body.username });

    const result = await login(body);

    // 🔍 诊断日志 2：确认 skipDataUnwrap 是否生效（关键！）
    // 如果这里打印出 undefined，说明 client.ts 的 skipDataUnwrap 没起作用
    console.log('[BFF Login] Raw result success field:', result?.success);
    console.log('[BFF Login] Raw result keys:', Object.keys(result || {}));

    // ✅ 业务失败判断
    if (!result?.success) {
      console.warn('[BFF Login] Business failed:', result);
      return NextResponse.json(
        { error: result?.message || '登录失败', code: result?.code },
        { status: 400 }
      );
    }

    // ✅ Token 防御性校验
    const { access_token, refresh_token, expires_in } = result.data || {};
    if (!access_token) {
      console.error('[BFF Login] Token missing despite success=true:', result);
      return NextResponse.json(
        { error: '登录凭证获取异常', code: -1 },
        { status: 500 }
      );
    }

    // ✅ 成功响应
    const response = NextResponse.json({ success: true });

    response.cookies.set('access_token', access_token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: expires_in || 86400,
      path: '/',
    });

    if (refresh_token) {
      response.cookies.set('refresh_token', refresh_token, {
        httpOnly: true,
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'lax',
        maxAge: 60 * 60 * 24 * 30,
        path: '/',
      });
    }

    console.log('[BFF Login] ✅ Success, cookie set');
    return response;

  } catch (error) {
    // 🔍 诊断日志 3：确认是否是 apiFetch 抛出的异常
    if (error instanceof ApiBusinessError) {
      console.error('[BFF Login] ❌ ApiBusinessError caught:', {
        message: error.message,
        code: error.code,
        raw: error,
      });
      return NextResponse.json(
        { error: error.message, code: error.code },
        { status: 400 }
      );
    }

    console.error('[BFF Login] ❌ Unexpected Error:', error);
    return NextResponse.json(
      { error: '登录服务内部异常' },
      { status: 500 }
    );
  }
}