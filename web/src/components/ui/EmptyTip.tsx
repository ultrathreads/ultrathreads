// components/ui/EmptyTip.tsx
'use client';

type TipVariant = 'empty' | 'error';

interface EmptyTipProps {
  text: string;
  variant?: TipVariant;
  className?: string;
}

export default function EmptyTip({
  text,
  variant = 'empty',
  className = '',
}: EmptyTipProps) {
  const baseClass = variant === 'error' ? 'error-tip' : 'empty-tip';
  return <div className={`${baseClass} ${className}`}>{text}</div>;
}