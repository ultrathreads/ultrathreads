// src/app/(main)/settings/profile/page.tsx
'use client';

import { useState, useRef, useEffect, ChangeEvent, FormEvent } from 'react';
import { toast } from 'sonner';
import { getCurrentUser } from '@/services/auth.server';
import { updateUserProfile } from '@/services/user';
import { uploadAvatar } from '@/services/upload';
import { useAuth } from '@/hooks/use-auth';
import { ApiBusinessError } from '@/lib/api/client';
import type { CurrentUser } from '@/types/auth';

/** 从 CurrentUser 中提取表单所需字段，避免手动维护重复类型 */
type ProfileForm = Pick<CurrentUser, 'nickname' | 'avatar' | 'website' | 'description'>;

const DEFAULT_FORM: ProfileForm = {
  nickname: '',
  avatar: '',
  website: '',
  description: '',
};

export default function ProfilePage() {
  // ✅ 获取全局刷新方法，用于保存后同步 Header 状态
  const { refreshUser } = useAuth();

  const [userId, setUserId] = useState<number | null>(null);
  const [form, setForm] = useState<ProfileForm>(DEFAULT_FORM);
  const [previewUrl, setPreviewUrl] = useState('');
  const [uploading, setUploading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [loading, setLoading] = useState(true);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // 初始化表单并缓存用户 ID
  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const user = await getCurrentUser();
        if (cancelled) return;
        setUserId(user.id);
        setForm({
          nickname: user.nickname,
          avatar: user.avatar,
          website: user.website,
          description: user.description,
        });
        setPreviewUrl(user.avatar);
      } catch (err) {
        if (!cancelled) {
          toast.error(err instanceof ApiBusinessError ? err.message : '加载用户资料失败');
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  const handleFieldChange = (field: keyof ProfileForm) =>
    (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
      setForm((prev) => ({ ...prev, [field]: e.target.value }));

  // ========== 头像上传 ==========
  const handleFileChange = async (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.type.startsWith('image/')) {
      toast.warning('请选择图片文件');
      return;
    }
    if (file.size > 3 * 1024 * 1024) {
      toast.warning('图片大小不能超过 3MB');
      return;
    }

    const localPreview = URL.createObjectURL(file);
    setPreviewUrl(localPreview);

    try {
      setUploading(true);
      await toast.promise(uploadAvatar(file), {
        loading: '头像上传中...',
        success: (result) => {
          setForm((prev) => ({ ...prev, avatar: result.url }));
          setPreviewUrl(result.url);
          return '头像已上传，点击保存后生效';
        },
        error: (err) => {
          setPreviewUrl(form.avatar);
          return err instanceof ApiBusinessError ? err.message : '头像上传异常';
        },
      });
    } finally {
      setUploading(false);
      if (localPreview.startsWith('blob:')) URL.revokeObjectURL(localPreview);
      if (fileInputRef.current) fileInputRef.current.value = '';
    }
  };

  // ========== 表单提交 ==========
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!form.nickname.trim()) {
      toast.warning('昵称不能为空');
      return;
    }
    if (userId === null) {
      toast.error('用户信息未加载完成，请稍后重试');
      return;
    }

    try {
      setSaving(true);
      await toast.promise(
        // ✅ 核心修复：先更新后端，再刷新全局 AuthContext
        updateUserProfile(userId, form).then(async () => {
          await refreshUser();
        }),
        {
          loading: '保存中...',
          success: '个人资料已更新',
          error: (err) =>
            err instanceof ApiBusinessError ? err.message : '保存失败，请稍后重试',
        },
      );
    } finally {
      setSaving(false);
    }
  };

  const displayAvatar = previewUrl || form.avatar;

  // 加载中骨架屏，避免闪烁
  if (loading) {
    return (
      <div className="settings-card">
        <h2 className="settings-card-title">个人信息</h2>
        <p className="settings-card-desc">加载中...</p>
      </div>
    );
  }

  return (
    <div className="settings-card">
      <h2 className="settings-card-title">个人信息</h2>
      <p className="settings-card-desc">在这里管理你的公开资料，让其他用户更好地了解你。</p>

      <form className="settings-form" onSubmit={handleSubmit}>
        {/* 头像 */}
        <div className="form-group">
          <label className="form-label">头像</label>
          <div className="avatar-upload">
            <div className="avatar-preview">
              {displayAvatar ? (
                <img
                  src={displayAvatar}
                  alt="头像预览"
                  style={{
                    width: '100%',
                    height: '100%',
                    objectFit: 'cover',
                    borderRadius: '50%',
                  }}
                />
              ) : (
                <span style={{ fontSize: 32 }}>👤</span>
              )}
              {uploading && (
                <div
                  style={{
                    position: 'absolute',
                    inset: 0,
                    background: 'rgba(0,0,0,0.5)',
                    borderRadius: '50%',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    color: '#fff',
                    fontSize: 12,
                  }}
                >
                  上传中...
                </div>
              )}
            </div>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              onChange={handleFileChange}
              style={{ display: 'none' }}
            />
            <button
              type="button"
              className="upload-btn"
              onClick={() => fileInputRef.current?.click()}
              disabled={uploading}
            >
              {uploading ? '上传中...' : '更换头像'}
            </button>
          </div>
        </div>

        {/* 昵称 */}
        <div className="form-group">
          <label className="form-label">昵称</label>
          <input
            type="text"
            className="form-input"
            value={form.nickname}
            onChange={handleFieldChange('nickname')}
            maxLength={30}
          />
        </div>

        {/* 个人主页 */}
        <div className="form-group">
          <label className="form-label">个人主页</label>
          <input
            type="url"
            className="form-input"
            value={form.website}
            onChange={handleFieldChange('website')}
            placeholder="https://example.com"
          />
        </div>

        {/* 个人简介 */}
        <div className="form-group">
          <label className="form-label">个人简介</label>
          <textarea
            className="form-textarea"
            value={form.description}
            onChange={handleFieldChange('description')}
            placeholder="介绍一下你自己吧..."
            maxLength={200}
            style={{ minHeight: 'auto', height: '120px', resize: 'vertical' }}
          />
        </div>

        <button type="submit" className="save-btn" disabled={saving || uploading}>
          {saving ? '保存中...' : '保存修改'}
        </button>
      </form>
    </div>
  );
}