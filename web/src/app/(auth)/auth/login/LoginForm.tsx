// src/app/auth/login/LoginForm.tsx
'use client';

import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function LoginForm() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    // 模拟登录请求
    await new Promise((resolve) => setTimeout(resolve, 1200));

    console.log('登录数据:', { email, password });
    setIsSubmitting(false);
    alert('🎉 登录成功！');
    router.push('/'); // 登录成功后跳转回首页
  };

  return (
    <div className="auth-container">
      <div className="auth-header">
        <h1 className="auth-title">欢迎回来</h1>
        <p className="auth-subtitle">登录你的 UltraThreads 账号</p>
      </div>

      <form className="auth-form" onSubmit={handleSubmit}>
        <div className="form-group">
          <label className="form-label">邮箱地址</label>
          <input 
            type="email" 
            className="form-input" 
            placeholder="请输入邮箱..." 
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>

        <div className="form-group">
          <label className="form-label">密码</label>
          <input 
            type="password" 
            className="form-input" 
            placeholder="请输入密码..." 
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          {/* 👇 忘记密码链接 */}
          <Link href="/auth/forgot-password" className="auth-link forgot-password-link">
            忘记密码？
          </Link>
        </div>

        <div className="auth-actions">
          <button 
            type="submit" 
            className="auth-btn" 
            disabled={isSubmitting || !email || !password}
          >
            {isSubmitting ? '登录中...' : '登 录'}
          </button>
        </div>
      </form>

      <div className="auth-footer">
        还没有账号？<Link href="/auth/register" className="auth-link">立即注册</Link>
      </div>
    </div>
  );
}