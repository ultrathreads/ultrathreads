// src/app/auth/forgot-password/page.tsx
'use client';

import { useState, FormEvent } from 'react';
import Link from 'next/link';

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    // 模拟发送重置邮件的请求
    await new Promise((resolve) => setTimeout(resolve, 1200));

    console.log('发送重置密码邮件到:', email);
    setIsSubmitting(false);
    setIsSuccess(true); // 模拟发送成功，切换UI状态
  };

  return (
    <div className="main-body">
      <div className="auth-container">
        <div className="auth-header">
          <h1 className="auth-title">找回密码</h1>
          <p className="auth-subtitle">我们将帮你重置密码</p>
        </div>

        {/* 状态1：输入邮箱阶段 */}
        {!isSuccess ? (
          <form className="auth-form" onSubmit={handleSubmit}>
            <div className="auth-tip">
              请输入你注册时使用的邮箱地址，我们将向该邮箱发送包含密码重置链接的邮件。
            </div>

            <div className="form-group">
              <label className="form-label">邮箱地址</label>
              <input 
                type="email" 
                className="form-input" 
                placeholder="请输入注册邮箱..." 
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>

            <div className="auth-actions">
              <button 
                type="submit" 
                className="auth-btn" 
                disabled={isSubmitting || !email}
              >
                {isSubmitting ? '发送中...' : '发送重置链接'}
              </button>
            </div>
          </form>
        ) : (
          /* 状态2：发送成功后的提示阶段 */
          <div className="auth-form">
            <div className="auth-tip" style={{ backgroundColor: '#f0f9eb', borderColor: '#e1f3d8', color: '#67c23a' }}>
              ✉️ 重置密码链接已发送至 <strong>{email}</strong>。请检查你的收件箱（包括垃圾邮件文件夹），并按照邮件中的指引重置密码。
            </div>
            <div className="auth-actions">
              <Link href="/auth/login" className="auth-btn" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', textDecoration: 'none' }}>
                返回登录
              </Link>
            </div>
          </div>
        )}

        <div className="auth-footer">
          想起密码了？<Link href="/auth/login" className="auth-link">返回登录</Link>
        </div>
      </div>
    </div>
  );
}