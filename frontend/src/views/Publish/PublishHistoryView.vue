<template>
  <div>
    <!-- 页面头部 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center" style="gap: 2px;">
        <el-icon class="text-primary-500"><DocumentCopy /></el-icon>
        发布历史
      </h1>
      <p class="mt-1 text-sm text-gray-500">查看所有内容的发布记录</p>
    </div>

    <!-- 主内容卡片 -->
    <div class="bg-white rounded-xl p-5 shadow-xiaohongshu">
      <!-- 空状态 -->
      <div v-if="publishHistory.length === 0" class="text-center py-16">
        <div class="w-24 h-24 mx-auto mb-6 bg-gray-100 rounded-full flex items-center justify-center">
          <svg class="w-12 h-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-600 mb-2">暂无发布记录</h3>
        <p class="text-gray-400">当您发布内容后，发布记录会显示在这里</p>
      </div>

      <!-- 时间线 -->
      <div v-else>
        <el-timeline>
          <el-timeline-item
            v-for="item in publishHistory"
            :key="item.id"
            :timestamp="formatDate(item.publish_time)"
            placement="top"
            :color="item.status === 'success' ? '#67c23a' : '#f56c6c'"
            :size="item.status === 'success' ? 'large' : 'large'"
          >
            <div class="bg-gray-50 rounded-xl p-5 border border-gray-100 hover:border-xiaohongshu-red/30 transition-all">
              <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
                <!-- 左侧信息 -->
                <div class="flex-1">
                  <div class="flex items-center gap-3 mb-3">
                    <div
                      class="w-10 h-10 rounded-full flex items-center justify-center"
                      :class="item.status === 'success' ? 'bg-green-100' : 'bg-red-100'"
                    >
                      <el-icon
                        :size="20"
                        :color="item.status === 'success' ? '#67c23a' : '#f56c6c'"
                      >
                        <component :is="item.status === 'success' ? 'CircleCheck' : 'CircleClose'" />
                      </el-icon>
                    </div>
                    <div class="flex-1">
                      <h4 class="font-semibold text-gray-800 text-lg">{{ item.title }}</h4>
                      <div class="flex items-center gap-3 mt-1">
                        <el-tag
                          size="small"
                          :type="item.status === 'success' ? 'success' : 'danger'"
                        >
                          {{ item.status === 'success' ? '发布成功' : '发布失败' }}
                        </el-tag>
                        <span class="text-sm text-gray-500">平台：{{ item.platform }}</span>
                      </div>
                    </div>
                  </div>

                  <!-- 错误信息 -->
                  <div v-if="item.error" class="mt-3 p-3 bg-red-50 rounded-lg border border-red-100">
                    <div class="flex items-start gap-2">
                      <el-icon class="text-red-500 mt-0.5" :size="16"><Warning /></el-icon>
                      <div>
                        <span class="font-medium text-red-700">错误信息：</span>
                        <span class="text-red-600">{{ item.error }}</span>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- 右侧操作 -->
                <div class="flex items-center gap-2">
                  <el-button
                    size="small"
                    type="primary"
                    @click="handleView(item)"
                    class="bg-xiaohongshu-red border-none hover:opacity-90"
                  >
                    查看详情
                  </el-button>
                  <el-button
                    v-if="item.status === 'failed'"
                    size="small"
                    type="warning"
                    @click="handleRetry(item)"
                  >
                    重试发布
                  </el-button>
                </div>
              </div>
            </div>
          </el-timeline-item>
        </el-timeline>

        <!-- 分页 -->
        <div class="mt-8 flex justify-end">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
          />
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="mt-6 grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl shadow-xiaohongshu p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">总发布次数</p>
            <p class="text-2xl font-bold text-gray-800 mt-1">{{ stats.total }}</p>
          </div>
          <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
            </svg>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl shadow-xiaohongshu p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">发布成功</p>
            <p class="text-2xl font-bold text-green-600 mt-1">{{ stats.success }}</p>
          </div>
          <div class="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl shadow-xiaohongshu p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">发布失败</p>
            <p class="text-2xl font-bold text-red-600 mt-1">{{ stats.failed }}</p>
          </div>
          <div class="w-12 h-12 bg-red-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl shadow-xiaohongshu p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">成功率</p>
            <p class="text-2xl font-bold text-xiaohongshu-red mt-1">{{ stats.successRate }}%</p>
          </div>
          <div class="w-12 h-12 bg-xiaohongshu-red/10 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-xiaohongshu-red" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
            </svg>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { DocumentCopy, Warning } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { computed, reactive, ref } from 'vue'

const publishHistory = ref([
  {
    id: 1,
    title: 'Python 效率工具推荐',
    platform: '小红书',
    status: 'success',
    publish_time: '2026-03-13T10:30:00Z',
    error: null
  },
  {
    id: 2,
    title: '美食探店指南',
    platform: '小红书',
    status: 'failed',
    publish_time: '2026-03-13T09:20:00Z',
    error: '网络连接超时，请检查网络后重试'
  },
  {
    id: 3,
    title: '旅游攻略',
    platform: '小红书',
    status: 'success',
    publish_time: '2026-03-12T16:45:00Z',
    error: null
  }
])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 3
})

const stats = computed(() => {
  const total = publishHistory.value.length
  const success = publishHistory.value.filter(item => item.status === 'success').length
  const failed = publishHistory.value.filter(item => item.status === 'failed').length
  const successRate = total > 0 ? Math.round((success / total) * 100) : 0

  return {
    total,
    success,
    failed,
    successRate
  }
})

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return dayjs(dateStr).format('YYYY-MM-DD HH:mm:ss')
}

const handleRefresh = () => {
  ElMessage.success('刷新成功')
}

const handleView = (item: any) => {
  ElMessage.info('查看详情功能开发中')
}

const handleRetry = (item: any) => {
  ElMessage.success('重试发布成功')
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
}

const handlePageChange = (page: number) => {
  pagination.page = page
}
</script>
