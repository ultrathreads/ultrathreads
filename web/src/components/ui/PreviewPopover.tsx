'use client';
import { useEffect } from 'react';
import { apiFetch, ApiBusinessError } from '@/lib/api/client';
import type { PostEntity } from '@/types/domain'; // 👈 直接复用已有领域类型
import DOMPurify from 'dompurify';

export default function PreviewPopover() {
  useEffect(() => {
    let popover: HTMLDivElement | null = null;
    let currentBtn: HTMLElement | null = null;
    let abortController: AbortController | null = null;

    function ensurePopover() {
      if (popover) return popover!;
      popover = document.createElement('div');
      popover.className = 'preview-popover';
      popover.innerHTML = `
        <div class="preview-popover-header">
          <span class="preview-popover-title"></span>
          <button class="preview-popover-close" aria-label="关闭">✕</button>
        </div>
        <div class="preview-popover-body">
          <div class="preview-meta" style="display:none;font-size:12px;color:#7f8c8d;margin-bottom:8px"></div>
          <div class="preview-loading">加载中…</div>
          <div class="preview-content post-body" style="display:none"></div>
          <div class="preview-error" style="display:none"></div>
        </div>`;
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

    async function open(btn: HTMLElement) {
      const pop = ensurePopover();
      const postSlug = btn.dataset.postSlug;
      const fallbackTitle = btn.dataset.title || '无标题';

      // UI 重置
      pop.querySelector('.preview-popover-title')!.textContent = fallbackTitle;
      const metaEl = pop.querySelector('.preview-meta') as HTMLDivElement;
      const loadingEl = pop.querySelector('.preview-loading') as HTMLDivElement;
      const contentEl = pop.querySelector('.preview-content') as HTMLDivElement;
      const errorEl = pop.querySelector('.preview-error') as HTMLDivElement;
      metaEl.style.display = 'none';
      loadingEl.style.display = '';
      contentEl.style.display = 'none';
      errorEl.style.display = 'none';

      currentBtn = btn;
      pop.classList.remove('open');
      position(btn);
      void pop.offsetHeight;
      pop.classList.add('open');

      abortController?.abort();
      abortController = new AbortController();

      if (!postSlug) {
        loadingEl.style.display = 'none';
        errorEl.textContent = '缺少帖子 Slug';
        errorEl.style.display = '';
        return;
      }

      try {
        // 👇 复用 apiFetch + PostEntity 类型，路径 /post/{slug}
        const data = await apiFetch<PostEntity>(`/post/${postSlug}`, {
          signal: abortController.signal,
          cacheStrategy: { cache: 'no-store' },
        });

        if (abortController.signal.aborted) return;

        // 更新标题 & 元信息（安全访问可选字段）
        pop.querySelector('.preview-popover-title')!.textContent = data.title;
        const metaParts = [
          data.user?.nickname,
          data.node?.name,
          `👁 ${data.viewCount ?? 0}`,
        ].filter(Boolean);
        metaEl.textContent = metaParts.join(' · ');
        metaEl.style.display = '';

        // 👇 安全渲染 HTML（content 在 PostEntity 中是可选字段）
        const rawHtml = data.content ?? '<p style="color:#999">暂无内容</p>';
        contentEl.innerHTML = DOMPurify.sanitize(rawHtml, {
          ALLOWED_TAGS: ['h1','h2','h3','h4','p','ul','ol','li','code','pre',
                         'blockquote','em','strong','a','br','img'],
          ALLOWED_ATTR: ['href', 'title', 'rel', 'src', 'alt', 'id'],
        });

        loadingEl.style.display = 'none';
        contentEl.style.display = '';
        position(btn);
      } catch (err: unknown) {
        if (abortController.signal.aborted) return;
        loadingEl.style.display = 'none';
        if (err instanceof ApiBusinessError) {
          errorEl.textContent = err.message;
        } else if (err instanceof DOMException && err.name === 'AbortError') {
          return;
        } else {
          errorEl.textContent = err instanceof Error ? err.message : '加载失败';
        }
        errorEl.style.display = '';
      }
    }

    function close() {
      if (!popover) return;
      popover.classList.remove('open');
      currentBtn = null;
      abortController?.abort();
      abortController = null;
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
      abortController?.abort();
      popover?.remove();
    };
  }, []);

  return null;
}