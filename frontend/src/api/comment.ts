import request from './request'
import type { Comment, Pagination } from './types'

export interface CommentListParams {
  page?: number
  page_size?: number
}

export interface CommentListResponse {
  items: Comment[]
  pagination: Pagination
}

export function getCommentList(articleId: number, params?: CommentListParams) {
  return request.get<CommentListResponse>(`/articles/${articleId}/comments`, { params })
}

export function createComment(articleId: number, data: { content: string; parent_id?: number }) {
  return request.post<Comment>(`/articles/${articleId}/comments`, data)
}

export function updateComment(id: number, data: { content: string }) {
  return request.put<Comment>(`/comments/${id}`, data)
}

export function deleteComment(id: number) {
  return request.delete(`/comments/${id}`)
}

export function toggleCommentLike(id: number) {
  return request.post<{ is_liked: boolean }>(`/comments/${id}/like`)
}

