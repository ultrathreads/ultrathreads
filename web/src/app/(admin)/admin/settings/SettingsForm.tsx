'use client';

import { useTransition } from 'react';
import { toast } from 'sonner';
import type { SiteSettings } from './actions';

interface Props {
  initialData: SiteSettings;
  saveAction: (formData: FormData) => Promise<{ success: boolean; message: string }>;
}

export function SettingsForm({ initialData, saveAction }: Props) {
  const [isPending, startTransition] = useTransition();

  const handleSubmit = (formData: FormData) => {
    startTransition(async () => {
      const result = await saveAction(formData);
      if (result.success) {
        toast.success(result.message);
      } else {
        toast.error(result.message);
      }
    });
  };

  return (
    <form action={handleSubmit} className="admin-settings-form">
      {/* 站点标题 */}
      <div className="admin-form-group">
        <label className="admin-form-label">
          站点标题 <span className="required">*</span>
        </label>
        <input
          name="siteTitle"
          className="admin-form-input"
          defaultValue={initialData.siteTitle}
          required
          placeholder="例如：UltraThreads"
        />
      </div>

      {/* 站点描述 */}
      <div className="admin-form-group">
        <label className="admin-form-label">站点描述</label>
        <textarea
          name="siteDescription"
          className="admin-form-textarea"
          defaultValue={initialData.siteDescription}
          rows={3}
          placeholder="一句话描述你的站点"
        />
      </div>

      {/* SEO 关键词 */}
      <div className="admin-form-group">
        <label className="admin-form-label">SEO 关键词</label>
        <input
          name="siteKeywords"
          className="admin-form-input"
          defaultValue={initialData.siteKeywords ?? ''}
          placeholder="多个关键词用英文逗号分隔"
        />
        <p className="admin-form-hint">用于搜索引擎优化，留空则不输出 meta keywords</p>
      </div>

      {/* 双列布局 */}
      <div className="admin-form-row">
        <div className="admin-form-group">
          <label className="admin-form-label">
            默认节点 ID <span className="required">*</span>
          </label>
          <input
            name="defaultNodeId"
            type="number"
            className="admin-form-input"
            defaultValue={initialData.defaultNodeId}
            required
            min={1}
          />
          <p className="admin-form-hint">发帖时未选择板块时的兜底节点</p>
        </div>

        <div className="admin-form-group">
          <label className="admin-form-label">自定义导航 (JSON)</label>
          <textarea
            name="siteNavs"
            className="admin-form-textarea admin-form-textarea-mono"
            defaultValue={initialData.siteNavs ?? ''}
            rows={4}
            placeholder='[{"label":"首页","href":"/"}]'
          />
          <p className="admin-form-hint">留空使用默认导航，填写需为合法 JSON 数组</p>
        </div>
      </div>

      {/* 推荐标签 */}
      <div className="admin-form-group">
        <label className="admin-form-label">推荐标签</label>
        <input
          name="recommendTags"
          className="admin-form-input"
          defaultValue={initialData.recommendTags?.join(', ') ?? ''}
          placeholder="多个标签用英文逗号分隔"
        />
        <p className="admin-form-hint">
          当前共 {initialData.recommendTags?.length ?? 0} 个标签，编辑后以逗号分隔保存
        </p>
      </div>

      {/* 提交按钮 */}
      <div className="admin-form-footer">
        <button
          type="submit"
          className="admin-btn admin-btn-primary"
          disabled={isPending}
        >
          {isPending ? '保存中...' : '保存设置'}
        </button>
      </div>
    </form>
  );
}