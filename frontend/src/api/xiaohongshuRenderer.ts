import { http, type ApiResponse } from '@/api/request'

export interface StyleInfo {
  key: string
  name: string
}

export interface RenderRequest {
  markdown_content: string
  style_key: string
  output_prefix?: string
  enable_smart_pagination?: boolean
  card_width?: number
  card_height?: number
  max_content_height?: number
}

export interface RenderData {
  success: boolean
  message: string
  images: string[]
  styles?: StyleInfo[]
}

/**
 * 获取可用样式列表
 */
export const getStyles = async (): Promise<ApiResponse<RenderData>> => {
  return await http.get('/xiaohongshu-renderer/styles')
}

/**
 * 渲染 Markdown 内容为图片
 */
export const renderMarkdown = async (request: RenderRequest): Promise<ApiResponse<RenderData>> => {
  return await http.post('/xiaohongshu-renderer/render', request)
}

/**
 * 获取渲染后的图片
 */
export const getRenderedImage = (imagePath: string): string => {
  // 返回图片的完整URL - 支持包含子目录的路径
  // 前端使用 /api/v1 作为 baseURL，所以图片也应该使用 /api/v1 路径
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000'
  // 注意：前端axios默认baseURL是 /api/v1，所以这里也应该使用 /api/v1
  const imageUrl = `${baseUrl}/api/v1/xiaohongshu-renderer/image/${imagePath}`
  console.log('getRenderedImage - 原始路径:', imagePath)
  console.log('getRenderedImage - 完整URL:', imageUrl)
  return imageUrl
}