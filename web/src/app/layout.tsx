// src/app/layout.tsx
import { headers } from 'next/headers';
import { I18nClientProvider } from '@/providers/I18nClientProvider';
import { SiteConfigProvider } from '@/providers/SiteConfigProvider'; 
import { AuthProvider } from '@/providers/AuthProvider';
import { getServerSession } from '@/lib/auth/session';
import type { Metadata } from 'next';
import { fetchSiteConfig } from '@/services/site-service';

export async function generateMetadata(): Promise<Metadata> {
  const config = await fetchSiteConfig();

  return {
    title: {
      default: config.siteTitle,
      // template 只放站点名，%s 会被子页面的 title 替换
      template: `%s | ${config.siteTitle}`,
    },
    // description 单独设置
    description: config.siteDescription,
  };
}

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const headersList = await headers();
  const locale = headersList.get('x-locale') || 'zh';
  const initialUser = await getServerSession();
  const config = await fetchSiteConfig(); 

  return (
    <html lang={locale} suppressHydrationWarning>
      <body>
        <SiteConfigProvider config={config}>
          <I18nClientProvider locale={locale}>
             <AuthProvider initialUser={initialUser}>
              {children}
            </AuthProvider>
          </I18nClientProvider>
        </SiteConfigProvider>
      </body>
    </html>
  );
}