// components/ui/Avatar.tsx
import React, { useState } from 'react';

interface AvatarProps extends React.ImgHTMLAttributes<HTMLImageElement> {
  /** 用户头像地址 */
  src?: string | null;
  /** 默认兜底头像 */
  fallback?: string;
  /** 图片 alt 属性 */
  alt?: string;
}

// 默认 SVG 头像 (灰色小人图标，Base64 格式避免额外网络请求)
const DefaultFallbackSVG = `data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%239ca3af'%3E%3Cpath d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/%3E%3C/svg%3E`;

const Avatar: React.FC<AvatarProps> = ({ 
  src, 
  fallback = DefaultFallbackSVG, 
  alt = 'avatar', 
  className, 
  ...rest 
}) => {
  // 使用 state 跟踪图片是否加载失败，防止死循环
  const [isError, setIsError] = useState(false);

  // 核心逻辑：如果 src 为空、null、undefined，或者加载失败，则使用 fallback
  // !!src 确保了空字符串 "" 也会被判定为 false
  const currentSrc = !!src && !isError ? src : fallback;

  const handleError = () => {
    if (!isError) {
      setIsError(true);
    }
  };

  return (
    <img
      src={currentSrc}
      alt={alt}
      className={className}
      onError={handleError}
      {...rest}
    />
  );
};

export default Avatar;