export interface User {
  id: number
  username: string
  nickname: string
  avatar_url: string
  bio: string
  role: string
  created_at: string
}

export interface Article {
  id: number
  title: string
  slug: string
  content?: string
  content_html?: string
  summary: string
  cover_image_url: string
  author: User
  editor?: User
  categories: Category[]
  tags: Tag[]
  view_count: number
  like_count: number
  comment_count: number
  is_featured: boolean
  is_liked?: boolean
  status: string
  published_at?: string
  created_at: string
  updated_at: string
}

export interface Category {
  id: number
  name: string
  slug: string
  description: string
}

export interface Tag {
  id: number
  name: string
  slug: string
}

export interface Pagination {
  page: number
  page_size: number
  total: number
  total_pages: number
}

export interface ArticleListParams {
  page?: number
  page_size?: number
  category_id?: number
  tag_id?: number
  author_id?: number
  status?: string
  keyword?: string
  sort?: string
  order?: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  nickname?: string
}

export interface LoginResponse {
  token: string
  refresh_token?: string
  expires_in: number
  user: UserResponse
}

export interface UserResponse {
  id: number
  username: string
  nickname: string
  avatar_url: string
  bio: string
  role: string
  created_at: string
}

