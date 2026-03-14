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
  role_id: number
  role?: Role
  status: number
  created_at: string
  updated_at: string
}

// 兼容旧格式，导出为 User
export type User = UserInfo

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
  title_options: string // JSON格式的备选标题数组
  selected_title_index: number // 选中的备选标题索引
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
// 内容历史记录相关接口
// ==========================================

// 历史记录类型
export interface ContentHistory {
  id: number
  content_id: number
  user_id: number
  type: string // 'create' | 'edit' | 'delete' | 'publish'
  title: string
  description: string
  tags: string
  title_options: string
  selected_title_index: number
  change_reason?: string
  created_at: string
}

// 历史记录列表参数
export interface HistoryListParams {
  content_id?: number
  page?: number
  page_size?: number
}

// 历史记录列表结果
export interface HistoryListResult {
  list: ContentHistory[]
  total: number
  page: number
  page_size: number
}

// 获取历史记录列表
export const getHistoryList = (params?: HistoryListParams) => {
  return http.get<HistoryListResult>('/content/histories/list', { params })
}

// 获取历史记录详情
export const getHistoryDetail = (id: number) => {
  return http.get<ContentHistory>(`/content/histories/${id}`)
}

// 恢复到历史版本
export const restoreHistory = (id: number) => {
  return http.post<Content>(`/content/histories/${id}/restore`)
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

// ==========================================
// 角色和权限相关接口
// ==========================================

// 角色信息
export interface Role {
  id: number
  name: string
  code: string
  description: string
  permissions: string // JSON格式的权限数组
  is_system: boolean
  created_at: string
  updated_at: string
}

// 权限信息
export interface Permission {
  id: number
  name: string
  code: string
  module: string
  description: string
  created_at: string
  updated_at: string
}

// 角色列表参数
export interface RoleListParams {
  page?: number
  page_size?: number
}

// 角色列表结果
export interface RoleListResult {
  list: Role[]
  total: number
  page: number
  page_size: number
}

// 获取角色列表
export const getRoleList = (params?: RoleListParams) => {
  return http.get<RoleListResult>('/roles', { params })
}

// 获取所有角色
export const getAllRoles = () => {
  return http.get<Role[]>('/roles/all')
}

// 获取角色详情
export const getRole = (id: number) => {
  return http.get<Role>(`/roles/${id}`)
}

// 创建角色参数
export interface CreateRoleParams {
  name: string
  code: string
  description?: string
  permissions?: string[]
}

// 创建角色
export const createRole = (data: CreateRoleParams) => {
  return http.post<Role>('/roles', data)
}

// 更新角色参数
export interface UpdateRoleParams {
  name?: string
  description?: string
  permissions?: string[]
}

// 更新角色
export const updateRole = (id: number, data: UpdateRoleParams) => {
  return http.put<Role>(`/roles/${id}`, data)
}

// 删除角色
export const deleteRole = (id: number) => {
  return http.delete(`/roles/${id}`)
}

// 获取所有权限
export const getPermissions = () => {
  return http.get<Permission[]>('/permissions')
}

// 更新用户角色参数
export interface UpdateUserRoleParams {
  role_id: number
}

// 更新用户角色
export const updateUserRole = (id: number, data: UpdateUserRoleParams) => {
  return http.put<User>(`/users/${id}/role`, data)
}

// 更新用户状态参数
export interface UpdateUserStatusParams {
  status: number
}

// 更新用户状态
export const updateUserStatus = (id: number, data: UpdateUserStatusParams) => {
  return http.put<User>(`/users/${id}/status`, data)
}
