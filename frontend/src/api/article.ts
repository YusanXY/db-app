import request from './request'
import type { Article, ArticleListParams, Pagination } from './types'

export function getArticleList(params: ArticleListParams) {
  return request.get<{ items: Article[]; pagination: Pagination }>('/articles', { params })
}

export function getArticleDetail(id: number) {
  return request.get<Article>(`/articles/${id}`)
}

export function createArticle(data: Partial<Article>) {
  return request.post<Article>('/articles', data)
}

export function updateArticle(id: number, data: Partial<Article>) {
  return request.put<Article>(`/articles/${id}`, data)
}

export function deleteArticle(id: number) {
  return request.delete(`/articles/${id}`)
}

