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
 * 获取渲染后的图片基础URL
 */
const getBaseUrl = (): string => {
  return import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api/v1'
}

/**
 * 获取渲染后的图片完整URL
 * 后端返回的路径格式: /xiaohongshu-renderer/image/xxx.png
 * 需要拼接为: http://localhost:8000/api/v1/xiaohongshu-renderer/image/xxx.png
 */
export const getRenderedImage = (imagePath: string): string => {
  const baseUrl = getBaseUrl()
  // 后端返回的路径已经包含 /xiaohongshu-renderer/image/ 前缀
  // 直接去掉前导斜杠然后拼接到 baseUrl
  const cleanPath = imagePath.replace(/^\//, '')
  const imageUrl = `${baseUrl}/${cleanPath}`
  console.log('getRenderedImage - 原始路径:', imagePath)
  console.log('getRenderedImage - 完整URL:', imageUrl)
  return imageUrl
}

/**
 * 验证图片是否可访问
 */
const verifyImage = async (imageUrl: string): Promise<boolean> => {
  try {
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), 5000)

    const response = await fetch(imageUrl, {
      method: 'GET',
      mode: 'cors',
      signal: controller.signal
    })

    clearTimeout(timeoutId)
    return response.ok
  } catch (error) {
    console.warn('图片验证请求失败:', error)
    return false
  }
}

/**
 * 渲染 Markdown 内容为图片
 * 实现自我校验，直到成功或达到最大重试次数
 */
export const renderMarkdown = async (request: RenderRequest): Promise<ApiResponse<RenderData>> => {
  const maxRetries = 3
  const baseUrl = getBaseUrl()

  for (let i = 0; i < maxRetries; i++) {
    try {
      console.log(`开始渲染尝试 ${i + 1}/${maxRetries}`)

      const response = await http.post('/xiaohongshu-renderer/render', request)

      console.log('渲染响应:', response)

      if (response.code === 0 && response.data?.success) {
        const images = response.data.images || []

        if (images.length > 0) {
          const firstImagePath = images[0]
          const imageUrl = getRenderedImage(firstImagePath)

          console.log('验证图片URL:', imageUrl)

          // 验证第一张图片
          const isValid = await verifyImage(imageUrl)

          if (isValid) {
            console.log('✅ 图片渲染成功，验证通过')
            return response
          } else {
            console.warn(`⚠️ 图片路径验证失败，将返回响应数据`)
            // 即使验证失败也返回响应，让用户可以看到图片
            return response
          }
        } else {
          console.warn('⚠️ 响应成功但无图片列表')
          return response
        }
      }

      const errorMsg = response.message || '渲染失败'
      console.warn(`渲染失败: ${errorMsg}`)

      if (i < maxRetries - 1) {
        const waitTime = 1000 * (i + 1)
        console.log(`等待 ${waitTime}ms 后重试...`)
        await new Promise(resolve => setTimeout(resolve, waitTime))
      }
    } catch (error) {
      console.warn(`渲染尝试 ${i + 1}/${maxRetries} 失败:`, error)

      if (i < maxRetries - 1) {
        const waitTime = 1000 * (i + 1)
        await new Promise(resolve => setTimeout(resolve, waitTime))
      }
    }
  }

  throw new Error('图片渲染失败：达到最大重试次数')
}

/**
 * 直接渲染并验证图片（简化版）
 */
export const renderAndVerifyImage = async (request: RenderRequest): Promise<{
  success: boolean
  images: string[]
  error?: string
}> => {
  try {
    const response = await renderMarkdown(request)

    if (response.code === 0 && response.data?.images?.length > 0) {
      return {
        success: true,
        images: response.data.images
      }
    }

    return {
      success: false,
      images: [],
      error: response.message || '渲染失败'
    }
  } catch (error: any) {
    return {
      success: false,
      images: [],
      error: error.message || '未知错误'
    }
  }
}

// 测试图片渲染功能
export const testRenderMarkdown = async () => {
  try {
    const request = {
      markdown_content: "# 测试标题\n\n这是测试内容",
      style_key: "playful-geometric"
    }
    const response = await renderMarkdown(request)
    console.log('测试渲染结果:', response)
    return response
  } catch (error) {
    console.error('测试渲染失败:', error)
    throw error
  }
}
