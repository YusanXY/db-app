import request from './request'
import type { Tag } from './types'

export interface TagListParams {
  keyword?: string
  sort?: string
  order?: string
  limit?: number
}

export function getTagList(params?: TagListParams) {
  return request.get<Tag[]>('/tags', { params })
}

export function getTagDetail(id: number) {
  return request.get<Tag>(`/tags/${id}`)
}

export function getTagBySlug(slug: string) {
  return request.get<Tag>(`/tags/slug/${slug}`)
}

export function createTag(data: Partial<Tag>) {
  return request.post<Tag>('/tags', data)
}

export function updateTag(id: number, data: Partial<Tag>) {
  return request.put<Tag>(`/tags/${id}`, data)
}

export function deleteTag(id: number) {
  return request.delete(`/tags/${id}`)
}

