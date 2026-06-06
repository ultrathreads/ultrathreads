import { Thread, ForumBoard, Tag, PageData } from '@/types';

const BOARDS: ForumBoard[] = [
  { name: '技术交流', icon: '💻', count: 128 },
  { name: '设计创意', icon: '🎨', count: 86 },
  { name: '产品运营', icon: '📱', count: 64 },
  { name: 'AI & 大模型', icon: '🤖', count: 215 },
  { name: '灌水闲聊', icon: '☕', count: 342 },
];

const TAGS: Tag[] = [
  { label: 'Vue3' }, { label: 'React' }, { label: 'TypeScript' },
  { label: 'Node.js' }, { label: 'CSS技巧' }, { label: '性能优化' },
];

function generateReplies(depth: number, maxDepth: number): Thread['replies'] {
  if (depth >= maxDepth) return [];
  const count = Math.floor(Math.random() * 3);
  return Array.from({ length: count }, (_, i) => ({
    id: `reply-${depth}-${i}-${Date.now()}`,
    title: `Re: 回复内容 ${depth}-${i}`,
    author: ['王芳', '赵强', '孙丽', '刘洋', '黄磊'][Math.floor(Math.random() * 5)],
    date: '2026-06-03, 04:30 PM',
    isRead: Math.random() > 0.5,
    replies: generateReplies(depth + 1, maxDepth),
  }));
}

export function getMockPageData(page: number, pageSize: number = 10): PageData {
  const totalItems = 86;
  const threads: Thread[] = Array.from({ length: pageSize }, (_, i) => {
    const idx = (page - 1) * pageSize + i;
    if (idx >= totalItems) return null;
    return {
      id: `${idx}`,
      title: [
        'Vue3 Composition API 最佳实践讨论',
        'Docker Compose 多服务编排实战经验',
        'Rust async runtime 选型对比（Tokio vs async-std）',
        'Test Thread (Sandbox)',
        'Release of version 2.4.99.1',
        'Next.js App Router 深度解析',
        'CSS Container Queries 实战',
        'WebAssembly 在前端的应用场景',
        '微前端架构落地踩坑记录',
        'GraphQL vs REST API 选型指南',
      ][idx % 10],
      author: ['李明', '陈伟', '吴昊', 'Alfie', 'cool', 'Tmmy T.', 'Rich', '周杰'][idx % 8],
      date: '2026-06-01, 02:20 PM',
      category: '技术交流',
      isRead: idx % 3 === 0,
      replies: generateReplies(0, 3),
    };
  }).filter(Boolean) as Thread[];

  return { threads, totalItems, currentPage: page, pageSize, boards: BOARDS, tags: TAGS };
}

// src/lib/mock-data.ts

// ... 保持原有的 Thread, PageData 类型和 getMockPageData 函数不变 ...

// 1. 新增帖子详情的类型定义
export interface PostDetailData {
  id: string;
  title: string;
  author: string;
  authorAvatar: string;
  date: string;
  tag: string;
  views: number;
  comments: number;
  content: string;
  replies: Thread[]; // 👈 改为 Thread[]
  likes: number;
  favorites: number;
}

// 2. 新增一个模拟的帖子详情数据
const mockPostDetail: PostDetailData = {
  id: '3021',
  title: 'Vue3 Composition API 最佳实践讨论',
  author: '李明',
  authorAvatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=LiMing',
  date: '2026-06-01, 02:20 PM',
  tag: '技术交流',
  views: 1286,
  comments: 8,
  likes: 42,
  favorites: 18,
  content: `
    <p>最近在团队项目中全面迁移到 Vue3 Composition API，总结了一些实战中的最佳实践和踩坑经验，分享给大家讨论。</p>
    <p>首先关于 <code>setup()</code> 语法糖的使用，推荐在所有新组件中统一采用 <code>&lt;script setup&gt;</code> 写法，不仅减少样板代码，还能获得更好的 TypeScript 类型推导支持。需要注意的是，顶层 await 仅在 SFC 的 script setup 中可用，普通 setup 函数不支持。</p>
    <p>其次关于响应式数据的选择：<code>ref</code> 适合基本类型和需要整体替换的对象，<code>reactive</code> 适合深层嵌套的配置对象。避免对 reactive 对象进行解构赋值，这会导致响应性丢失，应使用 <code>toRefs()</code> 代替。</p>
    <p>最后是关于组合式函数（Composables）的组织方式，建议按功能域拆分文件，每个 composable 只关注单一职责，并通过参数注入依赖而非直接 import 全局状态，这样更利于单元测试和复用。</p>
  `,
  replies: [
    {
      id: 'r1',
      title: 'Re: Vue3 Composition API 最佳实践讨论', // 补上 title 字段
      author: '前端小白',
      date: '2026-06-01, 03:00 PM',
      category: '技术交流',
      isRead: false,
      replies: [],
    },
  ],
};

// 3. 新增获取帖子详情的函数
export function getMockPostDetail(id: string): PostDetailData | null {
  // 在实际项目中，这里会是 API 调用
  if (id === mockPostDetail.id) {
    return mockPostDetail;
  }
  return null;
}

// src/lib/mock-data.ts

// ... 保持原有的类型定义和 getMockPageData 函数不变 ...

// 1. 新增 Board 类型定义
export interface Board {
  name: string;
  icon: string;
  count: number;
}

// 2. 导出 boards 模拟数据
export const mockBoards: Board[] = [
  { name: '技术交流', icon: '💻', count: 128 },
  { name: '设计创意', icon: '🎨', count: 86 },
  { name: '产品运营', icon: '📱', count: 64 },
  { name: 'AI & 大模型', icon: '🤖', count: 215 },
  { name: '灌水闲聊', icon: '☕', count: 342 },
];

// ... 保持原有的 PostDetailData 类型和 getMockPostDetail 函数不变 ...
// src/lib/mock-data.ts

// ... 保持原有的类型定义、getMockPageData、mockBoards 等不变 ...

// 1. 新增 Tag 类型定义
export interface Tag {
  label: string;
}

// 2. 导出 tags 模拟数据
export const mockTags: Tag[] = [
  { label: 'React' },
  { label: 'Vue3' },
  { label: 'Next.js' },
  { label: 'TypeScript' },
  { label: 'Node.js' },
  { label: 'AI大模型' },
  { label: '性能优化' },
  { label: '前端架构' },
];

// ... 保持原有的 PostDetailData 类型和 getMockPostDetail 函数不变 ...
