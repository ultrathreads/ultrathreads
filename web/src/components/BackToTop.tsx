'use client';
import { useState, useEffect } from 'react';

export default function BackToTop() {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const el = document.getElementById('mainContent');
    if (!el) return;
    const handler = () => setVisible(el.scrollTop > 300);
    el.addEventListener('scroll', handler, { passive: true });
    return () => el.removeEventListener('scroll', handler);
  }, []);

  return (
    <button
      id="back-to-top"
      className={visible ? 'visible' : ''}
      onClick={() => document.getElementById('mainContent')?.scrollTo({ top: 0, behavior: 'smooth' })}
    >
      ↑
    </button>
  );
}