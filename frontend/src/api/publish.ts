import { http, type ApiResponse } from './request'

export interface PublishRecord {
  id: number
  user_id: number
  content_id: number
  status: number // 0:待发布, 1:发布中, 2:成功, 3:失败
  error_msg?: string
  scheduled_at: string
  published_at?: string
  created_at: string
  updated_at: string
  content?: any
}

export interface SchedulePublishRequest {
  content_id: number
  publish_time: string // RFC3339格式
  frequency?: string // once, daily, weekly
}

export interface PublishNowRequest {
  content_id: number
}

export interface PublishListResponse {
  list: PublishRecord[]
  total: number
  page: number
  page_size: number
}

/**
 * 立即发布内容
 */
export const publishNow = async (data: PublishNowRequest): Promise<ApiResponse<PublishRecord>> => {
  return await http.post('/publish/now', data)
}

/**
 * 定时发布内容
 */
export const schedulePublish = async (data: SchedulePublishRequest): Promise<ApiResponse<PublishRecord>> => {
  return await http.post('/publish/schedule', data)
}

/**
 * 获取发布记录列表
 */
export const getPublishRecords = async (params: {
  page?: number
  page_size?: number
}): Promise<ApiResponse<PublishListResponse>> => {
  return await http.get('/publish/list', { params })
}

/**
 * 获取发布记录详情
 */
export const getPublishRecord = async (id: number): Promise<ApiResponse<PublishRecord>> => {
  return await http.get(`/publish/${id}`)
}

/**
 * 取消发布
 */
export const cancelPublish = async (id: number): Promise<ApiResponse<any>> => {
  return await http.post(`/publish/${id}/cancel`)
}

/**
 * 重试发布
 */
export const retryPublish = async (id: number): Promise<ApiResponse<PublishRecord>> => {
  return await http.post(`/publish/${id}/retry`)
}
