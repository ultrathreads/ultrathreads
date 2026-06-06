'use client';
import { useEffect } from 'react';

export default function PreviewPopover() {
  useEffect(() => {
    // 复用原始 JS 中的锚定预览气泡逻辑
    // 由于是纯客户端 DOM 操作，直接在 useEffect 中注入
    let popover: HTMLDivElement | null = null;
    let currentBtn: HTMLElement | null = null;

    function ensurePopover() {
      if (popover) return popover!;
      popover = document.createElement('div');
      popover.className = 'preview-popover';
      popover.innerHTML = `
        <div class="preview-popover-header">
          <span class="preview-popover-title"></span>
          <button class="preview-popover-close" aria-label="关闭">✕</button>
        </div>
        <div class="preview-popover-body"></div>`;
      document.body.appendChild(popover);
      popover.querySelector('.preview-popover-close')!.addEventListener('click', close);
      document.addEventListener('click', (e) => {
        if (popover?.classList.contains('open') &&
            !popover.contains(e.target as Node) &&
            !(e.target as HTMLElement).closest('.preview-btn')) close();
      });
      document.addEventListener('keydown', (e) => { if (e.key === 'Escape') close(); });
      window.addEventListener('resize', reposition);
      window.addEventListener('scroll', reposition, true);
      return popover;
    }

    function reposition() {
      if (!popover?.classList.contains('open') || !currentBtn) return;
      position(currentBtn);
    }

    function position(btn: HTMLElement) {
      const rect = btn.getBoundingClientRect();
      const popW = popover!.offsetWidth, popH = popover!.offsetHeight, gap = 10;
      let left = rect.left + window.scrollX;
      if (left + popW > window.innerWidth - 16) left = window.innerWidth - popW - 16 + window.scrollX;
      if (left < 16 + window.scrollX) left = 16 + window.scrollX;
      let top = rect.bottom + window.scrollY + gap;
      popover!.classList.remove('flip-y');
      if (top + popH > window.innerHeight + window.scrollY - 16) {
        top = rect.top + window.scrollY - popH - gap;
        popover!.classList.add('flip-y');
      }
      popover!.style.left = left + 'px';
      popover!.style.top = top + 'px';
    }

    function open(btn: HTMLElement) {
      const pop = ensurePopover();
      pop.querySelector('.preview-popover-title')!.textContent = btn.dataset.title || '无标题';
      pop.querySelector('.preview-popover-body')!.textContent = '演示内容：' + (btn.dataset.title || '');
      currentBtn = btn;
      pop.classList.remove('open');
      position(btn);
      void pop.offsetHeight;
      pop.classList.add('open');
    }

    function close() {
      if (!popover) return;
      popover.classList.remove('open');
      currentBtn = null;
    }

    const captureHandler = (e: MouseEvent) => {
      const btn = (e.target as HTMLElement).closest('.preview-btn') as HTMLElement;
      if (!btn) return;
      e.preventDefault();
      e.stopPropagation();
      if (currentBtn === btn && popover?.classList.contains('open')) { close(); return; }
      open(btn);
    };
    document.addEventListener('click', captureHandler, true);

    return () => {
      document.removeEventListener('click', captureHandler, true);
      popover?.remove();
    };
  }, []);

  return null;
}