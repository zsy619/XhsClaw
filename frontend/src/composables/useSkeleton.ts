/**
 * 骨架屏 Composable
 * 提供加载状态的骨架屏效果
 */
import { ref, computed, type Ref } from 'vue'

/**
 * 骨架屏配置
 */
interface SkeletonOptions {
  /** 是否显示骨架屏 */
  show?: boolean
  /** 骨架屏行数，默认 3 */
  rows?: number
  /** 是否显示头像，默认 false */
  avatar?: boolean
  /** 是否显示标题，默认 true */
  title?: boolean
  /** 是否显示段落，默认 true */
  paragraph?: boolean
  /** 自定义加载状态 */
  loading?: Ref<boolean>
}

/**
 * useSkeleton Hook
 * 管理骨架屏状态
 * 
 * @example
 * ```ts
 * // 基本使用
 * const { showSkeleton } = useSkeleton()
 * 
 * // 自定义配置
 * const { showSkeleton } = useSkeleton({
 *   rows: 5,
 *   avatar: true,
 *   title: true
 * })
 * ```
 */
export function useSkeleton(options: SkeletonOptions = {}) {
  const {
    rows = 3,
    avatar = false,
    title = true,
    paragraph = true,
    loading
  } = options

  // 内部加载状态
  const internalLoading = ref(true)

  // 计算最终的加载状态
  const showSkeleton = computed(() => {
    if (loading !== undefined) {
      return loading.value
    }
    return internalLoading.value
  })

  /**
   * 设置加载状态
   */
  function setLoading(value: boolean) {
    internalLoading.value = value
  }

  /**
   * 开始加载
   */
  function start() {
    internalLoading.value = true
  }

  /**
   * 结束加载
   */
  function stop() {
    internalLoading.value = false
  }

  return {
    showSkeleton,
    setLoading,
    start,
    stop,
    // 配置参数
    skeletonConfig: {
      rows,
      avatar,
      title,
      paragraph
    }
  }
}

export default useSkeleton
