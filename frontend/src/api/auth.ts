import request from './request'
import type { LoginRequest, RegisterRequest, LoginResponse, UserResponse } from './types'

export function register(data: RegisterRequest) {
  return request.post<UserResponse>('/auth/register', data)
}

export function login(data: LoginRequest) {
  return request.post<LoginResponse>('/auth/login', data)
}

export function getCurrentUser() {
  return request.get<UserResponse>('/auth/me')
}

