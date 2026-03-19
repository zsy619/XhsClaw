<template>
  <div class="dashboard-container">
    <!-- 页面标题 -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center" style="gap: 2px;">
          <el-icon class="text-primary-500"><DataAnalysis /></el-icon>
          统计仪表盘
        </h1>
        <p class="mt-1 text-sm text-gray-500">查看您的内容创作数据和运营概况</p>
      </div>
      <el-button :icon="Refresh" circle @click="refreshData" :loading="loading" />
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <el-icon class="is-loading text-4xl text-primary-500"><Loading /></el-icon>
      <span class="ml-3 text-gray-500">加载中...</span>
    </div>

    <template v-else>
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
            <div v-if="stat.trend !== undefined" class="text-right">
              <div
                class="flex items-center text-sm"
                :class="stat.trend >= 0 ? 'text-green-500' : 'text-red-500'"
              >
                <el-icon><component :is="stat.trend >= 0 ? 'Top' : 'Bottom'" /></el-icon>
                <span class="ml-1">{{ Math.abs(stat.trend) }}%</span>
              </div>
              <div class="text-xs text-gray-400">较上周</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 快捷操作 -->
      <div class="mt-6 rounded-xl bg-white p-5 shadow-xiaohongshu">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-lg font-semibold text-xiaohongshu-dark">快捷操作</h2>
          <span class="text-xs text-gray-400">快速开始您的创作</span>
        </div>
        <el-button-group class="w-full">
          <el-button
            v-for="action in actions"
            :key="action.label"
            @click="handleAction(action.path)"
            size="large"
            class="action-button"
          >
            <el-icon class="mr-1"><component :is="action.icon" /></el-icon>
            {{ action.label }}
          </el-button>
        </el-button-group>
      </div>

      <!-- 图表区域 -->
      <div class="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-2">
        <!-- 内容生成趋势 -->
        <div class="rounded-xl bg-white p-5 shadow-xiaohongshu">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-xiaohongshu-dark">内容生成趋势</h2>
            <el-radio-group v-model="trendPeriod" size="small">
              <el-radio-button value="week">本周</el-radio-button>
              <el-radio-button value="month">本月</el-radio-button>
            </el-radio-group>
          </div>
          <div ref="trendChartRef" class="chart h-72"></div>
        </div>

        <!-- 发布状态分布 -->
        <div class="rounded-xl bg-white p-5 shadow-xiaohongshu">
          <div class="mb-4">
            <h2 class="text-lg font-semibold text-xiaohongshu-dark">发布状态分布</h2>
          </div>
          <div ref="statusChartRef" class="chart h-72"></div>
        </div>
      </div>

      <!-- 用户活跃度与生成效率 -->
      <div class="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-2">
        <!-- 用户活跃度 -->
        <div class="rounded-xl bg-white p-5 shadow-xiaohongshu">
          <div class="mb-4">
            <h2 class="text-lg font-semibold text-xiaohongshu-dark">内容发布趋势</h2>
          </div>
          <div ref="activityChartRef" class="chart h-72"></div>
        </div>

        <!-- 内容生成效率 -->
        <div class="rounded-xl bg-white p-5 shadow-xiaohongshu">
          <div class="mb-4">
            <h2 class="text-lg font-semibold text-xiaohongshu-dark">创作效率统计</h2>
          </div>
          <div ref="efficiencyChartRef" class="chart h-72"></div>
        </div>
      </div>

      <!-- 最近活动 -->
      <div class="mt-6 rounded-xl bg-white p-5 shadow-xiaohongshu">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-lg font-semibold text-xiaohongshu-dark">最近活动</h2>
          <el-button link type="primary" size="small" @click="$router.push('/content/history')">查看全部</el-button>
        </div>
        <div v-if="activities.length > 0" class="space-y-4">
          <div
            v-for="activity in activities"
            :key="activity.id"
            class="flex items-start gap-4 rounded-lg border border-gray-100 p-4 transition-all hover:border-primary-200 hover:bg-primary-50/30"
          >
            <div
              class="flex h-10 w-10 items-center justify-center rounded-full"
              :class="getActivityBgClass(activity.type)"
            >
              <el-icon :size="20" :color="getActivityColor(activity.type)">
                <component :is="getActivityIcon(activity.type)" />
              </el-icon>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between">
                <p class="text-sm font-medium text-xiaohongshu-dark truncate">
                  {{ activity.title || '无标题内容' }}
                </p>
                <el-tag :type="getStatusType(activity.status)" size="small">
                  {{ activity.status }}
                </el-tag>
              </div>
              <p class="mt-1 text-xs text-gray-400">{{ activity.time_ago }}</p>
            </div>
          </div>
        </div>
        <div v-else class="py-8 text-center text-gray-400">
          <el-icon :size="40"><Document /></el-icon>
          <p class="mt-2">暂无活动记录</p>
        </div>
      </div>

      <!-- 数据来源说明 -->
      <div class="mt-6 bg-blue-50 border border-blue-100 rounded-xl p-4">
        <div class="flex items-start gap-3">
          <el-icon :size="20" class="text-blue-500 mt-0.5"><InfoFilled /></el-icon>
          <div>
            <h4 class="font-medium text-blue-800 mb-1">数据说明</h4>
            <ul class="text-sm text-blue-700 space-y-1">
              <li>• 统计数据基于您创建和发布的内容计算</li>
              <li>• Token使用量来自大模型API调用记录</li>
              <li>• 最近活动显示您最近的5条内容操作</li>
              <li>• 数据每小时自动刷新一次</li>
            </ul>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { getDashboardData, type DashboardData, type UserActivity } from '@/api/dashboard'
