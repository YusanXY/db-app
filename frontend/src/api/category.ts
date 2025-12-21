import request from './request'
import type { Category } from './types'

export interface CategoryListParams {
  parent_id?: number
  is_active?: boolean
  tree?: boolean
}

export function getCategoryList(params?: CategoryListParams) {
  return request.get<Category[]>('/categories', { params })
}

export function getCategoryDetail(id: number) {
  return request.get<Category>(`/categories/${id}`)
}

export function getCategoryBySlug(slug: string) {
  return request.get<Category>(`/categories/slug/${slug}`)
}

export function createCategory(data: Partial<Category>) {
  return request.post<Category>('/categories', data)
}

export function updateCategory(id: number, data: Partial<Category>) {
  return request.put<Category>(`/categories/${id}`, data)
}

export function deleteCategory(id: number) {
  return request.delete(`/categories/${id}`)
}

