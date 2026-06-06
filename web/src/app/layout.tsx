// src/app/layout.tsx
import './globals.css';
import { headers } from 'next/headers';
import { I18nClientProvider } from '@/components/I18nClientProvider';
import { getServerTranslation } from '@/lib/i18n-server';
import type { Metadata } from 'next';

export async function generateMetadata(): Promise<Metadata> {
  const headersList = await headers();
  const locale = headersList.get('x-locale') || 'zh';

  const t = await getServerTranslation(locale, ['common']);

  return {
    title: {
      default: t('common:title'),       // ← 从 i18n 文件动态获取
      template: `%s | ${t('common:title')}`,
    },
  };
}

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const headersList = await headers();
  const locale = headersList.get('x-locale') || 'zh';

  return (
    <html lang={locale}>
      <body>
        <I18nClientProvider locale={locale}>{children}</I18nClientProvider>
      </body>
    </html>
  );
}