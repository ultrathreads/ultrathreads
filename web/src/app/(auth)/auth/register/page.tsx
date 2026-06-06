// src/app/auth/register/page.tsx
import type { Metadata } from 'next';
import RegisterForm from './RegisterForm';

export const metadata: Metadata = {
  title: '创建新账号',
};

export default function RegisterPage() {
  return (
    <div className="main-body">
      <RegisterForm />
    </div>
  );
}