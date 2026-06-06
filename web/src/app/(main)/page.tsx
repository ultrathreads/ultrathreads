import type { Metadata } from 'next';
import { getServerTranslation } from '@/lib/i18n-server';
import ThreadTree from '@/components/ThreadTree';
import Pagination from '@/components/Pagination';
import { getMockPageData } from '@/lib/mock-data';

// 👇 注意：searchParams 现在是 Promise 类型
interface Props {
  searchParams: Promise<{ page?: string }>;
}

// 👇 动态生成 metadata，支持 await searchParams
export async function generateMetadata({ searchParams }: Props): Promise<Metadata> {
  const { page: pageStr } = await searchParams;
  const page = Math.max(1, parseInt(pageStr || '1', 10));

  return {
    title: page > 1 ? `帖子列表 - 第 ${page} 页` : '帖子列表',
  };
}

export default async function HomePage({ searchParams }: Props) {
  const t = await getServerTranslation(['common','home']);
  // 👇 必须先 await 解包
  const { page: pageStr } = await searchParams;
  const page = Math.max(1, parseInt(pageStr || '1', 10));

  const data = getMockPageData(page);

  return (
    <>
      <ThreadTree threads={data.threads} />
      <Pagination
        totalItems={data.totalItems}
        pageSize={data.pageSize}
        currentPage={data.currentPage}
      />
    </>
  );
}