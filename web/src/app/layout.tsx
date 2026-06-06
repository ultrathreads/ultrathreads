// src/app/layout.tsx
import { headers } from 'next/headers';
import { I18nClientProvider } from '@/components/I18nClientProvider';
import { getServerTranslation } from '@/lib/i18n-server';
import type { Metadata } from 'next';
import { getSiteConfig } from '@/lib/api';

import './globals.css';


export async function generateMetadata(): Promise<Metadata> {
  const config = await getSiteConfig();

  return {
    title: {
      default: config.siteTitle,
      template: `%s | ${config.siteTitle}`,
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