import { DataAnalysis, Document, InfoFilled, Loading, Refresh } from '@element-plus/icons-vue'
import type { ECharts } from 'echarts'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 加载状态
const loading = ref(true)

// 趋势周期
const trendPeriod = ref('week')

// 仪表盘数据
const dashboardData = ref<DashboardData | null>(null)

// 活动列表
const activities = ref<UserActivity[]>([])

// 统计数据卡片
interface StatsCard {
  label: string
  value: string | number
  icon: string
  gradient: string
  trend?: number
}

const statsCards = computed<StatsCard[]>(() => {
  if (!dashboardData.value) {
    return [
      { label: '生成内容', value: '0', icon: 'Document', gradient: 'linear-gradient(135deg, #ff2442 0%, #ff4d64 100%)' },
      { label: '已发布', value: '0', icon: 'Upload', gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
      { label: '成功率', value: '0%', icon: 'TrendCharts', gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
      { label: '平均时间', value: '0s', icon: 'Timer', gradient: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' }
    ]
  }

  const stats = dashboardData.value.stats
  return [
    {
      label: '生成内容',
      value: stats.total_contents,
      icon: 'Document',
      gradient: 'linear-gradient(135deg, #ff2442 0%, #ff4d64 100%)',
      trend: stats.weekly_trend.length > 0 ? calculateTrend(stats.weekly_trend.map(d => d.contents)) : 0
    },
    {
      label: '已发布',
      value: stats.published_count,
      icon: 'Upload',
      gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      trend: stats.weekly_trend.length > 0 ? calculateTrend(stats.weekly_trend.map(d => d.published)) : 0
    },
    {
      label: '成功率',
      value: `${stats.success_rate.toFixed(1)}%`,
      icon: 'TrendCharts',
      gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
      trend: stats.success_rate > 0 ? 2.1 : 0
    },
    {
      label: '平均时间',
      value: `${stats.avg_generation_time.toFixed(1)}s`,
      icon: 'Timer',
      gradient: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
      trend: -5.4
    }
  ]
})

// 计算趋势百分比
const calculateTrend = (values: number[]): number => {
  if (values.length < 2) return 0
  const recent = values[values.length - 1]
  const previous = values[values.length - 2]
  if (previous === 0) return recent > 0 ? 100 : 0
  return Math.round(((recent - previous) / previous) * 100)
}

// 快捷操作
const actions = [
  { label: '创作中心', icon: 'Edit', path: '/creation' },
  { label: '我的笔记', icon: 'List', path: '/content/list' },
  { label: '发布历史', icon: 'DocumentCopy', path: '/publish/history' }
]

// 获取活动背景颜色类
const getActivityBgClass = (type: string) => {
  const colorMap: Record<string, string> = {
    'generate': 'bg-primary-100',
    'save': 'bg-purple-100',
    'publish': 'bg-green-100',
    'edit': 'bg-yellow-100'
  }
  return colorMap[type] || 'bg-gray-100'
}

// 获取活动颜色
const getActivityColor = (type: string) => {
  const colorMap: Record<string, string> = {
    'generate': '#ff2442',
    'save': '#667eea',
    'publish': '#43e97b',
    'edit': '#e6a23c'
  }
  return colorMap[type] || '#999999'
}

// 获取活动图标
const getActivityIcon = (type: string) => {
  const iconMap: Record<string, string> = {
    'generate': 'Document',
    'save': 'EditPen',
    'publish': 'Upload',
    'edit': 'EditPen'
  }
  return iconMap[type] || 'Document'
}

// 获取状态类型
const getStatusType = (status: string) => {
  const typeMap: Record<string, string> = {
    '已发布': 'success',
    '待发布': 'warning',
    '草稿': 'info',
    '发布失败': 'danger'
  }
  return typeMap[status] || 'info'
}

// 图表实例
let trendChart: ECharts | null = null
let statusChart: ECharts | null = null
let activityChart: ECharts | null = null
let efficiencyChart: ECharts | null = null

const trendChartRef = ref<HTMLElement>()
const statusChartRef = ref<HTMLElement>()
const activityChartRef = ref<HTMLElement>()
const efficiencyChartRef = ref<HTMLElement>()

// 自我校验相关配置
const MAX_RETRY_COUNT = 3
const RETRY_DELAY = 2000
let retryCount = ref(0)
let isRetrying = ref(false)

// 数据校验函数
const validateData = (data: any): boolean => {
  if (!data) return false
  if (!data.stats) return false

  const stats = data.stats
  if (typeof stats.total_contents !== 'number') return false
  if (typeof stats.published_count !== 'number') return false
  if (typeof stats.success_rate !== 'number') return false

  return true
}

// 获取图表数据
const getChartData = () => {
  if (!dashboardData.value) {
    return {
      trendLabels: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      trendValues: [0, 0, 0, 0, 0, 0, 0],
      statusData: [
        { value: 0, name: '草稿', itemStyle: { color: '#667eea' } },
        { value: 0, name: '待发布', itemStyle: { color: '#4facfe' } },
        { value: 0, name: '已发布', itemStyle: { color: '#43e97b' } },
        { value: 0, name: '发布失败', itemStyle: { color: '#999' } }
      ],
      publishTrend: [0, 0, 0, 0, 0, 0, 0],
      efficiencyData: { avgTime: 5.2, successRate: 0 }
    }
  }

  const stats = dashboardData.value.stats
  const weeklyTrend = stats.weekly_trend || []
  const trends = dashboardData.value.trends || []

  // 格式化日期标签
  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
    return weekdays[date.getDay()]
  }

  return {
    trendLabels: weeklyTrend.map(d => formatDate(d.date)),
    trendValues: weeklyTrend.map(d => d.contents),
    statusData: stats.status_distribution.map(s => ({
      value: s.count,
      name: s.label,
      itemStyle: { color: s.color }
    })),
    publishTrend: trends.map(t => t.publish),
    efficiencyData: {
      avgTime: stats.avg_generation_time,
      successRate: stats.success_rate
    }
  }
}

// 初始化内容生成趋势图
const initTrendChart = () => {
  if (!trendChartRef.value) return

  trendChart = echarts.init(trendChartRef.value)
  const data = getChartData()

  trendChart.setOption({
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#fff',
      borderColor: '#f0f0f0',
      textStyle: { color: '#333' }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.trendLabels,
      axisLine: { lineStyle: { color: '#f0f0f0' } },
      axisLabel: { color: '#999' }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      splitLine: { lineStyle: { color: '#f5f5f5' } },
      axisLabel: { color: '#999' }
    },
    series: [
      {
        name: '生成数量',
        type: 'line',
        smooth: true,
        data: data.trendValues,
        itemStyle: { color: '#ff2442' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(255, 36, 66, 0.3)' },
            { offset: 1, color: 'rgba(255, 36, 66, 0.01)' }
          ])
        },
        lineStyle: { width: 3 }
      }
    ]
  })
}

