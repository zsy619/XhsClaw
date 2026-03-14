/**
 * 通用异步操作 Composable
 * 提供统一的异步操作处理，包括加载状态、错误处理等
 */
import { ref, type Ref } from 'vue'
import { toast } from './useToast'

/**
 * 异步操作配置
 */
interface AsyncActionOptions<T> {
  /** 成功时的提示消息，设为 false 则不显示 */
  successMessage?: string | false
  /** 错误时的提示消息，设为 false 则不显示 */
  errorMessage?: string | false
  /** 是否显示加载状态，默认 true */
  showLoading?: boolean
  /** 成功回调 */
  onSuccess?: (data: T) => void
  /** 失败回调 */
  onError?: (error: any) => void
  /** 最终回调（无论成功失败都会执行） */
  onFinally?: () => void
}

/**
 * 异步操作返回结果
 */
interface AsyncActionReturn {
  /** 加载状态 */
  loading: Ref<boolean>
  /** 错误信息 */
  error: Ref<any>
  /** 执行异步操作的函数 */
  execute: <T>(action: Promise<T>, options?: AsyncActionOptions<T>) => Promise<T | null>
}

/**
 * useAsyncAction Hook
 * 处理异步操作的加载状态、错误处理和提示
 * 
 * @example
 * ```ts
 * import { useAsyncAction } from '@/composables/useAsyncAction'
 * 
 * const { loading, error, execute } = useAsyncAction()
 * 
 * // 基本使用
 * const data = await execute(fetchData())
 * 
 * // 带提示消息
 * const data = await execute(fetchData(), {
 *   successMessage: '获取成功',
 *   errorMessage: '获取失败'
 * })
 * 
 * // 不显示加载状态
 * const data = await execute(fetchData(), {
 *   showLoading: false
 * })
 * 
 * // 带回调函数
 * const data = await execute(fetchData(), {
 *   onSuccess: (data) => {
 *     console.log('成功获取数据:', data)
 *   },
 *   onError: (error) => {
 *     console.error('获取失败:', error)
 *   },
 *   onFinally: () => {
 *     console.log('操作完成')
 *   }
 * })
 * 
 * // 在表单提交中使用
 * const submitForm = async () => {
 *   const result = await execute(submitApi(formData), {
 *     successMessage: '提交成功',
 *     onSuccess: () => {
 *       // 重置表单
 *       resetForm()
 *     }
 *   })
 *   
 *   if (result) {
 *     // 处理成功后的逻辑
 *   }
 * }
 * ```
 */
export function useAsyncAction(): AsyncActionReturn {
  const loading = ref(false)
  const error = ref<any>(null)

  /**
   * 执行异步操作
   * @param action 异步操作 Promise
   * @param options 配置选项
   * @returns 操作结果，失败返回 null
   */
  async function execute<T>(
    action: Promise<T>,
    options?: AsyncActionOptions<T>
  ): Promise<T | null> {
    const {
      successMessage,
      errorMessage = '操作失败，请稍后重试',
      showLoading = true,
      onSuccess,
      onError,
      onFinally
    } = options || {}

    // 重置错误状态
    error.value = null

    // 设置加载状态
    if (showLoading) {
      loading.value = true
    }

    try {
      // 执行异步操作
      const data = await action

      // 成功提示
      if (successMessage !== false) {
        toast.success(successMessage || '操作成功')
      }

      // 成功回调
      onSuccess?.(data)

      return data
    } catch (err: any) {
      // 设置错误状态
      error.value = err

      // 错误提示
      if (errorMessage !== false) {
        const message = err?.message || err || errorMessage
        toast.error(message)
      }

      // 错误回调
      onError?.(err)

      return null
    } finally {
      // 恢复加载状态
      if (showLoading) {
        loading.value = false
      }

      // 最终回调
      onFinally?.()
    }
  }

  return {
    loading,
    error,
    execute
  }
}

export default useAsyncAction
