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
      <div v-if="loading" class="text-center py-16">
        <el-icon class="w-16 h-16 text-gray-300" :size="64">
          <Loading />
        </el-icon>
        <p class="text-gray-500 mt-4">加载中...</p>
      </div>
      
      <div v-else-if="publishHistory.length === 0" class="text-center py-16">
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
            :timestamp="formatDate(item.scheduled_at || item.published_at)"
            placement="top"
            :color="getStatusColor(item.status)"
            :size="getStatusSize(item.status)"
          >
            <div class="bg-gray-50 rounded-xl p-5 border border-gray-100 hover:border-xiaohongshu-red/30 transition-all">
              <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
                <!-- 左侧信息 -->
                <div class="flex-1">
                  <div class="flex items-center gap-3 mb-3">
                    <div
                      class="w-10 h-10 rounded-full flex items-center justify-center"
                      :class="getStatusBgClass(item.status)"
                    >
                      <el-icon
                        :size="20"
                        :color="getStatusIconColor(item.status)"
                      >
                        <component :is="getStatusIcon(item.status)" />
                      </el-icon>
                    </div>
                    <div class="flex-1">
                      <h4 class="font-semibold text-gray-800 text-lg">{{ item.content?.title || '未命名内容' }}</h4>
                      <div class="flex items-center gap-3 mt-1">
                        <el-tag
                          size="small"
                          :type="getStatusTagType(item.status)"
                        >
                          {{ getStatusText(item.status) }}
                        </el-tag>
                        <span class="text-sm text-gray-500">平台：小红书</span>
                      </div>
                    </div>
                  </div>

                  <!-- 错误信息 -->
                  <div v-if="item.status === 3 && item.error_msg" class="mt-3 p-3 bg-red-50 rounded-lg border border-red-100">
                    <div class="flex items-start gap-2">
                      <el-icon class="text-red-500 mt-0.5" :size="16"><Warning /></el-icon>
                      <div>
                        <span class="font-medium text-red-700">错误信息：</span>
                        <span class="text-red-600">{{ item.error_msg }}</span>
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
                    v-if="item.status === 0"
                    size="small"
                    type="warning"
                    @click="handleCancel(item)"
                  >
                    取消发布
                  </el-button>
                  <el-button
                    v-if="item.status === 3"
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

    <!-- 发布详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="发布详情"
      width="720px"
      :close-on-click-modal="false"
    >
      <div v-loading="detailLoading" class="detail-content">
        <div v-if="currentRecord" class="space-y-6">
          <!-- 基本信息 -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <p class="text-sm text-gray-500 mb-1">发布记录ID</p>
              <p class="font-medium text-gray-800">{{ currentRecord.id }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500 mb-1">发布状态</p>
              <el-tag :type="getStatusTagType(currentRecord.status)">
                {{ getStatusText(currentRecord.status) }}
              </el-tag>
            </div>
          </div>

          <!-- 内容信息 -->
          <div class="bg-gray-50 rounded-lg p-4">
            <p class="text-sm text-gray-500 mb-2">关联内容</p>
            <p class="font-medium text-gray-800">{{ currentRecord.content?.title || '未命名内容' }}</p>
            <p class="text-sm text-gray-600 mt-2 line-clamp-3">{{ currentRecord.content?.description || '暂无描述' }}</p>
          </div>

          <!-- 时间信息 -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <p class="text-sm text-gray-500 mb-1">创建时间</p>
              <p class="font-medium text-gray-800">{{ formatDate(currentRecord.created_at) }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500 mb-1">计划发布时间</p>
              <p class="font-medium text-gray-800">{{ formatDate(currentRecord.scheduled_at) }}</p>
            </div>
            <div v-if="currentRecord.published_at">
              <p class="text-sm text-gray-500 mb-1">实际发布时间</p>
              <p class="font-medium text-gray-800">{{ formatDate(currentRecord.published_at) }}</p>
            </div>
          </div>

          <!-- 错误信息 -->
          <div v-if="currentRecord.status === 3 && currentRecord.error_msg" class="bg-red-50 rounded-lg p-4 border border-red-200">
            <p class="text-sm font-medium text-red-700 mb-2">错误信息</p>
            <p class="text-sm text-red-600">{{ currentRecord.error_msg }}</p>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end">
          <el-button @click="showDetailDialog = false">关闭</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { DocumentCopy, Warning, Loading, CircleCheck, CircleClose, Clock } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, reactive, ref, onMounted } from 'vue'
import { getPublishRecords, cancelPublish, retryPublish, getPublishRecord, type PublishRecord } from '@/api/publish'

const loading = ref(false)
const publishHistory = ref<PublishRecord[]>([])
const showDetailDialog = ref(false)
const currentRecord = ref<PublishRecord | null>(null)
const detailLoading = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const stats = computed(() => {
  const total = publishHistory.value.length
  const success = publishHistory.value.filter(item => item.status === 2).length
  const failed = publishHistory.value.filter(item => item.status === 3).length
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

const getStatusColor = (status: number) => {
  switch (status) {
    case 0: return '#f59e0b' // 待发布 - 黄色
    case 1: return '#3b82f6' // 发布中 - 蓝色
    case 2: return '#67c23a' // 成功 - 绿色
    case 3: return '#f56c6c' // 失败 - 红色
    default: return '#909399' // 默认 - 灰色
  }
}

const getStatusSize = (status: number) => {
  return status === 2 ? 'large' : 'large'
}

const getStatusBgClass = (status: number) => {
  switch (status) {
    case 0: return 'bg-yellow-100'
    case 1: return 'bg-blue-100'
    case 2: return 'bg-green-100'
    case 3: return 'bg-red-100'
    default: return 'bg-gray-100'
  }
}

const getStatusIconColor = (status: number) => {
  switch (status) {
    case 0: return '#f59e0b'
    case 1: return '#3b82f6'
    case 2: return '#67c23a'
    case 3: return '#f56c6c'
    default: return '#909399'
  }
}

const getStatusIcon = (status: number) => {
  switch (status) {
    case 0: return Clock
    case 1: return Loading
    case 2: return CircleCheck
    case 3: return CircleClose
    default: return Clock
  }
}

const getStatusTagType = (status: number) => {
  switch (status) {
    case 0: return 'warning'
    case 1: return 'primary'
    case 2: return 'success'
    case 3: return 'danger'
    default: return 'info'
  }
}

const getStatusText = (status: number) => {
  switch (status) {
    case 0: return '待发布'
    case 1: return '发布中'
    case 2: return '发布成功'
    case 3: return '发布失败'
    default: return '未知状态'
  }
}

const loadPublishHistory = async () => {
  loading.value = true
  try {
    const res = await getPublishRecords({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res.data) {
      publishHistory.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('加载发布历史失败:', error)
    ElMessage.error('加载发布历史失败')
  } finally {
    loading.value = false
  }
}

const handleView = async (item: PublishRecord) => {
  detailLoading.value = true
  try {
    const res = await getPublishRecord(item.id)
    if (res.data) {
      currentRecord.value = res.data
      showDetailDialog.value = true
    }
  } catch (error) {
    console.error('获取详情失败:', error)
    ElMessage.error('获取详情失败')
  } finally {
    detailLoading.value = false
  }
}

const handleCancel = async (item: PublishRecord) => {
  try {
    await ElMessageBox.confirm('确定要取消发布吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await cancelPublish(item.id)
    ElMessage.success('已取消发布')
    loadPublishHistory()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消发布失败:', error)
      ElMessage.error('取消发布失败')
    }
  }
}

const handleRetry = async (item: PublishRecord) => {
  try {
    await ElMessageBox.confirm('确定要重试发布吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await retryPublish(item.id)
    ElMessage.success('已开始重试发布')
    loadPublishHistory()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重试发布失败:', error)
      ElMessage.error('重试发布失败')
    }
  }
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadPublishHistory()
}

const handlePageChange = (page: number) => {
  pagination.page = page
  loadPublishHistory()
}

onMounted(() => {
  loadPublishHistory()
})
</script>