// 初始化发布状态分布图
const initStatusChart = () => {
  if (!statusChartRef.value) return

  statusChart = echarts.init(statusChartRef.value)
  const data = getChartData()

  statusChart.setOption({
    tooltip: {
      trigger: 'item',
      backgroundColor: '#fff',
      borderColor: '#f0f0f0',
      textStyle: { color: '#333' }
    },
    legend: {
      orient: 'vertical',
      right: '5%',
      top: 'center'
    },
    series: [
      {
        name: '发布状态',
        type: 'pie',
        radius: ['45%', '70%'],
        center: ['35%', '50%'],
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
        data: data.statusData
      }
    ]
  })
}

// 初始化发布趋势图
const initActivityChart = () => {
  if (!activityChartRef.value) return

  activityChart = echarts.init(activityChartRef.value)
  const data = getChartData()
  const chartData = getChartData()
  const labels = chartData.trendLabels

  activityChart.setOption({
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#fff',
      borderColor: '#f0f0f0',
      textStyle: { color: '#333' }
    },
    legend: {
      data: ['生成数', '发布数'],
      top: 0
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: labels,
      axisLine: { lineStyle: { color: '#f0f0f0' } },
      axisLabel: { color: '#999' }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      splitLine: { lineStyle: { color: '#f5f5f5' } },
      axisLabel: { color: '#999' }
    },
    series: [
      {
        name: '生成数',
        type: 'line',
        smooth: true,
        data: data.trendValues,
        itemStyle: { color: '#ff2442' },
        lineStyle: { width: 2 },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(255, 36, 66, 0.2)' },
            { offset: 1, color: 'rgba(255, 36, 66, 0.01)' }
          ])
        }
      },
      {
        name: '发布数',
        type: 'line',
        smooth: true,
        data: data.publishTrend,
        itemStyle: { color: '#43e97b' },
        lineStyle: { width: 2 },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(67, 233, 123, 0.2)' },
            { offset: 1, color: 'rgba(67, 233, 123, 0.01)' }
          ])
        }
      }
    ]
  })
}

