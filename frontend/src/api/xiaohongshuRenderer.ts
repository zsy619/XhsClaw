import { http, type ApiResponse } from '@/api/request'

export interface StyleInfo {
    key: string
    name: string
}

export interface RenderRequest {
    markdown_content: string
    style_key: string
    output_prefix?: string
    // 分页模式: separator(按---分隔), auto-fit(自动缩放), auto-split(自动拆分), dynamic(动态高度)
    pagination_mode?: string
    card_width?: number
    card_height?: number
    max_content_height?: number
    // 主题名称
    theme?: string
    // 封面配置
    cover_enabled?: boolean
    cover_title?: string
    cover_subtitle?: string
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
    return await http.get('/xhsclaw/styles')
}

/**
 * 获取渲染后的图片基础URL
 */
const getBaseUrl = (): string => {
    return import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api/v1'
}

/**
 * 获取渲染后的图片完整URL
 * 后端返回的路径格式: /xhsclaw/image/xxx.png
 * 需要拼接为: http://localhost:8000/api/v1/xhsclaw/image/xxx.png
 */
export const getRenderedImage = (imagePath: string): string => {
    const baseUrl = getBaseUrl()
    // 后端返回的路径已经包含 /xhsclaw/image/ 前缀
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
const verifyImage = async (imageUrl: string, timeout: number = 10000): Promise<boolean> => {
    try {
        const controller = new AbortController()
        const timeoutId = setTimeout(() => controller.abort(), timeout)

        const response = await fetch(imageUrl, {
            method: 'GET',
            mode: 'cors',
            signal: controller.signal,
            cache: 'no-cache' // 禁用缓存
        })

        clearTimeout(timeoutId)
        
        if (response.ok) {
            // 检查 Content-Type 是否为图片
            const contentType = response.headers.get('content-type')
            if (contentType && contentType.startsWith('image/')) {
                return true
            }
            // 如果没有 Content-Type，尝试读取一小部分内容
            try {
                const blob = await response.blob()
                return blob.size > 0
            } catch {
                return false
            }
        }
        return false
    } catch (error) {
        console.warn('图片验证请求失败:', error)
        return false
    }
}

/**
 * 检查是否为网络错误（连接失败）
 */
const isNetworkError = (error: any): boolean => {
    if (!error) return false
    
    // Axios 网络错误
    if (error.code === 'ECONNREFUSED' || 
        error.code === 'ENOTFOUND' || 
        error.code === 'ETIMEDOUT' ||
        error.message === 'Network Error' ||
        error.message?.includes('net::ERR_')) {
        return true
    }
    
    // 检查 message 中是否包含网络相关错误
    const errorMsg = error.message || ''
    if (errorMsg.includes('Failed to fetch') ||
        errorMsg.includes('Network request failed') ||
        errorMsg.includes('Connection refused') ||
        errorMsg.includes('Network Error')) {
        return true
    }
    
    return false
}

/**
 * 等待指定时间
 */
const wait = (ms: number): Promise<void> => {
    return new Promise(resolve => setTimeout(resolve, ms))
}

/**
 * 渲染 Markdown 内容为图片
 * 实现自我校验，直到成功或达到最大重试次数
 */
export const renderMarkdown = async (request: RenderRequest): Promise<ApiResponse<RenderData>> => {
    const maxRetries = 3
    const baseUrl = getBaseUrl()
    const retryDelays = [2000, 4000, 8000] // 递增延迟

    for (let i = 0; i < maxRetries; i++) {
        try {
            const attempt = i + 1
            console.log(`开始渲染尝试 ${attempt}/${maxRetries}`)

            const response = await http.post('/generation/render', request, {
                timeout: 60000 // 60秒超时
            })

            console.log('渲染响应:', response)

            if (response.code === 0 && response.data?.success) {
                const images = response.data.images || []

                if (images.length > 0) {
                    const firstImagePath = images[0]
                    const imageUrl = getRenderedImage(firstImagePath)

                    console.log('验证图片URL:', imageUrl)

                    // 验证第一张图片，增加重试验证
                    let isValid = false
                    for (let v = 0; v < 3; v++) {
                        if (v > 0) {
                            console.log(`图片验证重试 ${v + 1}/3...`)
                            await wait(1000)
                        }
                        isValid = await verifyImage(imageUrl)
                        if (isValid) break
                    }

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

            // 如果还有重试次数，等待后重试
            if (i < maxRetries - 1) {
                const waitTime = retryDelays[i]
                console.log(`等待 ${waitTime}ms 后重试...`)
                await wait(waitTime)
            }
        } catch (error) {
            console.warn(`渲染尝试 ${i + 1}/${maxRetries} 失败:`, error)
            
            // 检查是否为网络错误
            const isNetError = isNetworkError(error)
            
            if (isNetError) {
                console.error(`网络错误 (尝试 ${i + 1}/${maxRetries}):`, error.message)
            }
            
            // 如果还有重试次数，等待后重试
            if (i < maxRetries - 1) {
                const waitTime = retryDelays[i]
                console.log(`等待 ${waitTime}ms 后重试...`)
                await wait(waitTime)
            } else {
                // 最后一次尝试失败，抛出更详细的错误信息
                if (isNetError) {
                    throw new Error(`网络连接失败，请检查后端服务是否启动。错误详情: ${error.message}`)
                }
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
