// src/app/layout.tsx
import { headers } from 'next/headers';
import { I18nClientProvider } from '@/components/providers/I18nClientProvider';
import { AuthProvider } from '@/components/providers/AuthProvider';
import type { Metadata } from 'next';
import { fetchSiteConfig } from '@/services/site-service';

import './globals.css';

export async function generateMetadata(): Promise<Metadata> {
  const config = await fetchSiteConfig();

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
        {/* ✅ 保持 I18nClientProvider 在最外层，确保 AuthProvider 及其子组件能使用多语言 */}
        <I18nClientProvider locale={locale}>
          <AuthProvider>
            {children}
          </AuthProvider>
        </I18nClientProvider>
      </body>
    </html>
  );
}