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
      
      <!-- 主题选择器 -->
      <div class="flex items-center gap-4">
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
      </div>
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
            <span class="text-blue-600">{{ row.prompt_tokens.toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="completion_tokens" label="输出Tokens" width="120" align="right">
          <template #default="{ row }">
            <span class="text-green-600">{{ row.completion_tokens.toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="total_tokens" label="总Tokens" width="120" align="right">
          <template #default="{ row }">
            <span class="font-semibold">{{ row.total_tokens.toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="cost" label="费用(美元)" width="100" align="right">
          <template #default="{ row }">
            <span class="text-orange-600">${{ row.cost.toFixed(6) }}</span>
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
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import type { ECharts } from 'echarts'
import { Coin, Refresh } from '@element-plus/icons-vue'
import { 
  getUserTokenStats, 
  getUserDailyStats, 
  getUserStatsByModel, 
  getUserTokenUsage,
  type TokenUsageStats,
  type UserTokenUsage,
  type TokenUsageByModel,
  type TokenUsage
} from '@/api/token_usage'

// 主题列表
const themes = [
  { value: 'default', label: '默认', previewColor: '#ffffff' },
  { value: 'xiaohongshu', label: '小红书', previewColor: '#ffe4e6' },
  { value: 'terminal', label: '终端', previewColor: '#1e1e1e' },
  { value: 'ocean', label: '海洋', previewColor: '#0ea5e9' },
]

// 主题预览颜色映射
const themePreviewColors: Record<string, string> = {
  default: '#ffffff',
  xiaohongshu: '#ffe4e6',
  terminal: '#1e1e1e',
  ocean: '#0ea5e9',
}

// 响应式数据
const loading = ref(false)
const dailyPeriod = ref('7')
const selectedTheme = ref('xiaohongshu')

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

// 格式化日期
const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 主题切换
const handleThemeChange = (theme: string) => {
  selectedTheme.value = theme
  ElMessage.success(`已切换主题: ${theme}`)
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
    ElMessage.error(error.message || '加载统计数据失败')
  }
}

// 加载使用记录
const loadUsageRecords = async () => {
  loading.value = true
  try {
    const res = await getUserTokenUsage(50)
    if (res.code === 0 && res.data) {
      usageRecords.value = res.data as TokenUsage[]
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载使用记录失败')
  } finally {
    loading.value = false
  }
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
  padding: 20px;
}
</style>
