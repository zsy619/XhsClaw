/**
 * 未保存离开提示 Composable
 * 防止用户在有未保存内容时意外离开页面
 */
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { ElMessageBox } from 'element-plus'

/**
 * 配置选项
 */
interface UseConfirmLeaveOptions {
  /** 是否启用，默认 true */
  enabled?: boolean
  /** 提示消息 */
  message?: string
  /** 标题 */
  title?: string
  /** 确认按钮文字 */
  confirmButtonText?: string
  /** 取消按钮文字 */
  cancelButtonText?: string
  /** 保存回调 */
  onSave?: () => Promise<void> | void
  /** 离开回调 */
  onLeave?: () => void
}

/**
 * useConfirmLeave Hook
 * 管理未保存内容的离开确认
 * 
 * @example
 * ```ts
 * import { useConfirmLeave } from '@/composables/useConfirmLeave'
 * 
 * const { isDirty, enable, disable, confirmLeave } = useConfirmLeave({
 *   message: '您有未保存的更改，确定要离开吗？',
 *   title: '确认离开',
 *   onSave: async () => {
 *     await saveData()
 *   }
 * })
 * 
 * // 标记为已修改
 * const onContentChange = () => {
 *   isDirty.value = true
 * }
 * 
 * // 保存后标记为干净
 * const onSave = async () => {
 *   await saveData()
 *   isDirty.value = false
 * }
 * ```
 */
export function useConfirmLeave(options: UseConfirmLeaveOptions = {}) {
  const {
    enabled = true,
    message = '您有未保存的更改，确定要离开吗？',
    title = '确认离开',
    confirmButtonText = '保存并离开',
    cancelButtonText = '不保存离开',
    onSave,
    onLeave
  } = options

  // 是否有未保存的更改
  const isDirty = ref(false)
  
  // 是否阻止离开
  const blockLeave = ref(false)

  /**
   * 处理离开页面
   */
  async function handleBeforeUnload(e: BeforeUnloadEvent) {
    if (!enabled || !isDirty.value) return

    // 现代浏览器需要设置 returnValue
    e.preventDefault()
    e.returnValue = message
    
    return message
  }

  /**
   * 处理路由离开
   */
  async function handleRouteLeave(): Promise<boolean> {
    if (!enabled || !isDirty.value) return true

    try {
      const action = await ElMessageBox.confirm(message, title, {
        confirmButtonText: confirmButtonText,
        cancelButtonText: cancelButtonText,
        distinguishCancelAndClose: true,
        type: 'warning'
      })

      if (action === 'confirm') {
        // 用户选择保存
        if (onSave) {
          await onSave()
        }
        isDirty.value = false
        return true
      } else {
        // 用户选择不保存
        onLeave?.()
        isDirty.value = false
        return true
      }
    } catch {
      // 用户取消操作
      return false
    }
  }

  /**
   * 启用离开确认
   */
  function enableConfirm() {
    blockLeave.value = true
  }

  /**
   * 禁用离开确认
   */
  function disableConfirm() {
    blockLeave.value = false
  }

  /**
   * 标记为已保存
   */
  function markAsSaved() {
    isDirty.value = false
  }

  /**
   * 标记为未保存
   */
  function markAsUnsaved() {
    isDirty.value = true
  }

  // 监听浏览器关闭/刷新
  onMounted(() => {
    if (enabled) {
      window.addEventListener('beforeunload', handleBeforeUnload)
    }
  })

  // 清理事件监听
  onUnmounted(() => {
    window.removeEventListener('beforeunload', handleBeforeUnload)
  })

  // 监听 enabled 状态变化
  watch(enabled, (newVal) => {
    if (newVal) {
      window.addEventListener('beforeunload', handleBeforeUnload)
    } else {
      window.removeEventListener('beforeunload', handleBeforeUnload)
    }
  })

  return {
    isDirty,
    blockLeave,
    enableConfirm,
    disableConfirm,
    markAsSaved,
    markAsUnsaved,
    handleRouteLeave
  }
}

export default useConfirmLeave
