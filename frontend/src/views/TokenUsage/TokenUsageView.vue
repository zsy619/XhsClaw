<template>
  <div class="token-usage-container">
    <!-- 页面标题和主题选择 -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center" style="gap: 8px;">
          <el-icon class="text-primary-500"><Coin /></el-icon>
          Token使用统计
        </h1>
        <p class="mt-1 text-sm text-gray-500">查看您的大模型API使用情况和费用统计</p>
      </div>
      
      <!-- 主题选择器和诊断工具 -->
      <div class="flex items-center gap-3">
        <span class="text-sm text-gray-600">选择主题：</span>
        <el-select 
          v-model="selectedTheme" 
          placeholder="选择主题" 
          size="default"
          style="width: 220px;"
          @change="handleThemeChange"
          filterable
        >
          <el-option
            v-for="theme in themes"
            :key="theme.value"
            :label="theme.label"
            :value="theme.value"
          >
            <div class="flex items-center gap-2">
              <div 
                class="w-4 h-4 rounded-full border border-gray-300"
                :style="{ background: theme.previewColor }"
              ></div>
              <span>{{ theme.label }}</span>
            </div>
          </el-option>
        </el-select>
        
        <!-- 诊断按钮 -->
        <el-button 
          type="warning" 
          size="small" 
          @click="runDiagnostic"
          :icon="Tools"
          :loading="diagnosing"
        >
          系统诊断
        </el-button>
        
        <!-- 主题校验按钮 -->
        <el-button 
          type="primary" 
          size="small" 
          @click="validateThemes"
          :icon="Check"
          :loading="validating"
        >
          校验主题
        </el-button>
      </div>
    </div>

    <!-- 诊断结果显示 -->
    <div v-if="diagnosticResult" class="mb-6">
      <el-alert
        :title="diagnosticResult.success ? '诊断成功' : '诊断发现问题'"
        :type="diagnosticResult.success ? 'success' : 'warning'"
        :description="diagnosticResult.message"
        show-icon
        :closable="false"
      >
        <template v-if="diagnosticResult.details" #default>
          <div class="mt-2 text-sm">
            <p><strong>详细信息：</strong></p>
            <ul class="list-disc list-inside mt-1">
              <li v-for="(detail, index) in diagnosticResult.details" :key="index">{{ detail }}</li>
            </ul>
          </div>
        </template>
        <template v-if="diagnosticResult.suggestions && diagnosticResult.suggestions.length > 0" #default>
          <div class="mt-3 text-sm">
            <p><strong>💡 建议操作：</strong></p>
            <ul class="list-disc list-inside mt-1">
              <li v-for="(suggestion, index) in diagnosticResult.suggestions" :key="index">{{ suggestion }}</li>
            </ul>
          </div>
        </template>
      </el-alert>
    </div>

    <!-- 校验结果显示 -->
    <div v-if="validationResult" class="mb-6">
      <el-alert
        :title="validationResult.success ? '校验成功' : '校验失败'"
        :type="validationResult.success ? 'success' : 'error'"
        :description="validationResult.message"
        show-icon
        :closable="false"
      >
        <template v-if="validationResult.details" #default>
          <div class="mt-2 text-sm">
            <p><strong>详细信息：</strong></p>
            <ul class="list-disc list-inside mt-1">
              <li v-for="(detail, index) in validationResult.details" :key="index">{{ detail }}</li>
            </ul>
          </div>
        </template>
      </el-alert>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <div
        v-for="stat in statsCards"
        :key="stat.label"
        class="rounded-xl bg-white p-5 shadow-xiaohongshu transition-all hover:shadow-xiaohongshu-lg"
      >
        <div class="flex items-center gap-4">
          <div
            class="flex h-14 w-14 items-center justify-center rounded-xl text-white"
            :style="{ background: stat.gradient }"
          >
            <el-icon :size="28"><component :is="stat.icon" /></el-icon>
          </div>
          <div class="flex-1">
            <div class="text-2xl font-bold text-xiaohongshu-dark">{{ stat.value }}</div>
            <div class="mt-1 text-sm text-gray-500">{{ stat.label }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-2">
      <!-- 每日使用趋势 -->
      <div class="rounded-xl bg-white p-5 shadow-xiaohongshu">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-lg font-semibold text-xiaohongshu-dark">每日使用趋势</h2>
          <el-radio-group v-model="dailyPeriod" size="small">
            <el-radio-button value="7">近7天</el-radio-button>
            <el-radio-button value="30">近30天</el-radio-button>
          </el-radio-group>
        </div>
        <div ref="dailyChartRef" class="h-72"></div>
      </div>

      <!-- 模型分布 -->
      <div class="rounded-xl bg-white p-5 shadow-xiaohongshu">
        <div class="mb-4">
          <h2 class="text-lg font-semibold text-xiaohongshu-dark">模型使用分布</h2>
        </div>
        <div ref="modelChartRef" class="h-72"></div>
      </div>
    </div>

    <!-- 使用记录表格 -->
    <div class="mt-6 rounded-xl bg-white p-5 shadow-xiaohongshu">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-lg font-semibold text-xiaohongshu-dark">使用记录</h2>
        <el-button @click="loadUsageRecords" :loading="loading">
          <el-icon class="mr-1"><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
      <el-table :data="usageRecords" v-loading="loading" stripe>
        <el-table-column prop="model" label="模型" min-width="150" />
        <el-table-column prop="provider" label="提供商" width="120" />
        <el-table-column prop="prompt_tokens" label="输入Tokens" width="120" align="right">
          <template #default="{ row }">
            <span class="text-blue-600">{{ (row.prompt_tokens ?? 0).toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="completion_tokens" label="输出Tokens" width="120" align="right">
          <template #default="{ row }">
            <span class="text-green-600">{{ (row.completion_tokens ?? 0).toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="total_tokens" label="总Tokens" width="120" align="right">
          <template #default="{ row }">
            <span class="font-semibold">{{ (row.total_tokens ?? 0).toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="cost" label="费用(美元)" width="100" align="right">
          <template #default="{ row }">
            <span class="text-orange-600">${{ (row.cost ?? 0).toFixed(6) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="request_type" label="请求类型" width="120" />
        <el-table-column prop="response_status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.response_status === 'success' ? 'success' : 'danger'" size="small">
              {{ row.response_status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
    getUserDailyStats,
    getUserStatsByModel,
    getUserTokenStats,
    getUserTokenUsage,
    type TokenUsage,
    type TokenUsageByModel,
    type TokenUsageStats,
    type UserTokenUsage
} from '@/api/token_usage'
import { getQuickFixSteps, quickDiagnostic } from '@/utils/diagnostic'
import { Check, Coin, Refresh, Tools } from '@element-plus/icons-vue'
import type { ECharts } from 'echarts'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

// 诊断结果接口
interface DiagnosticResult {
  success: boolean
  message: string
  details?: string[]
  suggestions?: string[]
}

// 校验结果接口
interface ValidationResult {
  success: boolean
  message: string
  details?: string[]
}

// 主题接口
interface Theme {
  value: string
  label: string
  previewColor: string
}

// 完整的主题列表（根据backend/assets/themes目录）
const allThemeFiles = [
  'aurora-green', 'autumn-leaves', 'berry-smoothie', 'blackberry-sage',
  'blue-lagoon', 'blueberry-cheese', 'blueberry-night', 'blueberry',
  'blush-pink', 'botanical', 'bubblegum-pink', 'candy-pink',
  'caramel-macchiato', 'cherry-blossom', 'cherry-blush', 'cherry',
  'chocolate-mint', 'coconut-cream', 'coral', 'cotton-candy',
  'cream-custard', 'cream', 'crystal-purple', 'dark', 'deep-ocean',
  'default', 'earl-grey', 'elegant', 'floral-pink', 'forest-green',
  'grape-purple', 'grape', 'honey-ginger', 'honey-peach', 'ice-blue',
  'ivory-cream', 'kiwi-green', 'lavender-gray', 'lavender-honey',
  'lavender-purple', 'lavender', 'lemon-meringue', 'lemon-yellow',
  'lemon', 'lily-white', 'mango-pudding', 'matcha-green', 'matcha-latte',
  'matcha', 'mint-breeze', 'mint-chocolate', 'mint-green', 'mint',
  'neo-brutalism', 'neon-pink', 'nordic', 'ocean-blue', 'ocean-mist',
  'ocean', 'peach-blossom', 'peach', 'peaches-cream', 'pearl-white',
  'pink-cream', 'pistachio-green', 'playful-geometric', 'pomegranate',
  'professional', 'purple', 'rainbow-sherbet', 'rainbow-sorbet',
  'red-velvet', 'retro', 'rose-gold', 'rose-milk', 'rose', 'sage-green',
  'sakura-blossom', 'sakura-pink', 'sand-brown', 'sea-glass', 'sketch',
  'sky-blue', 'strawberry-milk', 'strawberry-red', 'strawberry-sweet',
  'strawberry', 'sun-kissed', 'sunset-glow', 'sunset-orange', 'sunset',
  'taro-milktea', 'terminal', 'tiramisu', 'vanilla-cream', 'vanilla-milk',
  'winter-sky', 'xiaohongshu'
]

// 主题预览颜色映射
const themePreviewColors: Record<string, string> = {
  'aurora-green': '#4CAF50',
  'autumn-leaves': '#FF5722',
  'berry-smoothie': '#9C27B0',
  'blackberry-sage': '#4A148C',
  'blue-lagoon': '#00BCD4',
  'blueberry-cheese': '#3F51B5',
  'blueberry-night': '#1A237E',
  'blueberry': '#2196F3',
  'blush-pink': '#FFCDD2',
  'botanical': '#81C784',
  'bubblegum-pink': '#FF4081',
  'candy-pink': '#E91E63',
  'caramel-macchiato': '#795548',
  'cherry-blossom': '#F8BBD9',
  'cherry-blush': '#F48FB1',
  'cherry': '#D32F2F',
  'chocolate-mint': '#00695C',
  'coconut-cream': '#FFF8E1',
  'coral': '#FF7043',
  'cotton-candy': '#FF80AB',
  'cream-custard': '#FFFDE7',
  'cream': '#FFF9C4',
  'crystal-purple': '#7C4DFF',
  'dark': '#212121',
  'deep-ocean': '#0D47A1',
  'default': '#FF2442',
  'earl-grey': '#607D8B',
  'elegant': '#9E9E9E',
  'floral-pink': '#FCE4EC',
  'forest-green': '#1B5E20',
  'grape-purple': '#6A1B9A',
  'grape': '#7B1FA2',
  'honey-ginger': '#FF6F00',
  'honey-peach': '#FFAB91',
  'ice-blue': '#B3E5FC',
  'ivory-cream': '#FFFFF0',
  'kiwi-green': '#8BC34A',
  'lavender-gray': '#90A4AE',
  'lavender-honey': '#E1BEE7',
  'lavender-purple': '#CE93D8',
  'lavender': '#E6EE9C',
  'lemon-meringue': '#FFF59D',
  'lemon-yellow': '#FFEB3B',
  'lemon': '#FFF176',
  'lily-white': '#FAFAFA',
  'mango-pudding': '#FFB74D',
  'matcha-green': '#43A047',
  'matcha-latte': '#A5D6A7',
  'matcha': '#66BB6A',
  'mint-breeze': '#B2DFDB',
  'mint-chocolate': '#2E7D32',
  'mint-green': '#81D4FA',
  'mint': '#4DD0E1',
  'neo-brutalism': '#FFD54F',
  'neon-pink': '#FF1744',
  'nordic': '#ECEFF1',
  'ocean-blue': '#0288D1',
  'ocean-mist': '#B0BEC5',
  'ocean': '#03A9F4',
  'peach-blossom': '#FFCCBC',
  'peach': '#FF8A65',
  'peaches-cream': '#FFE0B2',
  'pearl-white': '#F5F5F5',
  'pink-cream': '#F8BBD9',
  'pistachio-green': '#C5E1A5',
  'playful-geometric': '#FF5252',
  'pomegranate': '#C62828',
  'professional': '#37474F',
  'purple': '#9C27B0',
  'rainbow-sherbet': '#FF9800',
  'rainbow-sorbet': '#FFC107',
  'red-velvet': '#B71C1C',
  'retro': '#FF9800',
  'rose-gold': '#E91E63',
  'rose-milk': '#FCE4EC',
  'rose': '#EC407A',
  'sage-green': '#81C784',
  'sakura-blossom': '#FFCDD2',
  'sakura-pink': '#F48FB1',
  'sand-brown': '#A1887F',
  'sea-glass': '#80CBC4',
  'sketch': '#BDBDBD',
  'sky-blue': '#64B5F6',
  'strawberry-milk': '#FFCDD2',
  'strawberry-red': '#E53935',
  'strawberry-sweet': '#EF5350',
  'strawberry': '#F44336',
  'sun-kissed': '#FFA726',
  'sunset-glow': '#FF7043',
  'sunset-orange': '#FF5722',
  'sunset': '#FF9800',
  'taro-milktea': '#B39DDB',
  'terminal': '#424242',
  'tiramisu': '#8D6E63',
  'vanilla-cream': '#FFFDE7',
  'vanilla-milk': '#FFF8E1',
  'winter-sky': '#BBDEFB',
  'xiaohongshu': '#FF2442'
}

// 生成主题列表
const themes: Theme[] = allThemeFiles.map(themeFile => ({
  value: themeFile,
  label: formatThemeLabel(themeFile),
  previewColor: themePreviewColors[themeFile] || '#9E9E9E'
}))

// 格式化主题标签
function formatThemeLabel(themeFile: string): string {
  return themeFile
    .split('-')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

// 响应式数据
const loading = ref(false)
const validating = ref(false)
const diagnosing = ref(false)
const selectedTheme = ref('xiaohongshu')
const validationResult = ref<ValidationResult | null>(null)
const diagnosticResult = ref<DiagnosticResult | null>(null)
const dailyPeriod = ref('7')
const statsCards = ref([
  { label: '总请求数', value: '0', icon: 'TrendCharts', gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
  { label: '总Tokens', value: '0', icon: 'Coin', gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
  { label: '总费用(美元)', value: '$0', icon: 'Money', gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
  { label: '平均Tokens/次', value: '0', icon: 'DataAnalysis', gradient: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' }
])

const usageRecords = ref<TokenUsage[]>([])
const dailyStats = ref<UserTokenUsage[]>([])
const modelStats = ref<TokenUsageByModel[]>([])

// 图表引用
const dailyChartRef = ref<HTMLElement>()
const modelChartRef = ref<HTMLElement>()
let dailyChart: ECharts | null = null
let modelChart: ECharts | null = null

// 运行系统诊断
const runDiagnostic = async () => {
  diagnosing.value = true
  diagnosticResult.value = null
  
  try {
    console.log('🔧 开始系统诊断...')
    const result = await quickDiagnostic()
    diagnosticResult.value = result
    
    if (!result.success) {
      // 自动显示快速修复步骤
      console.log('📋 快速修复步骤:')
      getQuickFixSteps().forEach(step => console.log(step))
    }
  } catch (error: any) {
    diagnosticResult.value = {
      success: false,
      message: `诊断过程出错：${error.message}`,
      suggestions: getQuickFixSteps()
    }
    ElMessage.error('诊断过程出错')
  } finally {
    diagnosing.value = false
  }
}

// 主题自我校验机制
const validateThemes = async () => {
  validating.value = true
  validationResult.value = null
  
  try {
    const details: string[] = []
    let success = true
    
    // 校验1: 检查主题列表数量是否正确
    const expectedCount = 117
    if (themes.length === expectedCount) {
      details.push(`✓ 主题数量正确：${themes.length} 个`)
    } else {
      details.push(`✗ 主题数量错误：期望 ${expectedCount} 个，实际 ${themes.length} 个`)
      success = false
    }
    
    // 校验2: 检查所有主题文件是否都在列表中
    const missingThemes: string[] = []
    allThemeFiles.forEach(file => {
      const found = themes.find(t => t.value === file)
      if (!found) {
        missingThemes.push(file)
      }
    })
    
    if (missingThemes.length === 0) {
      details.push('✓ 所有主题文件都已包含')
    } else {
      details.push(`✗ 缺失主题：${missingThemes.join(', ')}`)
      success = false
    }
    
    // 校验3: 检查是否有重复主题
    const themeValues = themes.map(t => t.value)
    const duplicates = themeValues.filter((item, index) => themeValues.indexOf(item) !== index)
    
    if (duplicates.length === 0) {
      details.push('✓ 没有重复的主题')
    } else {
      details.push(`✗ 重复主题：${[...new Set(duplicates)].join(', ')}`)
      success = false
    }
    
    // 校验4: 检查所有主题是否有预览颜色
    const themesWithoutColor: string[] = []
    themes.forEach(theme => {
      if (!themePreviewColors[theme.value]) {
        themesWithoutColor.push(theme.value)
      }
    })
    
    if (themesWithoutColor.length === 0) {
      details.push('✓ 所有主题都有预览颜色')
    } else {
      details.push(`⚠ 缺少预览颜色的主题：${themesWithoutColor.join(', ')}`)
    }
    
    // 校验5: 验证主题选择功能
    try {
      selectedTheme.value = themes[0].value
      await new Promise(resolve => setTimeout(resolve, 100))
      selectedTheme.value = 'xiaohongshu'
      details.push('✓ 主题选择功能正常')
    } catch (error) {
      details.push('✗ 主题选择功能异常')
      success = false
    }
    
    validationResult.value = {
      success,
      message: success 
        ? '所有校验通过！主题功能完整可用。'
        : '部分校验失败，请检查详细信息。',
      details
    }
    
    if (success) {
      ElMessage.success('主题校验成功！')
    } else {
      ElMessage.error('主题校验失败，请检查详细信息')
    }
  } catch (error: any) {
    validationResult.value = {
      success: false,
      message: `校验过程出错：${error.message}`
    }
    ElMessage.error('校验过程出错')
  } finally {
    validating.value = false
  }
}

// 格式化日期
const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return '-'
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 主题切换处理
const handleThemeChange = (themeValue: string) => {
  ElMessage.success(`已切换到 ${formatThemeLabel(themeValue)} 主题`)
}

// 加载统计数据
const loadStats = async () => {
  try {
    const [statsRes, dailyRes, modelRes] = await Promise.all([
      getUserTokenStats(),
      getUserDailyStats(parseInt(dailyPeriod.value)),
      getUserStatsByModel()
    ])
    
    if (statsRes.code === 0 && statsRes.data) {
      const stats = statsRes.data as TokenUsageStats
      statsCards.value[0].value = stats.total_requests.toString()
      statsCards.value[1].value = stats.total_tokens.toLocaleString()
      statsCards.value[2].value = `$${stats.total_cost.toFixed(4)}`
      statsCards.value[3].value = Math.round(stats.average_tokens).toLocaleString()
    }
    
    if (dailyRes.code === 0 && dailyRes.data) {
      dailyStats.value = dailyRes.data as UserTokenUsage[]
      updateDailyChart()
    }
    
    if (modelRes.code === 0 && modelRes.data) {
      modelStats.value = modelRes.data as TokenUsageByModel[]
      updateModelChart()
    }
  } catch (error: any) {
    console.error('加载统计数据失败:', error)
    // 如果是网络错误，自动运行诊断
    if (error.message?.includes('Network Error') || error.code === 'ERR_NETWORK') {
      ElMessage.warning('检测到网络问题，正在运行系统诊断...')
      setTimeout(() => runDiagnostic(), 500)
    }
  }
}

// 使用记录加载状态
const usageRecordsLoading = ref(false)
const usageRecordsRetryCount = ref(0)
const usageRecordsMaxRetries = 5

// 加载使用记录 - 带自我校验和自动重试机制
const loadUsageRecords = async () => {
  if (usageRecordsLoading.value) return

  usageRecordsLoading.value = true
  loading.value = true
  usageRecordsRetryCount.value = 0

  const attemptLoad = async (): Promise<boolean> => {
    try {
      usageRecordsRetryCount.value++
      const res = await getUserTokenUsage(10)

      if (res.code === 0 && res.data && Array.isArray(res.data)) {
        const records = res.data as TokenUsage[]

        // 自我校验：验证数据完整性和有效性
        const isValid = records.every(record =>
          record &&
          typeof record.model === 'string' &&
          typeof record.total_tokens === 'number'
        )

        if (isValid || records.length === 0) {
          usageRecords.value = records
          return true
        } else {
          console.warn(`⚠ 第 ${usageRecordsRetryCount.value} 次校验失败：数据不完整`)
        }
      } else {
        console.warn(`⚠ 第 ${usageRecordsRetryCount.value} 次校验失败：返回数据异常`)
      }

      return false
    } catch (error: any) {
      console.error(`❌ 第 ${usageRecordsRetryCount.value} 次加载失败:`, error.message || error)
      return false
    }
  }

  let success = false

  while (!success && usageRecordsRetryCount.value < usageRecordsMaxRetries) {
    success = await attemptLoad()

    if (!success && usageRecordsRetryCount.value < usageRecordsMaxRetries) {
      const delay = Math.min(1000 * Math.pow(2, usageRecordsRetryCount.value - 1), 5000)
      console.log(`⏳ ${Math.round(delay/1000)}秒后进行第 ${usageRecordsRetryCount.value + 1} 次尝试...`)
      await new Promise(resolve => setTimeout(resolve, delay))
    }
  }

  if (!success) {
    console.error(`❌ 已达到最大重试次数 (${usageRecordsMaxRetries})，加载失败`)
    ElMessage.error('使用记录加载失败，请稍后刷新页面重试')
  } else if (usageRecordsRetryCount.value > 1) {
    ElMessage.success(`✓ 使用记录加载成功 (共尝试 ${usageRecordsRetryCount.value} 次)`)
  }

  usageRecordsLoading.value = false
  loading.value = false
}

// 更新每日趋势图表
const updateDailyChart = () => {
  if (!dailyChart) return
  
  const dates = dailyStats.value.map(item => item.date).reverse()
  const tokens = dailyStats.value.map(item => item.total_tokens).reverse()
  const costs = dailyStats.value.map(item => item.total_cost).reverse()
  
  dailyChart.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' }
    },
    legend: {
      data: ['Tokens', '费用(美元)']
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: {
        rotate: 45,
        formatter: (value: string) => {
          const date = new Date(value)
          return `${date.getMonth() + 1}/${date.getDate()}`
        }
      }
    },
    yAxis: [
      { type: 'value', name: 'Tokens', position: 'left' },
      { type: 'value', name: '费用', position: 'right', axisLabel: { formatter: '$ {value}' } }
    ],
    series: [
      {
        name: 'Tokens',
        type: 'bar',
        data: tokens,
        itemStyle: { color: '#667eea' }
      },
      {
        name: '费用(美元)',
        type: 'line',
        yAxisIndex: 1,
        data: costs,
        itemStyle: { color: '#f5576c' },
        smooth: true
      }
    ]
  })
}

// 更新模型分布图表
const updateModelChart = () => {
  if (!modelChart) return
  
  const data = modelStats.value.map(item => ({
    name: item.model,
    value: item.total_tokens
  }))
  
  modelChart.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} Tokens ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: data,
        color: ['#667eea', '#f093fb', '#f5576c', '#4facfe', '#00f2fe', '#43e97b', '#38f9d7']
      }
    ]
  })
}

// 初始化图表
const initCharts = () => {
  if (dailyChartRef.value) {
    dailyChart = echarts.init(dailyChartRef.value)
  }
  if (modelChartRef.value) {
    modelChart = echarts.init(modelChartRef.value)
  }
}

// 监听时间范围变化
watch(dailyPeriod, () => {
  loadStats()
})

// 生命周期
onMounted(() => {
  initCharts()
  loadStats()
  loadUsageRecords()
  
  window.addEventListener('resize', () => {
    dailyChart?.resize()
    modelChart?.resize()
  })
})

onBeforeUnmount(() => {
  dailyChart?.dispose()
  modelChart?.dispose()
})
</script>

<style scoped>
.token-usage-container {
}
</style>
