import { NextResponse } from 'next/server';

export async function POST() {
  const response = NextResponse.json({ success: true });

  // ✅ 核心：通过设置 maxAge: 0 和空值来强制浏览器销毁 HttpOnly Cookie
  // 注意：path、domain、sameSite 等属性必须与登录时 set cookie 的配置完全一致，否则浏览器不会匹配并删除该 Cookie
  
  response.cookies.set('access_token', '', {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'lax',
    maxAge: 0,      // 🔑 立即过期
    path: '/',      // 必须与登录时的 path 保持一致
  });

  response.cookies.set('refresh_token', '', {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'lax',
    maxAge: 0,      // 🔑 立即过期
    path: '/',      // 必须与登录时的 path 保持一致
  });

  console.log('[BFF Logout] ✅ Cookies cleared');
  return response;
}