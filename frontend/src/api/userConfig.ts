import service from './request'

// 用户配置请求参数
export interface UserConfigRequest {
  llm_api_key?: string
  llm_base_url?: string
  llm_model?: string
  xiaohongshu_cookie?: string
  xiaohongshu_user_id?: string
  xiaohongshu_token?: string
  default_publish_time?: string
  auto_publish_enabled?: boolean
}

// 用户配置响应
export interface UserConfigResponse {
  id: number
  user_id: number
  llm_api_key: string
  llm_base_url: string
  llm_model: string
  xiaohongshu_cookie: string
  xiaohongshu_user_id: string
  xiaohongshu_token: string
  default_publish_time: string
  auto_publish_enabled: boolean
  created_at: string
  updated_at: string
}

// 获取用户配置
export const getUserConfig = () => {
  return service.get<UserConfigResponse>('/user/config')
}

// 更新用户配置
export const updateUserConfig = (data: UserConfigRequest) => {
  return service.put<UserConfigResponse>('/user/config', data)
}