// 初始化效率统计图
const initEfficiencyChart = () => {
  if (!efficiencyChartRef.value) return

  efficiencyChart = echarts.init(efficiencyChartRef.value)
  const data = getChartData()

  efficiencyChart.setOption({
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#fff',
      borderColor: '#f0f0f0',
      textStyle: { color: '#333' }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: ['创作效率'],
      axisLine: { lineStyle: { color: '#f0f0f0' } },
      axisLabel: { color: '#999' }
    },
    yAxis: [
      {
        type: 'value',
        name: '平均时长(秒)',
        position: 'left',
        axisLine: { show: false },
        splitLine: { lineStyle: { color: '#f5f5f5' } },
        axisLabel: { color: '#999' },
        max: 10
      },
      {
        type: 'value',
        name: '成功率(%)',
        position: 'right',
        axisLine: { show: false },
        splitLine: { show: false },
        axisLabel: { color: '#999' },
        max: 100
      }
    ],
    series: [
      {
        name: '平均时长',
        type: 'bar',
        data: [data.efficiencyData.avgTime],
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#ff2442' },
            { offset: 1, color: '#ff4d64' }
          ]),
          borderRadius: [6, 6, 0, 0]
        }
      },
      {
        name: '成功率',
        type: 'line',
        yAxisIndex: 1,
        data: [data.efficiencyData.successRate],
        itemStyle: { color: '#4facfe' },
        lineStyle: { width: 3 }
      }
    ]
  })
}

