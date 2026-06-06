// src/app/auth/login/page.tsx
import type { Metadata } from 'next';
import LoginForm from './LoginForm';

// 👇 静态 metadata，构建时即可确定，完美支持 SEO
export const metadata: Metadata = {
  title: '用户登录',
};

// 👇 极简的服务端页面组件
export default function LoginPage() {
  return (
    <div className="main-body">
      <LoginForm />
    </div>
  );
}