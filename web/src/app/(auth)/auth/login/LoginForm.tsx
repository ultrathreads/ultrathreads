'use client';

import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { toast } from 'sonner';
import { useAuth } from '@/hooks/use-auth';

const USERNAME_REGEX = /^[a-zA-Z0-9_-]{3,20}$/;

export default function LoginForm() {
  const router = useRouter();
  const { refreshUser } = useAuth();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  // ✅ 分离两种错误状态
  const [fieldErrors, setFieldErrors] = useState<{ username?: string; password?: string }>({});
  const [serverError, setServerError] = useState<string | null>(null);

  // ✅ 实时校验：仅在用户输入过时才显示格式错误
  const usernameTouched = username.length > 0;
  const isUsernameFormatValid = USERNAME_REGEX.test(username);

  const validateFields = (): boolean => {
    const errors: typeof fieldErrors = {};
    if (!username.trim()) errors.username = '请输入用户名';
    else if (!isUsernameFormatValid) errors.username = '用户名需为3-20位字母、数字、下划线或连字符';
    if (!password) errors.password = '请输入密码';

    setFieldErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setServerError(null);

    // ✅ 前端校验失败 → 仅行内提示，不弹 Toast
    if (!validateFields()) return;

    setIsSubmitting(true);
    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });

      if (!res.ok) {
        const data = await res.json().catch(() => null);
        // ✅ 服务端业务错误 → 顶部横幅，持久展示
        throw new Error(data?.error || `登录失败 (${res.status})`);
      }

      await refreshUser();
      // ✅ 成功跳转用 Toast 确认（可选，也可省略直接跳转）
      toast.success('登录成功，欢迎回来 👋');
      router.push('/');
      router.refresh();
    } catch (err) {
      const message = err instanceof Error ? err.message : '未知错误';
      // ✅ 区分业务错误和网络异常
      if (message.includes('网络') || message.includes('timeout') || message.includes('500')) {
        toast.error(message);
      } else {
        setServerError(message);
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  // ✅ 输入时清除对应字段的行内错误 + 服务端错误
  const handleUsernameChange = (val: string) => {
    setUsername(val);
    setServerError(null);
    if (fieldErrors.username) setFieldErrors((prev) => ({ ...prev, username: undefined }));
  };

  const handlePasswordChange = (val: string) => {
    setPassword(val);
    setServerError(null);
    if (fieldErrors.password) setFieldErrors((prev) => ({ ...prev, password: undefined }));
  };

  const canSubmit = !isSubmitting && username.trim() && password && isUsernameFormatValid;

  return (
    <div className="auth-container">
      <div className="auth-header">
        <h1 className="auth-title">欢迎回来</h1>
        <p className="auth-subtitle">登录你的 UltraThreads 账号</p>
      </div>

      {/* ✅ 服务端错误横幅：持久展示，role="alert" 保证无障碍可读 */}
      {serverError && (
        <div className="auth-error" role="alert">
          {serverError}
        </div>
      )}

      <form className="auth-form" onSubmit={handleSubmit} noValidate>
        <div className="form-group">
          <label className="form-label">用户名</label>
          <input
            type="text"
            className={`form-input ${fieldErrors.username ? 'form-error' : ''}`}
            placeholder="请输入用户名..."
            autoComplete="username"
            value={username}
            onChange={(e) => handleUsernameChange(e.target.value)}
            aria-invalid={!!fieldErrors.username}
            aria-describedby={fieldErrors.username ? 'username-error' : undefined}
          />
          {/* ✅ 行内错误：格式校验失败时显示 */}
          {fieldErrors.username && (
            <p id="username-error" className="form-hint form-hint--error" role="alert">
              {fieldErrors.username}
            </p>
          )}
          {/* ✅ 未出错时的格式提示（仅在用户开始输入且格式正确时隐藏） */}
          {!fieldErrors.username && usernameTouched && !isUsernameFormatValid && (
            <p className="form-hint form-hint--error">
              用户名需为3-20位字母、数字、下划线或连字符
            </p>
          )}
        </div>

        <div className="form-group">
          <label className="form-label">密码</label>
          <input
            type="password"
            className={`form-input ${fieldErrors.password ? 'form-error' : ''}`}
            placeholder="请输入密码..."
            autoComplete="current-password"
            value={password}
            onChange={(e) => handlePasswordChange(e.target.value)}
            aria-invalid={!!fieldErrors.password}
            aria-describedby={fieldErrors.password ? 'password-error' : undefined}
          />
          {fieldErrors.password && (
            <p id="password-error" className="form-hint form-hint--error" role="alert">
              {fieldErrors.password}
            </p>
          )}
          <Link href="/auth/forgot-password" className="auth-link forgot-password-link">
            忘记密码？
          </Link>
        </div>

        <div className="auth-actions">
          <button
            type="submit"
            className="auth-btn"
            disabled={!canSubmit}
            aria-disabled={!canSubmit}
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