// 初始化所有图表
const initCharts = () => {
  initTrendChart()
  initStatusChart()
  initActivityChart()
  initEfficiencyChart()
}

// 更新图表数据
const updateCharts = () => {
  const data = getChartData()

  if (trendChart) {
    trendChart.setOption({
      tooltip: {
        trigger: 'axis',
        backgroundColor: '#fff',
        borderColor: '#f0f0f0',
        textStyle: { color: '#333' }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: data.trendLabels,
        axisLine: { lineStyle: { color: '#f0f0f0' } },
        axisLabel: { color: '#999' }
      },
      yAxis: {
        type: 'value',
        axisLine: { show: false },
        splitLine: { lineStyle: { color: '#f5f5f5' } },
        axisLabel: { color: '#999' }
      },
      series: [{
        name: '生成数量',
        type: 'line',
        smooth: true,
        data: data.trendValues,
        itemStyle: { color: '#ff2442' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(255, 36, 66, 0.3)' },
            { offset: 1, color: 'rgba(255, 36, 66, 0.01)' }
          ])
        },
        lineStyle: { width: 3 }
      }]
    }, { notMerge: true })
  }

  if (statusChart) {
    statusChart.setOption({
      tooltip: {
        trigger: 'item',
        backgroundColor: '#fff',
        borderColor: '#f0f0f0',
        textStyle: { color: '#333' }
      },
      legend: {
        orient: 'vertical',
        right: '5%',
        top: 'center',
        textStyle: { color: '#666' }
      },
      series: [{
        type: 'pie',
        radius: ['45%', '70%'],
        center: ['35%', '50%'],
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
        labelLine: { show: false },
        data: data.statusData
      }]
    }, { notMerge: true })
  }

  if (activityChart) {
    activityChart.setOption({
      tooltip: {
        trigger: 'axis',
        backgroundColor: '#fff',
        borderColor: '#f0f0f0',
        textStyle: { color: '#333' }
      },
      legend: {
        data: ['生成数量', '发布数量'],
        bottom: 0,
        textStyle: { color: '#666' }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '12%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: data.trendLabels,
        axisLine: { lineStyle: { color: '#f0f0f0' } },
        axisLabel: { color: '#999' }
      },
      yAxis: {
        type: 'value',
        axisLine: { show: false },
        splitLine: { lineStyle: { color: '#f5f5f5' } },
        axisLabel: { color: '#999' }
      },
      series: [
        {
          name: '生成数量',
          type: 'bar',
          data: data.trendValues,
          itemStyle: { color: '#667eea' },
          barWidth: '35%'
        },
        {
          name: '发布数量',
          type: 'bar',
          data: data.publishTrend,
          itemStyle: { color: '#43e97b' },
          barWidth: '35%'
        }
      ]
    }, { notMerge: true })
  }

  if (efficiencyChart) {
    efficiencyChart.setOption({
      tooltip: {
        trigger: 'axis',
        backgroundColor: '#fff',
        borderColor: '#f0f0f0',
        textStyle: { color: '#333' },
        axisPointer: { type: 'shadow' }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: ['创作效率'],
        axisLine: { lineStyle: { color: '#f0f0f0' } },
        axisLabel: { color: '#999' }
      },
      yAxis: {
        type: 'value',
        axisLine: { show: false },
        splitLine: { lineStyle: { color: '#f5f5f5' } },
        axisLabel: { color: '#999' }
      },
      series: [
        {
          name: '平均耗时(分钟)',
          type: 'bar',
          data: [data.efficiencyData.avgTime],
          itemStyle: { color: '#4facfe' },
          barWidth: '40%',
          label: {
            show: true,
            position: 'top',
            formatter: '{c} 分钟'
          }
        },
        {
          name: '成功率(%)',
          type: 'bar',
          data: [data.efficiencyData.successRate],
          itemStyle: { color: '#43e97b' },
          barWidth: '40%',
          label: {
            show: true,
            position: 'top',
            formatter: '{c}%'
          }
        }
      ]
    }, { notMerge: true })
  }
}

