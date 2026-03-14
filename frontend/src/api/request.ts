import router from '@/router'
import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'

// 响应接口
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 扩展 AxiosRequestConfig 类型，添加重试计数
interface AxiosRequestConfigWithRetry extends InternalAxiosRequestConfig {
  __retryCount?: number
}

// 创建 axios 实例
const service: AxiosInstance = axios.create({
  baseURL: 'http://localhost:8000/api/v1',
  timeout: 30000,
  // 不设置默认 Content-Type，让 axios 根据数据类型自动设置
})

// 请求拦截器
service.interceptors.request.use(
  (config: AxiosRequestConfigWithRetry) => {
    const isAuthEndpoint = config.url?.includes('/auth/login') || config.url?.includes('/auth/register')
    
    // 登录和注册接口不添加 token
    if (!isAuthEndpoint) {
      const token = localStorage.getItem('token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
    }
    
    // 添加请求时间戳，防止缓存
    if (config.method === 'get') {
      config.params = {
        ...config.params,
        _t: Date.now(),
      }
    }
    
    return config
  },
  (error) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const res = response.data
    
    // 检查是否有 skipInterceptor 标记（用于登录等特殊接口）
    const config = response.config as any
    if (config.skipInterceptor) {
      // 直接返回数据，不进行 code 检查
      return res
    }
    
    // 如果返回的状态码不是 0，说明接口有错误
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      
      // 10002: 未授权，跳转到登录页
      if (res.code === 10002) {
        localStorage.removeItem('token')
        router.push('/login')
      }
      
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    
    return res
  },
  (error) => {
    console.error('响应错误:', error)
    
    const config = error.config as any
    const isAuthEndpoint = config?.url?.includes('/auth/login') || config?.url?.includes('/auth/register')
    
    // 登录和注册接口的 401 错误不需要特殊处理，直接返回错误给调用方
    if (isAuthEndpoint && error.response?.status === 401) {
      return Promise.reject(error)
    }
    
    let message = '网络错误，请稍后重试'
    
    if (error.response) {
      switch (error.response.status) {
        case 400:
          message = '请求参数错误'
          break
        case 401:
          message = '未授权，请重新登录'
          localStorage.removeItem('token')
          router.push('/login')
          break
        case 403:
          message = '拒绝访问'
          break
        case 404:
          message = '请求资源不存在'
          break
        case 500:
          message = '服务器内部错误'
          break
        case 502:
          message = '网关错误'
          break
        case 503:
          message = '服务不可用'
          break
        case 504:
          message = '网关超时'
          break
        default:
          message = `连接错误${error.response.status}`
      }
    } else if (error.request) {
      // 超时重试逻辑
      const config = error.config as AxiosRequestConfigWithRetry
      if (config && !config.__retryCount) {
        config.__retryCount = config.__retryCount || 0
        
        // 最多重试 2 次
        if (config.__retryCount < 2) {
          config.__retryCount++
          console.log(`请求重试，第${config.__retryCount}次`)
          
          // 延迟 500ms 后重试
          return new Promise((resolve) => {
            setTimeout(() => {
              resolve(service(config))
            }, 500)
          })
        }
      }
      
      message = '无法连接到服务器，请检查网络连接'
    }
    
    // 登录和注册接口不显示全局错误消息，让调用方自己处理
    if (!isAuthEndpoint) {
      ElMessage.error(message)
    }
    
    return Promise.reject(error)
  }
)

// 导出请求方法
export const http = {
  get<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return service.get(url, config)
  },

  post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return service.post(url, data, config)
  },

  put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return service.put(url, data, config)
  },

  delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return service.delete(url, config)
  },

  upload<T = any>(url: string, data: FormData, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return service.post(url, data, {
      ...config,
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
  },
}

export default service
