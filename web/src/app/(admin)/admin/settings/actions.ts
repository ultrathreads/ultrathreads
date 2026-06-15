'use server';

import { revalidatePath } from 'next/cache';
import { apiFetch, ApiError, ApiConnectionError } from '@/lib/api/client';

export interface SiteSettings {
  siteTitle: string;
  siteDescription: string;
  siteKeywords: string | null;
  siteNavs: string | null;
  defaultNodeId: number;
  recommendTags: string[];
}

/**
 * 获取站点设置
 * auth: true → 自动从 server cookie 读取 access_token
 * noCache: true → 编辑页始终拉取最新数据
 */
export async function getSettings(): Promise<SiteSettings> {
  return apiFetch<SiteSettings>('/admin/settings', {
    auth: true,
    noCache: true,
  });
}

/**
 * 保存站点设置
 * 将 FormData 转为后端期望的 JSON 格式后提交
 */
export async function saveSettings(
  formData: FormData
): Promise<{ success: boolean; message: string }> {
  const tagsRaw = (formData.get('recommendTags') as string) || '';

  const payload = {
    siteTitle: formData.get('siteTitle') as string,
    siteDescription: formData.get('siteDescription') as string,
    siteKeywords: (formData.get('siteKeywords') as string) || null,
    siteNavs: (formData.get('siteNavs') as string) || null,
    defaultNodeId: Number(formData.get('defaultNodeId')),
    recommendTags: tagsRaw
      .split(',')
      .map((t) => t.trim())
      .filter(Boolean),
  };

  try {
    await apiFetch('/admin/settings', {
      method: 'POST',
      auth: true,
      body: JSON.stringify(payload),
    });

    revalidatePath('/admin/settings');
    return { success: true, message: '设置已保存' };
  } catch (err) {
    // ✅ 复用统一的异常类，精确区分业务错误与网络故障
    if (err instanceof ApiError) {
      return { success: false, message: err.message };
    }
    if (err instanceof ApiConnectionError) {
      return { success: false, message: err.message };
    }
    return { success: false, message: '未知异常，请稍后重试' };
  }
}