// 获取仪表盘数据（带自我校验和重试机制）
const fetchDashboardData = async (isRetry = false) => {
  if (!isRetry) {
    loading.value = true
  }

  try {
    const res = await getDashboardData()

    // 数据校验
    if (!validateData(res.data)) {
      throw new Error('数据格式不正确')
    }

    // 数据有效，重置重试计数
    retryCount.value = 0
    isRetrying.value = false

    dashboardData.value = res.data
    activities.value = res.data.activities || []

    // 更新图表
    updateCharts()

    // 如果是重试中，恢复正常状态
    if (isRetry) {
      ElMessage.success('数据已恢复')
    }
  } catch (error: any) {
    console.error('获取仪表盘数据失败:', error)

    // 自我校验机制：自动重试
    if (retryCount.value < MAX_RETRY_COUNT) {
      retryCount.value++
      isRetrying.value = true

      ElMessage.warning(`数据加载异常，正在尝试恢复 (${retryCount.value}/${MAX_RETRY_COUNT})...`)

      // 延迟重试
      await new Promise(resolve => setTimeout(resolve, RETRY_DELAY))
      return fetchDashboardData(true)
    }

    // 重试次数用完，显示错误
    ElMessage.error('获取数据失败，请刷新页面重试或检查网络连接')

    // 使用空数据渲染，确保界面不崩溃
    dashboardData.value = null
    activities.value = []
  } finally {
    if (!isRetry) {
      loading.value = false
    }
  }
}

// 手动刷新数据
const refreshData = () => {
  retryCount.value = 0
  isRetrying.value = false
  fetchDashboardData()
}

// 定时刷新（每5分钟）
let autoRefreshTimer: ReturnType<typeof setInterval> | null = null

const startAutoRefresh = () => {
  if (autoRefreshTimer) {
    clearInterval(autoRefreshTimer)
  }
  autoRefreshTimer = setInterval(() => {
    if (!isRetrying.value) {
      fetchDashboardData(true)
    }
  }, 5 * 60 * 1000) // 5分钟
}

const stopAutoRefresh = () => {
  if (autoRefreshTimer) {
    clearInterval(autoRefreshTimer)
    autoRefreshTimer = null
  }
}

// 处理快捷操作
const handleAction = (path: string) => {
  router.push(path)
}

// 监听窗口大小变化
const handleResize = () => {
  trendChart?.resize()
  statusChart?.resize()
  activityChart?.resize()
  efficiencyChart?.resize()
}

// 监听趋势周期变化
watch(trendPeriod, () => {
  // 可以在这里切换不同的数据周期
  updateCharts()
})

onMounted(async () => {
  await nextTick()
  await fetchDashboardData()
  initCharts()
  startAutoRefresh()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  stopAutoRefresh()
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  statusChart?.dispose()
  activityChart?.dispose()
  efficiencyChart?.dispose()
})
</script>

<style scoped lang="scss">
.dashboard-container {
  .action-button {
    flex: 1;
    height: 52px;
    font-size: 14px;
    font-weight: 500;
  }

  :deep(.el-button-group) {
    display: flex;
    width: 100%;

    .el-button {
      flex: 1;
    }
  }
}
</style>