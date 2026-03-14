import service, { type ApiResponse } from './request'

// 登录请求参数
export interface LoginParams {
  username: string
  password: string
}

// 注册请求参数
export interface RegisterParams {
  username: string
  nickname: string
  email: string
  password: string
}

// 用户信息
export interface UserInfo {
  id: number
  username: string
  email?: string
  nickname?: string
  avatar?: string
  is_active: boolean
  created_at: string
  last_login?: string
}

// 登录响应
export interface LoginResponse {
  access_token: string
  token_type: string
  expires_in: number
  user: UserInfo
}

// 登录（使用 URLSearchParams 格式）
export const login = (params: URLSearchParams) => {
  return service.post<LoginResponse>('/auth/login', params, {
    // 登录接口跳过统一拦截器处理
    skipInterceptor: true
  })
}

// 注册
export const register = (data: RegisterParams) => {
  return service.post('/auth/register', data, {
    // 注册接口跳过统一拦截器处理
    skipInterceptor: true
  })
}

// 获取当前用户信息
export const getCurrentUser = () => {
  return service.get<UserInfo>('/auth/me')
}

// 登出
export const logout = () => {
  return service.post('/auth/logout')
}

// 修改密码
export const changePassword = (data: { old_password: string; new_password: string }) => {
  return service.put('/auth/change-password', data)
}

// 获取用户资料
export const getProfile = () => {
  return service.get('/auth/profile')
}
