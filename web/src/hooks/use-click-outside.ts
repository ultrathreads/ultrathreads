// src/hooks/use-click-outside.ts
import { useEffect, useRef, useState } from 'react';

/**
 * 通用点击外部关闭 Hook
 * @param initialOpen - 初始展开状态，默认 false
 */
export function useClickOutside(initialOpen = false) {
  const [isOpen, setIsOpen] = useState(initialOpen);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!isOpen) return;

    function handleClickOutside(event: MouseEvent) {
      if (ref.current && !ref.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isOpen]);

  const toggle = () => setIsOpen(prev => !prev);
  const close = () => setIsOpen(false);
  const open = () => setIsOpen(true);

  return { isOpen, toggle, close, open, ref };
}