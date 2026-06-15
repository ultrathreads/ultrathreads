'use server';

import { revalidatePath } from 'next/cache';
import { apiFetch, ApiError, ApiConnectionError } from '@/lib/api/client';

export interface NodePayload {
  name: string;
  description?: string;
  icon?: string;
  sortNo: number;
  status: number;
}

export interface PageInfo {
  page: number;
  pageSize: number;
  total: number;
}

export interface NodeItem {
  id: number;
  name: string;
  description: string;
  icon: string;
  sortNo: number;
  status: number;
  topicCount: number;
  createTime: number;
}

// 统一异常处理
function handleApiError(err: unknown): { success: false; message: string } {
  if (err instanceof ApiError) return { success: false, message: err.message };
  if (err instanceof ApiConnectionError) return { success: false, message: err.message };
  return { success: false, message: '未知错误，请稍后重试' };
}

const REVALIDATE_PATH = '/(admin)/admin/nodes';

/**
 * 获取节点列表
 * 后端响应格式: { page: {...}, results: [...] }
 */
export async function getNodes(page = 1, pageSize = 20) {
  try {
    const envelope = await apiFetch<{ page: PageInfo; results: NodeItem[] }>(
      `/admin/nodes?page=${page}&pageSize=${pageSize}`,
      { auth: true, noCache: true, skipDataUnwrap: true }
    );

    // 从完整响应中提取 results 数组
    const nodes = Array.isArray(envelope.results) ? envelope.results : [];
    const pageInfo = envelope.page ?? { page: 1, pageSize: 20, total: 0 };

    return { nodes, pageInfo };
  } catch (err) {
    console.error('[getNodes] Error:', err);
    return { 
      nodes: [], 
      pageInfo: { page: 1, pageSize: 20, total: 0 } 
    };
  }
}

export async function getNode(id: number) {
  try {
    // 单个节点接口可能返回 { data: {...} } 或直接返回对象
    // 先用 skipDataUnwrap 安全获取
    const envelope = await apiFetch<any>(`/admin/nodes/${id}`, {
      auth: true,
      noCache: true,
      skipDataUnwrap: true,
    });
    
    if (!envelope || typeof envelope !== 'object' || !('id' in envelope)) {
      throw new Error(`获取节点 #${id} 失败：响应格式异常`);
    }

    // 兼容多种可能的单条数据响应格式
    return envelope;
  } catch (err) {
    throw err;
  }
}

export async function createNode(payload: NodePayload) {
  try {
    const data = await apiFetch<any>('/admin/nodes', {
      method: 'POST',
      body: JSON.stringify(payload),
      auth: true,
    });
    revalidatePath(REVALIDATE_PATH);
    return { success: true as const, data };
  } catch (err) {
    return handleApiError(err);
  }
}

export async function updateNode(id: number, payload: NodePayload) {
  try {
    const data = await apiFetch<any>(`/admin/nodes/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
      auth: true,
    });
    revalidatePath(REVALIDATE_PATH);
    return { success: true as const, data };
  } catch (err) {
    return handleApiError(err);
  }
}

export async function updateNodeSort(items: { id: number; sortNo: number }[]) {
  return apiFetch<{ updated: number }>('/admin/nodes/sort', {
    method: 'PUT',
    auth: true,
    noCache: true,
    body: JSON.stringify({ items }),
  });
}


export async function deleteNode(id: number) {
  try {
    await apiFetch(`/admin/nodes/${id}`, { method: 'DELETE', auth: true });
    revalidatePath(REVALIDATE_PATH);
    return { success: true as const };
  } catch (err) {
    return handleApiError(err);
  }
}