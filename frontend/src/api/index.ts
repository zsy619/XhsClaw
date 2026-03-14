// API接口定义
import { http, type ApiResponse } from './request'

// ==========================================
// 用户相关接口
// ==========================================

// 用户登录
export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  user: UserInfo
}

export interface UserInfo {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  role: string
  status: number
  created_at: string
  updated_at: string
}

export const login = (data: LoginParams) => {
  return http.post<LoginResult>('/auth/login', data)
}

// 用户注册
export interface RegisterParams {
  username: string
  email: string
  password: string
  nickname?: string
}

export const register = (data: RegisterParams) => {
  return http.post<LoginResult>('/auth/register', data)
}

// 获取用户信息
export const getUserInfo = () => {
  return http.get<UserInfo>('/user/info')
}

// 获取用户列表
export interface UserListParams {
  page?: number
  page_size?: number
}

export interface UserListResult {
  list: UserInfo[]
  total: number
  page: number
  page_size: number
}

export const getUserList = (params?: UserListParams) => {
  return http.get<UserListResult>('/users', { params })
}

// ==========================================
// 内容相关接口
// ==========================================

// 生成内容
export interface GenerateContentParams {
  skill_content: string
  count?: number
  length?: string // short | medium | long
}

export interface ContentItem {
  title: string
  description: string
  tags: string[]
}

export interface GenerateContentResult {
  contents: ContentItem[]
}

export const generateContent = (data: GenerateContentParams) => {
  return http.post<GenerateContentResult>('/content/generate', data)
}

// 保存内容
export const saveContent = (data: ContentItem) => {
  return http.post<Content>('/content/save', data)
}

// 内容详情
export interface Content {
  id: number
  user_id: number
  title: string
  description: string
  tags: string
  images: string
  status: number
  publish_time?: string
  created_at: string
  updated_at: string
}

export const getContent = (id: number) => {
  return http.get<Content>(`/content/${id}`)
}

// 内容列表
export interface ContentListParams {
  page?: number
  page_size?: number
  status?: number
}

export interface ContentListResult {
  list: Content[]
  total: number
  page: number
  page_size: number
}

export const getContentList = (params?: ContentListParams) => {
  return http.get<ContentListResult>('/content/list', { params })
}

// 更新内容
export interface UpdateContentParams {
  title?: string
  description?: string
  tags?: string[]
  images?: string[]
  status?: number
  publish_time?: string
}

export const updateContent = (id: number, data: UpdateContentParams) => {
  return http.put<Content>(`/content/${id}`, data)
}

// 删除内容
export const deleteContent = (id: number) => {
  return http.delete(`/content/${id}`)
}

// ==========================================
// 发布相关接口
// ==========================================

// 发布记录
export interface PublishRecord {
  id: number
  user_id: number
  content_id: number
  status: number
  error_msg?: string
  scheduled_at: string
  published_at?: string
  created_at: string
  updated_at: string
  content?: Content
}

// 定时发布
export interface SchedulePublishParams {
  content_id: number
  publish_time: string
  frequency?: string
}

export const schedulePublish = (data: SchedulePublishParams) => {
  return http.post<PublishRecord>('/publish/schedule', data)
}

// 立即发布
export interface PublishNowParams {
  content_id: number
}

export const publishNow = (data: PublishNowParams) => {
  return http.post<PublishRecord>('/publish/now', data)
}

// 发布记录详情
export const getPublishRecord = (id: number) => {
  return http.get<PublishRecord>(`/publish/${id}`)
}

// 发布记录列表
export interface PublishListParams {
  page?: number
  page_size?: number
}

export interface PublishListResult {
  list: PublishRecord[]
  total: number
  page: number
  page_size: number
}

export const getPublishList = (params?: PublishListParams) => {
  return http.get<PublishListResult>('/publish/list', { params })
}

// 取消发布
export const cancelPublish = (id: number) => {
  return http.delete(`/publish/${id}/cancel`)
}
