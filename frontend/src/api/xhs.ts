import { http } from './request'

// 小红书配置接口
export interface XHSConfig {
  id: number
  name: string
  xhs_user_id: string
  is_default: boolean
  is_enabled: boolean
  status: 'pending' | 'active' | 'expired' | 'error'
  last_login_at?: string
  description: string
  sort_order: number
  created_at: string
}

export interface XHSConfigRequest {
  name: string
  cookie?: string
  xhs_user_id?: string
  token?: string
  device_id?: string
  is_default?: boolean
  is_enabled?: boolean
  description?: string
  sort_order?: number
}

// 小红书配置列表
export const getXHSConfigs = (params?: { page?: number; page_size?: number }) => {
  return http.get<{ items: XHSConfig[]; total: number }>('/xhs/configs', { params })
}

// 获取单个小红书配置
export const getXHSConfig = (id: number) => {
  return http.get<XHSConfig>(`/xhs/configs/${id}`)
}

// 创建小红书配置
export const createXHSConfig = (data: XHSConfigRequest) => {
  return http.post<XHSConfig>('/xhs/configs', data)
}

// 更新小红书配置
export const updateXHSConfig = (id: number, data: XHSConfigRequest) => {
  return http.put<XHSConfig>(`/xhs/configs/${id}`, data)
}

// 删除小红书配置
export const deleteXHSConfig = (id: number) => {
  return http.delete<void>(`/xhs/configs/${id}`)
}

// 验证小红书配置
export const verifyXHSConfig = (id: number) => {
  return http.post<{ success: boolean; message: string; user_id?: string }>(`/xhs/configs/${id}/verify`)
}

// 获取用户当前激活的XHS配置
export const getActiveXHSConfig = () => {
  return http.get<XHSConfig>('/xhs/active')
}
