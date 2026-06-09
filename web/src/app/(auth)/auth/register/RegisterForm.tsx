// src/app/auth/register/RegisterForm.tsx
'use client';
import { useState, FormEvent } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import type { RegisterFormData } from '@/types/auth';

export default function RegisterForm() {
  const router = useRouter();
  const [form, setForm] = useState<RegisterFormData>({
    username: '',
    email: '',
    password: '',
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  // 统一表单赋值
  const handleChange = (key: keyof RegisterFormData, val: string) => {
    setForm(prev => ({ ...prev, [key]: val }));
  };

  // 提交函数
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    // 模拟接口请求
    await new Promise(resolve => setTimeout(resolve, 1200));
    console.log('注册数据', form);

    setIsSubmitting(false);
    alert('🎉 注册成功！请登录。');
    router.push('/auth/login');
  };

  const { username, email, password } = form;
  const disabled = isSubmitting || !username || !email || !password;

  return (
    <div className="auth-container">
      <div className="auth-header">
        <h1 className="auth-title">创建账号</h1>
        <p className="auth-subtitle">加入 UltraThreads 社区</p>
      </div>

      <form className="auth-form" onSubmit={handleSubmit}>
        <div className="form-group">
          <label className="form-label">用户名</label>
          <input
            type="text"
            className="form-input"
            placeholder="请输入用户名..."
            required
            value={username}
            onChange={e => handleChange('username', e.target.value)}
          />
        </div>
        <div className="form-group">
          <label className="form-label">邮箱地址</label>
          <input
            type="email"
            className="form-input"
            placeholder="请输入邮箱..."
            required
            value={email}
            onChange={e => handleChange('email', e.target.value)}
          />
        </div>
        <div className="form-group">
          <label className="form-label">密码</label>
          <input
            type="password"
            className="form-input"
            placeholder="设置你的密码..."
            required
            value={password}
            onChange={e => handleChange('password', e.target.value)}
          />
        </div>
        <div className="auth-actions">
          <button type="submit" className="auth-btn" disabled={disabled}>
            {isSubmitting ? '注册中...' : '注 册'}
          </button>
        </div>
      </form>

      <div className="auth-footer">
        已有账号？
        <Link href="/auth/login" className="auth-link">返回登录</Link>
      </div>
    </div>
  );
}