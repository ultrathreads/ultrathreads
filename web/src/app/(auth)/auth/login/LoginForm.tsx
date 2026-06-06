'use client';

import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/components/providers/AuthProvider'; // ✅ 引入全局状态 Hook

// 用户名验证：3-20位，仅允许字母、数字、下划线、连字符
const USERNAME_REGEX = /^[a-zA-Z0-9_-]{3,20}$/;

export default function LoginForm() {
  const router = useRouter();
  const { refreshUser } = useAuth(); // ✅ 获取刷新用户状态的方法
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  // 实时用户名合法性校验
  const isUsernameValid = USERNAME_REGEX.test(username);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setErrorMsg(null);

    // 提交前二次校验
    if (!isUsernameValid) {
      setErrorMsg('用户名需为3-20位字母、数字、下划线或连字符');
      setIsSubmitting(false);
      return;
    }

    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });

      if (!res.ok) {
        const data = await res.json().catch(() => null);
        throw new Error(data?.error || `登录失败 (${res.status})`);
      }

      // ✅ 核心改动：登录接口返回成功后，立即拉取 /user/current 并更新全局 Context
      await refreshUser();

      // 跳转首页并刷新服务端组件缓存
      router.push('/');
      router.refresh();
    } catch (err) {
      setErrorMsg(err instanceof Error ? err.message : '未知错误');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-header">
        <h1 className="auth-title">欢迎回来</h1>
        <p className="auth-subtitle">登录你的 UltraThreads 账号</p>
      </div>

      {errorMsg && (
        <div className="auth-error" role="alert">
          {errorMsg}
        </div>
      )}

      <form className="auth-form" onSubmit={handleSubmit}>
        <div className="form-group">
          <label className="form-label">用户名</label>
          <input
            type="text"
            className="form-input"
            placeholder="请输入用户名..."
            required
            autoComplete="username"
            value={username}
            onChange={(e) => {
              setUsername(e.target.value);
              setErrorMsg(null);
            }}
          />
          {username.length > 0 && !isUsernameValid && (
            <p className="form-hint form-hint--error">
              用户名需为3-20位字母、数字、下划线或连字符
            </p>
          )}
        </div>

        <div className="form-group">
          <label className="form-label">密码</label>
          <input
            type="password"
            className="form-input"
            placeholder="请输入密码..."
            required
            autoComplete="current-password"
            value={password}
            onChange={(e) => {
              setPassword(e.target.value);
              setErrorMsg(null);
            }}
          />
          <Link href="/auth/forgot-password" className="auth-link forgot-password-link">
            忘记密码？
          </Link>
        </div>

        <div className="auth-actions">
          <button
            type="submit"
            className="auth-btn"
            disabled={isSubmitting || !isUsernameValid || !password}
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