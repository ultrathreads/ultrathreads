'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { useAuth } from '@/hooks/use-auth';
import { loginClient } from '@/services/auth.client';
import AuthForm, { type AuthFieldConfig } from './AuthForm';

const USERNAME_REGEX = /^[a-zA-Z0-9_-]{3,20}$/;

const LOGIN_FIELDS: AuthFieldConfig[] = [
  {
    name: 'username',
    label: '用户名',
    placeholder: '请输入用户名...',
    autoComplete: 'username',
    required: true,
    validate: (v) =>
      !USERNAME_REGEX.test(v)
        ? '用户名需为3-20位字母、数字、下划线或连字符'
        : undefined,
  },
  {
    name: 'password',
    label: '密码',
    type: 'password',
    placeholder: '请输入密码...',
    autoComplete: 'current-password',
    required: true,
  },
];

export default function LoginForm() {
  const router = useRouter();
  const { refreshUser } = useAuth();

  const handleSubmit = async (values: Record<string, string>) => {
    await loginClient({
      username: values.username,
      password: values.password,
    });

    await refreshUser();
    toast.success('登录成功，欢迎回来 👋');
    router.push('/');
    router.refresh();
  };

  return (
    <AuthForm
      title="欢迎回来"
      subtitle="登录你的 UltraThreads 账号"
      fields={LOGIN_FIELDS}
      submitLabel="登 录"
      submittingLabel="登录中..."
      onSubmit={handleSubmit}
      renderExtraContent={(fieldName) =>
        fieldName === 'password' ? (
          <Link href="/auth/forgot-password" className="auth-link forgot-password-link">
            忘记密码？
          </Link>
        ) : null
      }
      footer={
        <>
          还没有账号？
          <Link href="/auth/register" className="auth-link">立即注册</Link>
        </>
      }
    />
  );
}