/**
 * 字符统计 Composable
 * 提供文本编辑器的字符统计功能
 */
import { ref, computed, watch, type Ref } from 'vue'

/**
 * 统计信息
 */
interface TextStats {
  /** 字符数（包含空格） */
  characters: number
  /** 字符数（不包含空格） */
  charactersNoSpace: number
  /** 单词数 */
  words: number
  /** 行数 */
  lines: number
  /** 段落数 */
  paragraphs: number
  /** 中文字符数 */
  chineseCharacters: number
  /** 英文字符数 */
  englishCharacters: number
  /** 数字字符数 */
  numberCharacters: number
  /** 特殊字符数 */
  specialCharacters: number
}

/**
 * 配置选项
 */
interface TextStatsOptions {
  /** 初始文本 */
  initialText?: string
  /** 是否监听变化，默认 true */
  watchChanges?: boolean
  /** 最大字符数限制 */
  maxLength?: number
  /** 字符数变化回调 */
  onCountChange?: (stats: TextStats) => void
}

/**
 * useTextStats Hook
 * 统计文本字符信息
 * 
 * @example
 * ```ts
 * import { useTextStats } from '@/composables/useTextStats'
 * 
 * const { text, stats, setStats } = useTextStats()
 * 
 * // 基本使用
 * const { text, stats } = useTextStats({
 *   initialText: 'Hello World',
 *   maxLength: 1000,
 *   onCountChange: (stats) => {
 *     console.log('字符数:', stats.characters)
 *   }
 * })
 * 
 * // 手动设置文本
 * setText('New content')
 * 
 * // 获取统计信息
 * console.log(stats.value.characters) // 字符数
 * console.log(stats.value.words) // 单词数
 * console.log(stats.value.lines) // 行数
 * ```
 */
export function useTextStats(options: TextStatsOptions = {}) {
  const {
    initialText = '',
    watchChanges = true,
    maxLength,
    onCountChange
  } = options

  // 文本内容
  const text = ref(initialText)
  
  // 统计信息
  const stats = ref<TextStats>({
    characters: 0,
    charactersNoSpace: 0,
    words: 0,
    lines: 0,
    paragraphs: 0,
    chineseCharacters: 0,
    englishCharacters: 0,
    numberCharacters: 0,
    specialCharacters: 0
  })

  // 是否超过限制
  const isOverLimit = computed(() => {
    if (!maxLength) return false
    return stats.value.characters > maxLength
  })

  // 剩余字符数
  const remainingChars = computed(() => {
    if (!maxLength) return null
    return maxLength - stats.value.characters
  })

  /**
   * 计算统计信息
   */
  function calculateStats(content: string): TextStats {
    const characters = content.length
    const charactersNoSpace = content.replace(/\s/g, '').length
    const words = content.trim() ? content.trim().split(/\s+/).length : 0
    const lines = content ? content.split(/\r\n|\r|\n/).length : 0
    const paragraphs = content ? content.split(/\n\s*\n/).filter(p => p.trim()).length : 0
    
    // 中文字符统计
    const chineseMatches = content.match(/[\u4e00-\u9fa5]/g)
    const chineseCharacters = chineseMatches ? chineseMatches.length : 0
    
    // 英文字符统计
    const englishMatches = content.match(/[a-zA-Z]/g)
    const englishCharacters = englishMatches ? englishMatches.length : 0
    
    // 数字字符统计
    const numberMatches = content.match(/[0-9]/g)
    const numberCharacters = numberMatches ? numberMatches.length : 0
    
    // 特殊字符统计（排除空格、中文、英文、数字）
    const specialMatches = content.match(/[^\u4e00-\u9fa5a-zA-Z0-9\s]/g)
    const specialCharacters = specialMatches ? specialMatches.length : 0

    return {
      characters,
      charactersNoSpace,
      words,
      lines,
      paragraphs,
      chineseCharacters,
      englishCharacters,
      numberCharacters,
      specialCharacters
    }
  }

  /**
   * 更新统计信息
   */
  function updateStats(content: string = text.value) {
    stats.value = calculateStats(content)
    
    // 触发回调
    onCountChange?.(stats.value)
  }

  /**
   * 设置文本内容
   */
  function setText(newText: string) {
    // 检查是否超过限制
    if (maxLength && newText.length > maxLength) {
      newText = newText.substring(0, maxLength)
    }
    
    text.value = newText
    updateStats(newText)
  }

  /**
   * 追加文本
   */
  function appendText(newText: string) {
    const combined = text.value + newText
    setText(combined)
  }

  /**
   * 清空文本
   */
  function clearText() {
    text.value = ''
    updateStats('')
  }

  /**
   * 获取统计信息
   */
  function getStats(): TextStats {
    return stats.value
  }

  /**
   * 获取格式化的统计信息
   */
  function getFormattedStats(): string {
    const s = stats.value
    return `字符：${s.characters} | 单词：${s.words} | 行：${s.lines} | 段落：${s.paragraphs}`
  }

  // 监听文本变化
  if (watchChanges) {
    watch(text, (newVal) => {
      updateStats(newVal)
    }, { immediate: true })
  }

  // 初始化统计
  if (initialText) {
    updateStats(initialText)
  }

  return {
    text,
    stats,
    isOverLimit,
    remainingChars,
    setText,
    appendText,
    clearText,
    getStats,
    getFormattedStats,
    updateStats
  }
}

export default useTextStats
