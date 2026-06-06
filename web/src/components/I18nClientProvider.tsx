// src/components/I18nClientProvider.tsx
'use client';

import '@/lib/i18n-client';
import { useEffect } from 'react';
import { I18nextProvider } from 'react-i18next';
import i18next from '@/lib/i18n-client';

export function I18nClientProvider({ locale, children }: { locale: string; children: React.ReactNode }) {
  useEffect(() => {
    if (i18next.language !== locale) {
      i18next.changeLanguage(locale);
    }
  }, [locale]);

  return <I18nextProvider i18n={i18next}>{children}</I18nextProvider>;
}