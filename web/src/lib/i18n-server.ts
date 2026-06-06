// src/lib/i18n-server.ts
import { cache } from 'react';
import { headers } from 'next/headers';
import i18next from 'i18next';
import zh from '../../messages/zh.json';
import en from '../../messages/en.json';

export const getServerTranslation = cache(async (namespace?: string) => {
  const headersList = await headers();
  const locale = headersList.get('x-locale') || 'zh';

  const instance = i18next.createInstance();
  await instance.init({
    lng: locale,
    ns: namespace ? [namespace] : undefined,
    defaultNS: namespace || 'translation',
    resources: {
      zh: { common: zh.common, home: zh.home },
      en: { common: en.common, home: en.home },
    },
    fallbackLng: 'zh',
    interpolation: { escapeValue: false },
  });

  return instance.t.bind(instance);
});