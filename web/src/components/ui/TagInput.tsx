'use client';

import { useState, useRef, useEffect, useCallback, KeyboardEvent } from 'react';
import { fetchTagSuggestions, type TagEntity } from '@/services/tag-service';

// 🟢 内联关闭图标，替代 lucide-react 的 <X />
function CloseIcon({ size = 12, className }: { size?: number; className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
    >
      <path d="M18 6 6 18" />
      <path d="m6 6 12 12" />
    </svg>
  );
}

interface TagInputProps {
  value: string[];
  onChange: (value: string[]) => void;
  placeholder?: string;
  recommendTags?: string[];
  maxTags?: number;
  maxTagLength?: number;
  className?: string;
}

export function TagInput({
  value,
  onChange,
  placeholder = '输入标签后按回车添加',
  recommendTags,
  maxTags = 5,
  maxTagLength = 20,
  className,
}: TagInputProps) {
  const [input, setInput] = useState('');
  const [suggestions, setSuggestions] = useState<TagEntity[]>([]);
  const [showRecommend, setShowRecommend] = useState(false);
  const [activeIndex, setActiveIndex] = useState(-1);
  const [isLoading, setIsLoading] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const wrapperRef = useRef<HTMLDivElement>(null);
  const debounceTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const closeTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  const isMaxReached = value.length >= maxTags;
  const showSuggestionDropdown = suggestions.length > 0 && !isMaxReached;
  const showRecommendPanel =
    !showSuggestionDropdown && showRecommend && !!recommendTags?.length && !isMaxReached;

  // ---------- 核心操作 ----------

  const addTag = useCallback(
    (tagName: string) => {
      const normalized = tagName.trim().replace(/^[,;]+|[,;]+$/g, '');
      if (
        !normalized ||
        normalized.length > maxTagLength ||
        value.length >= maxTags ||
        value.includes(normalized)
      ) {
        return;
      }
      onChange([...value, normalized]);
      setInput('');
      setSuggestions([]);
      setActiveIndex(-1);
    },
    [value, onChange, maxTags, maxTagLength]
  );

  const removeTag = useCallback(
    (tagName: string) => {
      onChange(value.filter((t) => t !== tagName));
    },
    [value, onChange]
  );

  // ---------- 防抖搜索 ----------

  const handleInputChange = useCallback(
    (nextValue: string) => {
      setInput(nextValue);
      setShowRecommend(false);
      setActiveIndex(-1);

      if (debounceTimer.current) clearTimeout(debounceTimer.current);

      if (!nextValue.trim()) {
        setSuggestions([]);
        return;
      }

      setIsLoading(true);
      debounceTimer.current = setTimeout(async () => {
        try {
          const results = await fetchTagSuggestions(nextValue.trim());
          setSuggestions(results.filter((t) => !value.includes(t.name)));
        } catch {
          setSuggestions([]);
        } finally {
          setIsLoading(false);
        }
      }, 300);
    },
    [value]
  );

  // ---------- 键盘交互 ----------

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if ((e.key === 'Backspace' || e.key === 'Delete') && !input && value.length > 0) {
      e.preventDefault();
      removeTag(value[value.length - 1]);
      return;
    }

    if (showSuggestionDropdown) {
      if (e.key === 'ArrowDown') {
        e.preventDefault();
        setActiveIndex((prev) => (prev < suggestions.length - 1 ? prev + 1 : prev));
        return;
      }
      if (e.key === 'ArrowUp') {
        e.preventDefault();
        setActiveIndex((prev) => (prev > 0 ? prev - 1 : -1));
        return;
      }
      if (e.key === 'Enter' && activeIndex >= 0) {
        e.preventDefault();
        addTag(suggestions[activeIndex].name);
        return;
      }
    }

    if (e.key === 'Enter' || e.key === ',' || e.key === ';') {
      e.preventDefault();
      addTag(input);
      return;
    }

    if (e.key === 'Escape') {
      setSuggestions([]);
      setActiveIndex(-1);
      setShowRecommend(false);
    }
  };

  // ---------- 点击外部关闭 ----------

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (wrapperRef.current && !wrapperRef.current.contains(e.target as Node)) {
        setShowRecommend(false);
        setSuggestions([]);
        setActiveIndex(-1);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // ---------- 清理定时器 ----------

  useEffect(() => {
    return () => {
      if (debounceTimer.current) clearTimeout(debounceTimer.current);
      if (closeTimer.current) clearTimeout(closeTimer.current);
    };
  }, []);

  // ---------- 推荐面板开关 ----------

  const openRecommend = () => {
    if (closeTimer.current) {
      clearTimeout(closeTimer.current);
      closeTimer.current = null;
    }
    if (recommendTags?.length && !isMaxReached) setShowRecommend(true);
  };

  const closeRecommend = () => {
    closeTimer.current = setTimeout(() => setShowRecommend(false), 150);
  };

  // ---------- 渲染 ----------

  return (
    <div ref={wrapperRef} className={`tag-input-wrapper ${className || ''}`}>
      <input type="hidden" name="tags" value={value.join(',')} readOnly />

      <div
        className={`tag-input-container ${isMaxReached ? 'is-disabled' : ''}`}
        onMouseDown={(e) => {
          if (e.target === e.currentTarget) {
            e.preventDefault();
            inputRef.current?.focus();
          }
        }}
      >
        {value.map((tag) => (
          <span key={tag} className="tag-input-badge">
            {tag}
            <button
              type="button"
              className="tag-input-badge-remove"
              onClick={() => removeTag(tag)}
            >
              {/* 🟢 替换为内联 SVG */}
              <CloseIcon size={12} />
            </button>
          </span>
        ))}

        <input
          ref={inputRef}
          className="tag-input-field"
          value={input}
          placeholder={isMaxReached ? '' : placeholder}
          disabled={isMaxReached}
          maxLength={maxTagLength}
          onChange={(e) => handleInputChange(e.target.value)}
          onKeyDown={handleKeyDown}
          onFocus={openRecommend}
          onBlur={closeRecommend}
          autoComplete="off"
        />

        {isLoading && <span className="tag-input-spinner" />}
      </div>

      {/* 自动补全下拉 */}
      {showSuggestionDropdown && (
        <div className="tag-input-dropdown">
          {suggestions.map((tag, idx) => (
            <button
              key={tag.id}
              type="button"
              className={`tag-input-option ${idx === activeIndex ? 'is-active' : ''}`}
              onMouseDown={(e) => e.preventDefault()}
              onMouseEnter={() => setActiveIndex(idx)}
              onClick={() => addTag(tag.name)}
            >
              {tag.name}
            </button>
          ))}
        </div>
      )}

      {/* 推荐标签面板 */}
      {showRecommendPanel && (
        <div className="tag-input-dropdown">
          <div className="tag-input-recommend-header">
            <span className="tag-input-recommend-title">推荐标签</span>
            <button
              type="button"
              className="tag-input-recommend-close"
              onMouseDown={(e) => e.preventDefault()}
              onClick={() => setShowRecommend(false)}
            >
              {/* 🟢 替换为内联 SVG */}
              <CloseIcon size={14} />
            </button>
          </div>
          <div className="tag-input-recommend-list">
            {recommendTags!.map((tag) => (
              <button
                key={tag}
                type="button"
                className="tag-input-recommend-item"
                onMouseDown={(e) => e.preventDefault()}
                onClick={() => {
                  addTag(tag);
                  setShowRecommend(false);
                }}
              >
                {tag}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}