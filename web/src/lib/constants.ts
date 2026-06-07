// src/lib/constants.ts

/**
 * 1. UI/UX 全局配置（避免组件中出现魔法数字）
 */
export const PAGINATION = {
  DEFAULT_PAGE_SIZE: 20,
  MAX_PAGE_SIZE: 100,
} as const;

export const TOAST_DURATION = 3000; // ms
export const DEBOUNCE_DELAY = 300;  // ms

/**
 * 2. 客户端路由路径（用于 Link、router.push、redirect）
 */
export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  NODE_DETAIL: (id: string) => `/nodes/${id}`,
  THREAD_DETAIL: (id: string) => `/threads/${id}`,
} as const;

/**
 * 3. 本地存储 Key（统一管理，防止拼写错误和冲突）
 */
export const STORAGE_KEYS = {
  AUTH_TOKEN: 'auth_token',
  THEME_PREFERENCE: 'theme_preference',
  I18N_LOCALE: 'i18n_locale',
} as const;

/**
 * 4. 纯前端展示用的枚举映射（Go 返回 code/status，前端负责翻译为 UI 文案）
 */
export const THREAD_STATUS_MAP = {
  OPEN: { label: '进行中', color: 'green' },
  CLOSED: { label: '已关闭', color: 'gray' },
  PINNED: { label: '置顶', color: 'blue' },
} as const;

/**
 * 5. 第三方服务公开配置（注意：绝不能放 Secret）
 */
export const OAUTH_PROVIDERS = ['github', 'google'] as const;