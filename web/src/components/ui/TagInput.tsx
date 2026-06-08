// src/components/ui/TagInput.tsx
'use client';

import { useState, useRef, useEffect, KeyboardEvent, ChangeEvent } from 'react';
// ✅ 改为从项目统一的 service 层导入，类型也复用 TagEntity
import { fetchTagSuggestions, type TagEntity } from '@/services/tag-service';

interface TagInputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
}

export function TagInput({ value, onChange, placeholder, className }: TagInputProps) {
  // ✅ 类型同步更新为 TagEntity
  const [suggestions, setSuggestions] = useState<TagEntity[]>([]);
  const [showDropdown, setShowDropdown] = useState(false);
  const [activeIndex, setActiveIndex] = useState(-1);
  const debounceTimer = useRef<NodeJS.Timeout | null>(null);
  const wrapperRef = useRef<HTMLDivElement>(null);

  const getCurrentInput = () => {
    const parts = value.split(',');
    return parts[parts.length - 1]?.trim() || '';
  };

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    onChange(newValue);

    const currentWord = newValue.split(',').pop()?.trim() || '';

    if (debounceTimer.current) clearTimeout(debounceTimer.current);

    if (!currentWord) {
      setSuggestions([]);
      setShowDropdown(false);
      return;
    }

    debounceTimer.current = setTimeout(async () => {
      // ✅ 调用统一 service，内部已处理异常兜底
      const results = await fetchTagSuggestions(currentWord);
      setSuggestions(results);
      setShowDropdown(results.length > 0);
      setActiveIndex(-1);
    }, 300);
  };

  const selectSuggestion = (tagName: string) => {
    const parts = value.split(',');
    parts[parts.length - 1] = ` ${tagName}`;
    const nextValue = parts.join(',') + ', ';
    onChange(nextValue.trim().replace(/,\s*$/, ''));

    setSuggestions([]);
    setShowDropdown(false);
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (!showDropdown) return;

    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setActiveIndex((prev) => (prev < suggestions.length - 1 ? prev + 1 : 0));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setActiveIndex((prev) => (prev > 0 ? prev - 1 : suggestions.length - 1));
    } else if (e.key === 'Enter' && activeIndex >= 0) {
      e.preventDefault();
      selectSuggestion(suggestions[activeIndex].name);
    } else if (e.key === 'Escape') {
      setShowDropdown(false);
    }
  };

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (wrapperRef.current && !wrapperRef.current.contains(e.target as Node)) {
        setShowDropdown(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <div ref={wrapperRef} style={{ position: 'relative' }}>
      <input
        type="text"
        className={className}
        placeholder={placeholder}
        value={value}
        onChange={handleInputChange}
        onKeyDown={handleKeyDown}
        onFocus={() => suggestions.length > 0 && setShowDropdown(true)}
        autoComplete="off"
      />

      {showDropdown && (
        <ul
          style={{
            position: 'absolute',
            top: '100%',
            left: 0,
            right: 0,
            zIndex: 50,
            margin: '4px 0 0',
            padding: 0,
            listStyle: 'none',
            background: '#fff',
            border: '1px solid #e2e8f0',
            borderRadius: 6,
            boxShadow: '0 4px 6px -1px rgba(0,0,0,.1)',
            maxHeight: 200,
            overflowY: 'auto',
          }}
        >
          {suggestions.map((tag, idx) => (
            <li
              key={tag.id}
              onClick={() => selectSuggestion(tag.name)}
              style={{
                padding: '8px 12px',
                cursor: 'pointer',
                fontSize: 14,
                background: idx === activeIndex ? '#edf2f7' : 'transparent',
              }}
              onMouseEnter={() => setActiveIndex(idx)}
            >
              {tag.name}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}