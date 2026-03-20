import { http } from './request'

// 大模型服务商配置接口
export interface LLMProvider {
  id: number
  name: string
  provider: string
  api_key?: string
  base_url: string
  model_name: string
  is_default: boolean
  is_enabled: boolean
  timeout: number
  retry_count: number
  description: string
  sort_order: number
  created_at: string
}

export interface LLMProviderRequest {
  name: string
  provider: string
  api_key?: string
  base_url?: string
  model_name?: string
  is_default?: boolean
  is_enabled?: boolean
  timeout?: number
  retry_count?: number
  description?: string
  sort_order?: number
}

// 大模型服务商列表
export const getLLMProviders = (params?: { page?: number; page_size?: number }) => {
  return http.get<{ items: LLMProvider[]; total: number }>('/llm/providers', { params })
}

// 获取单个服务商配置
export const getLLMProvider = (id: number) => {
  return http.get<LLMProvider>(`/llm/providers/${id}`)
}

// 创建服务商配置
export const createLLMProvider = (data: LLMProviderRequest) => {
  return http.post<LLMProvider>('/llm/providers', data)
}

// 更新服务商配置
export const updateLLMProvider = (id: number, data: LLMProviderRequest) => {
  return http.put<LLMProvider>(`/llm/providers/${id}`, data)
}

// 删除服务商配置
export const deleteLLMProvider = (id: number) => {
  return http.delete<void>(`/llm/providers/${id}`)
}

// 获取用户当前激活的LLM配置
export const getActiveLLMConfig = () => {
  return http.get<LLMProvider>('/llm/active')
}

// 测试LLM连接
export const testLLMConnection = (id: number) => {
  return http.post<{ success: boolean; message: string }>(`/llm/providers/${id}/test`)
}
