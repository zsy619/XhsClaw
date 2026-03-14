<template>
  <div class="dashboard-container">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center" style="gap: 2px;">
        <el-icon class="text-primary-500"><DataAnalysis /></el-icon>
        统计仪表盘
      </h1>
      <p class="mt-1 text-sm text-gray-500">查看您的内容创作数据和运营概况</p>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <div
        v-for="stat in stats"
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
          <div v-if="stat.trend" class="text-right">
            <div
              class="flex items-center text-sm"
              :class="stat.trend > 0 ? 'text-green-500' : 'text-red-500'"
            >
              <el-icon><component :is="stat.trend > 0 ? 'Top' : 'Bottom'" /></el-icon>
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
          <h2 class="text-lg font-semibold text-xiaohongshu-dark">用户活跃度</h2>
        </div>
        <div ref="activityChartRef" class="chart h-72"></div>
      </div>

      <!-- 内容生成效率 -->
      <div class="rounded-xl bg-white p-5 shadow-xiaohongshu">
        <div class="mb-4">
          <h2 class="text-lg font-semibold text-xiaohongshu-dark">内容生成效率</h2>
        </div>
        <div ref="efficiencyChartRef" class="chart h-72"></div>
      </div>
    </div>

    <!-- 最近活动 -->
    <div class="mt-6 rounded-xl bg-white p-5 shadow-xiaohongshu">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-lg font-semibold text-xiaohongshu-dark">最近活动</h2>
        <el-button link type="primary" size="small">查看全部</el-button>
      </div>
      <div class="space-y-4">
        <div
          v-for="activity in activities"
          :key="activity.id"
          class="flex items-start gap-4 rounded-lg border border-gray-100 p-4 transition-all hover:border-primary-200 hover:bg-primary-50/30"
        >
          <div
            class="flex h-10 w-10 items-center justify-center rounded-full"
            :class="getActivityBgClass(activity.color)"
          >
            <el-icon :size="20" :color="activity.color">
              <component :is="activity.icon" />
            </el-icon>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between">
              <p class="text-sm font-medium text-xiaohongshu-dark">
                {{ activity.text }}
              </p>
              <el-tag :type="activity.status" size="small">
                {{ activity.statusText }}
              </el-tag>
            </div>
            <p class="mt-1 text-xs text-gray-400">{{ activity.time }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import type { ECharts } from 'echarts'
import { Top, Bottom, Edit, List, Timer, Connection, DocumentCopy, Setting } from '@element-plus/icons-vue'

const router = useRouter()

// 趋势周期
const trendPeriod = ref('week')

// 统计数据
const stats = reactive([
  {
    label: '生成内容',
    value: '156',
    icon: 'Document',
    gradient: 'linear-gradient(135deg, #ff2442 0%, #ff4d64 100%)',
    trend: 12.5
  },
  {
    label: '已发布',
    value: '89',
    icon: 'Upload',
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    trend: 8.2
  },
  {
    label: '成功率',
    value: '94.2%',
    icon: 'TrendCharts',
    gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    trend: 2.1
  },
  {
    label: '平均时间',
    value: '5.2s',
    icon: 'Timer',
    gradient: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    trend: -5.4
  }
])

// 快捷操作 - 只保留最常用的功能
const actions = [
  { label: '创作中心', icon: 'Edit', path: '/creation' },
  { label: '我的笔记', icon: 'List', path: '/content/list' },
  { label: '发布历史', icon: 'DocumentCopy', path: '/publish/history' }
]

// 最近活动
const activities = [
  {
    id: 1,
    text: '生成了主题内容 "Python 效率工具推荐"',
    time: '5 分钟前',
    icon: 'Document',
    color: '#ff2442',
    status: 'success',
    statusText: '生成成功'
  },
  {
    id: 2,
    text: '创建了图片 "美食探店分享"',
    time: '10 分钟前',
    icon: 'Picture',
    color: '#667eea',
    status: 'success',
    statusText: '图片已生成'
  },
  {
    id: 3,
    text: '发布了内容 "5 款效率工具推荐"',
    time: '1 小时前',
    icon: 'Upload',
    color: '#43e97b',
    status: 'success',
    statusText: '发布成功'
  },
  {
    id: 4,
    text: '改写了内容 "周末旅游攻略"',
    time: '2 小时前',
    icon: 'EditPen',
    color: '#e6a23c',
    status: 'warning',
    statusText: '待发布'
  }
]

// 获取活动背景颜色类
const getActivityBgClass = (color: string) => {
  const colorMap: Record<string, string> = {
    '#ff2442': 'bg-primary-100',
    '#667eea': 'bg-purple-100',
    '#43e97b': 'bg-green-100',
    '#e6a23c': 'bg-yellow-100'
  }
  return colorMap[color] || 'bg-gray-100'
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

// 小红书主题色
const xiaohongshuColors = ['#ff2442', '#ff4d64', '#667eea', '#4facfe', '#43e97b']

// 初始化图表
const initCharts = () => {
  // 内容生成趋势图
  if (trendChartRef.value) {
    trendChart = echarts.init(trendChartRef.value)
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
        data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
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
          data: [12, 20, 15, 8, 7, 11, 13],
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

  // 发布状态分布图
  if (statusChartRef.value) {
    statusChart = echarts.init(statusChartRef.value)
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
          data: [
            { value: 89, name: '已发布', itemStyle: { color: '#ff2442' } },
            { value: 45, name: '草稿', itemStyle: { color: '#667eea' } },
            { value: 12, name: '待审核', itemStyle: { color: '#4facfe' } },
            { value: 10, name: '发布失败', itemStyle: { color: '#999' } }
          ]
        }
      ]
    })
  }

  // 用户活跃度图
  if (activityChartRef.value) {
    activityChart = echarts.init(activityChartRef.value)
    activityChart.setOption({
      tooltip: {
        trigger: 'axis',
        backgroundColor: '#fff',
        borderColor: '#f0f0f0',
        textStyle: { color: '#333' }
      },
      legend: {
        data: ['超级管理员', '内容管理员', '普通用户'],
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
        data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
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
          name: '超级管理员',
          type: 'line',
          smooth: true,
          data: [5, 8, 6, 4, 7, 9, 6],
          itemStyle: { color: '#ff2442' },
          lineStyle: { width: 2 }
        },
        {
          name: '内容管理员',
          type: 'line',
          smooth: true,
          data: [12, 15, 13, 10, 14, 16, 12],
          itemStyle: { color: '#667eea' },
          lineStyle: { width: 2 }
        },
        {
          name: '普通用户',
          type: 'line',
          smooth: true,
          data: [25, 30, 28, 22, 26, 32, 24],
          itemStyle: { color: '#43e97b' },
          lineStyle: { width: 2 }
        }
      ]
    })
  }

  // 内容生成效率图
  if (efficiencyChartRef.value) {
    efficiencyChart = echarts.init(efficiencyChartRef.value)
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
        data: ['方案1', '方案2', '方案3', '方案4'],
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
          axisLabel: { color: '#999' }
        },
        {
          type: 'value',
          name: '成功率(%)',
          position: 'right',
          axisLine: { show: false },
          splitLine: { show: false },
          axisLabel: { color: '#999' }
        }
      ],
      series: [
        {
          name: '平均时长',
          type: 'bar',
          data: [4.2, 3.8, 6.5, 5.9],
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
          data: [96, 95, 92, 93],
          itemStyle: { color: '#4facfe' },
          lineStyle: { width: 3 }
        }
      ]
    })
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

onMounted(() => {
  initCharts()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  statusChart?.dispose()
  activityChart?.dispose()
  efficiencyChart?.dispose()
})
</script>

<style scoped lang="scss">
.dashboard-container {
  .chart {
    width: 100%;
  }
  
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
