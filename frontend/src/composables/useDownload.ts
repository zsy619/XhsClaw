/**
 * 下载进度 Composable
 * 提供文件下载进度跟踪和提示
 */
import { ref, computed, type Ref } from 'vue'
import { toast } from './useToast'

/**
 * 下载进度信息
 */
interface DownloadProgress {
  /** 已下载的字节数 */
  loaded: number
  /** 总字节数 */
  total: number
  /** 下载进度百分比 (0-100) */
  percentage: number
  /** 下载速度 (字节/秒) */
  speed: number
  /** 剩余时间 (秒) */
  estimatedTime: number
  /** 是否完成 */
  isCompleted: boolean
}

/**
 * 下载配置
 */
interface DownloadOptions {
  /** 文件名 */
  filename?: string
  /** 是否显示进度提示，默认 true */
  showProgress?: boolean
  /** 成功提示消息 */
  successMessage?: string
  /** 错误提示消息 */
  errorMessage?: string
  /** 下载完成回调 */
  onComplete?: () => void
  /** 下载失败回调 */
  onError?: (error: any) => void
  /** 进度更新回调 */
  onProgress?: (progress: DownloadProgress) => void
}

/**
 * useDownload Hook
 * 管理文件下载进度
 * 
 * @example
 * ```ts
 * import { useDownload } from '@/composables/useDownload'
 * 
 * const { download, progress, isDownloading } = useDownload()
 * 
 * // 下载文件
 * await download('/api/export', {
 *   filename: 'data.xlsx',
 *   showProgress: true,
 *   successMessage: '下载完成',
 *   onProgress: (p) => {
 *     console.log(`下载进度：${p.percentage.toFixed(1)}%`)
 *   }
 * })
 * ```
 */
export function useDownload() {
  const isDownloading = ref(false)
  const progress = ref<DownloadProgress>({
    loaded: 0,
    total: 0,
    percentage: 0,
    speed: 0,
    estimatedTime: 0,
    isCompleted: false
  })

  const startTime = ref<number>(0)
  let lastLoaded = 0
  let lastTime = 0

  /**
   * 下载文件
   */
  async function download(
    url: string,
    options: DownloadOptions = {}
  ): Promise<Blob | null> {
    const {
      filename,
      showProgress = true,
      successMessage = '下载完成',
      errorMessage,
      onComplete,
      onError,
      onProgress
    } = options

    isDownloading.value = true
    progress.value = {
      loaded: 0,
      total: 0,
      percentage: 0,
      speed: 0,
      estimatedTime: 0,
      isCompleted: false
    }

    startTime.value = Date.now()
    lastLoaded = 0
    lastTime = startTime.value

    try {
      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })

      if (!response.ok) {
        throw new Error('下载失败')
      }

      // 获取总大小
      const contentLength = response.headers.get('content-length')
      const total = contentLength ? parseInt(contentLength, 10) : 0

      progress.value.total = total

      // 读取响应数据
      const reader = response.body?.getReader()
      if (!reader) {
        throw new Error('无法读取响应流')
      }

      const chunks: Uint8Array[] = []
      
      while (true) {
        const { done, value } = await reader.read()
        
        if (done) break
        
        chunks.push(value)
        
        // 更新进度
        lastLoaded += value.length
        const currentTime = Date.now()
        const timeDiff = (currentTime - lastTime) / 1000 || 1
        
        progress.value.loaded += value.length
        progress.value.percentage = total > 0 
          ? (progress.value.loaded / total) * 100 
          : 0
        progress.value.speed = lastLoaded / timeDiff
        progress.value.estimatedTime = total > 0
          ? (total - progress.value.loaded) / progress.value.speed
          : 0

        // 触发进度回调
        if (showProgress && onProgress) {
          onProgress(progress.value)
        }

        lastLoaded = 0
        lastTime = currentTime
      }

      // 合并数据块
      const blob = new Blob(chunks)
      
      // 标记完成
      progress.value.isCompleted = true
      progress.value.percentage = 100

      // 自动下载
      if (filename) {
        const downloadUrl = URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = downloadUrl
        link.download = filename
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        URL.revokeObjectURL(downloadUrl)
      }

      // 成功提示
      if (showProgress && successMessage) {
        toast.success(successMessage)
      }

      // 完成回调
      onComplete?.()

      return blob
    } catch (error: any) {
      // 错误处理
      const message = errorMessage || error?.message || '下载失败'
      if (showProgress) {
        toast.error(message)
      }
      
      onError?.(error)
      return null
    } finally {
      isDownloading.value = false
    }
  }

  /**
   * 重置进度
   */
  function reset() {
    isDownloading.value = false
    progress.value = {
      loaded: 0,
      total: 0,
      percentage: 0,
      speed: 0,
      estimatedTime: 0,
      isCompleted: false
    }
  }

  return {
    isDownloading,
    progress,
    download,
    reset
  }
}

export default useDownload
