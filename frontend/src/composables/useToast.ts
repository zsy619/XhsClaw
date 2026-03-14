/**
 * 全局 Toast 提示 Composable
 * 提供统一的提示消息管理
 */
import { ElMessage, type MessageOptions, type MessageHandler } from 'element-plus'

/**
 * Toast 类型
 */
type ToastType = 'success' | 'error' | 'warning' | 'info'

/**
 * Toast 配置选项
 */
interface ToastOptions extends Omit<MessageOptions, 'message' | 'type'> {
  /** 自动关闭时间，单位为毫秒，默认 3000ms */
  duration?: number
  /** 是否显示关闭按钮，默认 false */
  showClose?: boolean
}

/**
 * Toast 管理器
 */
class ToastManager {
  private currentMessage: MessageHandler | null = null

  /**
   * 关闭当前显示的 toast
   */
  close() {
    if (this.currentMessage) {
      this.currentMessage.close()
      this.currentMessage = null
    }
  }

  /**
   * 显示 toast 提示
   * @param message 提示消息
   * @param type 提示类型
   * @param options 配置选项
   * @returns MessageHandler
   */
  show(message: string, type: ToastType = 'info', options?: ToastOptions): MessageHandler {
    // 关闭之前的 toast
    this.close()

    const defaultOptions: ToastOptions = {
      duration: 3000,
      showClose: false,
      ...options
    }

    // 根据类型调用不同的方法
    const messageHandler = ElMessage[type]({
      message,
      type,
      ...defaultOptions
    })

    this.currentMessage = messageHandler
    return messageHandler
  }

  /**
   * 显示成功提示
   * @param message 提示消息
   * @param options 配置选项
   */
  success(message: string, options?: ToastOptions): MessageHandler {
    return this.show(message, 'success', options)
  }

  /**
   * 显示错误提示
   * @param message 提示消息
   * @param options 配置选项
   */
  error(message: string, options?: ToastOptions): MessageHandler {
    return this.show(message, 'error', options)
  }

  /**
   * 显示警告提示
   * @param message 提示消息
   * @param options 配置选项
   */
  warning(message: string, options?: ToastOptions): MessageHandler {
    return this.show(message, 'warning', options)
  }

  /**
   * 显示信息提示
   * @param message 提示消息
   * @param options 配置选项
   */
  info(message: string, options?: ToastOptions): MessageHandler {
    return this.show(message, 'info', options)
  }
}

// 创建全局单例
export const toast = new ToastManager()

/**
 * useToast Hook
 * 在组件中使用 toast 提示
 * 
 * @example
 * ```ts
 * import { useToast } from '@/composables/useToast'
 * 
 * const { toast } = useToast()
 * 
 * // 显示成功提示
 * toast.success('操作成功')
 * 
 * // 显示错误提示
 * toast.error('操作失败')
 * 
 * // 显示警告提示
 * toast.warning('请注意')
 * 
 * // 显示信息提示
 * toast.info('提示信息')
 * 
 * // 自定义配置
 * toast.success('操作成功', { duration: 5000, showClose: true })
 * 
 * // 关闭当前 toast
 * toast.close()
 * ```
 */
export function useToast() {
  return {
    toast
  }
}

export default toast
