// src/lib/services/node-service.ts
import { apiFetch } from '@/lib/api/client';
import type { NodeEntity, NodePageData, NodeDetailData } from '@/types/domain';

// ==================== 服务函数 ====================

/**
 * 获取所有论坛板块列表
 * - 成功：返回真实节点数据
 * - 失败：打印错误日志 + 返回空数组兜底，保证侧边栏不白屏
 */
export async function getAllNodes(): Promise<NodePageData> {
  try {
    // apiFetch 自动拆解 envelope.data，直接拿到 NodeEntity[]
    const data = await apiFetch<NodeEntity[]>('/nodes', {
      auth: false,
      cacheStrategy: { next: { tags: ['nodes'], revalidate: 60 } },
    });

    return {
      nodes: Array.isArray(data) ? data : [],
      error: null,
    };
  } catch (err) {
    console.error('[NodeService] Fetch nodes failed:', err);
    return {
      nodes: [],
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}

/**
 * 获取单个板块详情
 * - 用于点击板块时获取名称和简介，供右侧面板展示
 */
export async function getNodeDetail(nodeId: number): Promise<NodeDetailData> {
  try {
    const data = await apiFetch<NodeEntity>(`/node/${nodeId}`, {
      auth: false,
      cacheStrategy: { next: { tags: [`node-${nodeId}`], revalidate: 30 } },
    });

    return {
      node: data ?? null,
      error: null,
    };
  } catch (err) {
    console.error(`[NodeService] Fetch node ${nodeId} failed:`, err);
    return {
      node: null,
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}

/**
 * 标记指定节点为已读
 * POST /nodes/:nodeId/read （JWT 认证接口）
 */
export async function markNodeAsRead(nodeId: number | string): Promise<void> {
  // auth: true → 客户端自动携带 cookie，服务端自动读取 access_token
  // skipDataUnwrap 不需要，后端返回 success:true + data:null，解包后即为 void
  await apiFetch<null>(`/nodes/${nodeId}/read`, {
    method: 'POST',
    auth: true,
  });
}
