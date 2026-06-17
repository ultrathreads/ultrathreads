// src/lib/i18n-client.ts
'use client';

import i18next from 'i18next';
import { initReactI18next, useTranslation as useTranslationOrg } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

// ⚠️ 与服务端保持一致：直接导入完整的翻译 JSON
import zh from '../../../messages/zh.json';
import en from '../../../messages/en.json';

export const languages = ['zh', 'en'] as const;
export type Locale = (typeof languages)[number];
export const defaultLocale: Locale = 'zh';

const runsOnServerSide = typeof window === 'undefined';

i18next
  .use(initReactI18next)
  .use(LanguageDetector)
  .init({
    // ⚠️ 单命名空间：将所有 key 合并到默认的 translation 命名空间中
    resources: {
      zh: { translation: zh },
      en: { translation: en },
    },
    supportedLngs: languages,
    fallbackLng: defaultLocale,
    lng: defaultLocale,
    // 不再需要指定 ns 和 defaultNS，i18next 默认使用 'translation'
    interpolation: {
      escapeValue: false,
    },
    detection: {
      order: ['cookie', 'navigator'],
      caches: ['cookie'],
      lookupCookie: 'NEXT_LOCALE',
      cookieOptions: {
        path: '/',
        sameSite: 'lax',
        maxAge: 365 * 24 * 60 * 60,      // 1年
      },
    },
    // SSR 阶段固定使用默认语言，防止服务端与客户端初始渲染不一致
    ...(runsOnServerSide && {
      lng: defaultLocale,
      detection: { order: [] },
    }),
  });

/**
 * 类型安全的客户端翻译 Hook
 * 用法: const { t } = useTranslation();
 */
export function useTranslation(ns?: string | string[]) {
  return useTranslationOrg(ns);
}

export default i18next;