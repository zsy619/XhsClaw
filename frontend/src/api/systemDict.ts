import { http } from './request'

// 系统字典接口
export interface SystemDict {
  id: number
  category: string
  code: string
  name: string
  value: string
  description: string
  sort_order: number
  enabled: boolean
  extra: string
  created_at: string
}

export interface SystemDictRequest {
  category: string
  code: string
  name: string
  value?: string
  description?: string
  sort_order?: number
  enabled?: boolean
  extra?: string
}

// 获取字典分类列表
export const getDictCategories = () => {
  return http.get<string[]>('/dict/categories')
}

// 获取指定分类的字典项
export const getDictByCategory = (category: string) => {
  return http.get<SystemDict[]>(`/dict/category/${category}`)
}

// 获取所有字典
export const getAllDicts = (params?: { page?: number; page_size?: number }) => {
  return http.get<{ items: SystemDict[]; total: number }>('/dict/all', { params })
}

// 创建字典项
export const createDict = (data: SystemDictRequest) => {
  return http.post<SystemDict>('/dict', data)
}

// 更新字典项
export const updateDict = (id: number, data: SystemDictRequest) => {
  return http.put<SystemDict>(`/dict/${id}`, data)
}

// 删除字典项
export const deleteDict = (id: number) => {
  return http.delete<void>(`/dict/${id}`)
}

// 获取大模型服务商字典
export const getLLMProviderDicts = () => {
  return http.get<SystemDict[]>('/dict/category/llm_provider')
}

// 获取大模型字典
export const getLLMModelDicts = () => {
  return http.get<SystemDict[]>('/dict/category/llm_model')
}
