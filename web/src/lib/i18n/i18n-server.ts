// src/lib/i18n-server.ts
import { cache } from 'react';
import { headers } from 'next/headers';
import i18next from 'i18next';
import zh from '../../../messages/zh.json';
import en from '../../../messages/en.json';

export const getServerTranslation = cache(async () => {
  const headersList = await headers();
  const locale = headersList.get('x-locale') || 'zh';

  const instance = i18next.createInstance();
  await instance.init({
    lng: locale,
    // ⚠️ 单命名空间：不再需要 ns / defaultNS 参数
    resources: {
      zh: { translation: zh },
      en: { translation: en },
    },
    fallbackLng: 'zh',
    interpolation: { escapeValue: false },
  });

  return instance.t.bind(instance);
});