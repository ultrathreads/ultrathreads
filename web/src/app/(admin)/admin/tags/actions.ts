// src/app/(admin)/admin/tags/actions.ts
import { apiFetch } from '@/lib/api/client';

export async function getTags(status?: number) {
  const params = status !== undefined ? `?status=${status}` : '';
  return apiFetch<Tag[]>(`/admin/tags${params}`, { auth: true });
}

export async function createTag(data: Partial<Tag>) {
  return apiFetch<Tag>('/admin/tags', {
    method: 'POST', auth: true, noCache: true,
    body: JSON.stringify(data),
  });
}

export async function updateTag(id: number, data: Partial<Tag>) {
  return apiFetch<{ updated: boolean }>(`/admin/tags/${id}`, {
    method: 'PUT', auth: true, noCache: true,
    body: JSON.stringify(data),
  });
}

export async function deleteTag(id: number) {
  return apiFetch<{ deleted: boolean }>(`/admin/tags/${id}`, {
    method: 'DELETE', auth: true, noCache: true,
  });
}

export async function updateTagSort(items: { id: number; sortNo: number }[]) {
  return apiFetch<{ updated: number }>('/admin/tags/sort', {
    method: 'PUT', auth: true, noCache: true,
    body: JSON.stringify({ items }),
  });
}