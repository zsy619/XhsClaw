import service, { type ApiResponse } from './request'

// Token使用记录
export interface TokenUsage {
  id: number
  user_id: number
  model: string
  provider: string
  prompt_tokens: number
  completion_tokens: number
  total_tokens: number
  cost: number
  request_type: string
  request_content: string
  response_status: string
  error_message?: string
  created_at: string
}

// Token使用统计
export interface TokenUsageStats {
  total_requests: number
  success_requests: number
  failed_requests: number
  total_prompt_tokens: number
  total_completion_tokens: number
  total_tokens: number
  total_cost: number
  average_tokens: number
}

// 每日使用统计
export interface UserTokenUsage {
  date: string
  total_tokens: number
  total_cost: number
  request_count: number
}

// 按模型统计
export interface TokenUsageByModel {
  model: string
  total_tokens: number
  total_cost: number
  request_count: number
}

// 获取用户Token使用记录
export const getUserTokenUsage = (limit = 50): Promise<ApiResponse<TokenUsage[]>> => {
  return service.get(`/token-usage?limit=${limit}`)
}

// 获取用户Token使用统计
export const getUserTokenStats = (): Promise<ApiResponse<TokenUsageStats>> => {
  return service.get('/token-usage/stats')
}

// 获取用户每日Token使用统计
export const getUserDailyStats = (days = 30): Promise<ApiResponse<UserTokenUsage[]>> => {
  return service.get(`/token-usage/daily?days=${days}`)
}

// 获取用户按模型统计
export const getUserStatsByModel = (): Promise<ApiResponse<TokenUsageByModel[]>> => {
  return service.get('/token-usage/by-model')
}

// 获取全局Token使用统计（管理员）
export const getGlobalTokenStats = (): Promise<ApiResponse<TokenUsageStats>> => {
  return service.get('/admin/token-usage/global')
}

// 获取全局每日Token使用统计（管理员）
export const getGlobalDailyStats = (days = 30): Promise<ApiResponse<UserTokenUsage[]>> => {
  return service.get(`/admin/token-usage/global/daily?days=${days}`)
}
