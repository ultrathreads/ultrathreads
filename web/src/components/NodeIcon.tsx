// components/NodeIcon.tsx
import clsx from 'clsx';

// ✅ 修复：统一接口名称为 NodeIconProps
interface NodeIconProps {
  icon: string;
  className?: string;
}

// ✅ 修复：将 Props 替换为 NodeIconProps
export function NodeIcon({ icon, className }: NodeIconProps) {
  // 1. SVG 内容：直接渲染，不做任何过滤
  if (icon.trimStart().startsWith('<svg')) {
    return (
      <span
        className={clsx('inline-flex items-center justify-center w-5 h-5', className)}
        dangerouslySetInnerHTML={{ __html: icon }}
        aria-hidden="true"
      />
    );
  }

  // 2. URL：以 / 或 http 开头
  if (icon.startsWith('/') || icon.startsWith('http')) {
    return <img src={icon} alt="" className={clsx('w-5 h-5', className)} />;
  }

  // 3. 兜底：当作 Emoji 渲染
  return (
    <span role="img" aria-hidden="true" className={className}>
      {icon || '📁'}
    </span>
  );
}