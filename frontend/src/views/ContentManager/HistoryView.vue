<template>
  <div>
    <!-- 页面头部 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center" style="gap: 2px;">
        <el-icon class="text-primary-500"><Timer /></el-icon>
        创作记录
      </h1>
      <p class="text-gray-500 mt-1">查看内容的版本历史和修改记录</p>
    </div>

    <!-- 内容卡片 -->
    <div class="bg-white rounded-xl shadow-xiaohongshu p-5">
      <!-- 空状态 -->
      <div v-if="!hasData && !loading" class="text-center py-16">
        <div class="w-24 h-24 mx-auto mb-6 bg-gray-100 rounded-full flex items-center justify-center">
          <svg class="w-12 h-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-600 mb-2">暂无历史记录</h3>
        <p class="text-gray-400">当您创建或修改内容后，历史记录会显示在这里</p>
      </div>

      <!-- 历史记录列表 -->
      <div v-else class="space-y-4">
        <div
          v-for="item in historyData"
          :key="item.id"
          class="border border-gray-100 rounded-xl p-4 hover:border-xiaohongshu-red/30 hover:bg-xiaohongshu-red/5 transition-all"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <span class="px-3 py-1 rounded-full text-xs font-medium" :class="getTypeClass(item.type)">
                  {{ getTypeLabel(item.type) }}
                </span>
                <span class="text-sm text-gray-500">{{ formatDate(item.created_at) }}</span>
              </div>
              <h4 class="font-medium text-gray-800 mb-1">{{ getHistoryTitle(item) }}</h4>
              <p class="text-sm text-gray-500 line-clamp-2">{{ truncateText(item.description, 100) }}</p>
              <div v-if="item.change_reason" class="mt-2">
                <span class="text-xs text-gray-400">修改原因：{{ item.change_reason }}</span>
              </div>
            </div>
            <div class="flex items-center gap-2 ml-4">
              <el-button size="small" type="primary" link @click="handleView(item)">
                查看详情
              </el-button>
              <el-button size="small" type="success" link @click="handleRestore(item)">
                恢复此版本
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="text-center py-8">
        <el-icon class="animate-spin text-2xl text-xiaohongshu-red"><Loading /></el-icon>
        <p class="text-gray-500 mt-2">加载中...</p>
      </div>

      <!-- 分页 -->
      <div v-if="hasData && !loading" class="mt-6 flex justify-end">
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

    <!-- 小贴士 -->
    <div class="mt-6 bg-xiaohongshu-red/5 border border-xiaohongshu-red/20 rounded-xl p-4">
      <div class="flex items-start gap-3">
        <span class="text-xl">💡</span>
        <div>
          <h4 class="font-medium text-gray-800 mb-1">历史记录小贴士</h4>
          <ul class="text-sm text-gray-600 space-y-1">
            <li>• 系统会自动保存内容的每次修改记录</li>
            <li>• 您可以随时查看和恢复历史版本</li>
            <li>• 历史记录会保存原始内容、修改时间和操作人信息</li>
            <li>• 删除的内容也可以从历史记录中恢复</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- 查看详情对话框 -->
    <el-dialog
      v-model="viewDialogVisible"
      title="历史版本详情"
      width="700px"
      :close-on-click-modal="false"
    >
      <div v-if="currentHistory" class="space-y-4">
        <div class="flex items-center gap-3">
          <span class="px-3 py-1 rounded-full text-xs font-medium" :class="getTypeClass(currentHistory.type)">
            {{ getTypeLabel(currentHistory.type) }}
          </span>
          <span class="text-sm text-gray-500">{{ formatDate(currentHistory.created_at) }}</span>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">标题</label>
          <p class="text-gray-800">{{ currentHistory.title }}</p>
        </div>
        <div v-if="currentHistory.title_options && parseTitleOptions(currentHistory.title_options).length > 0">
          <label class="block text-sm font-medium text-gray-700 mb-1">备选标题</label>
          <div class="space-y-1">
            <el-tag
              v-for="(title, index) in parseTitleOptions(currentHistory.title_options)"
              :key="index"
              :type="index === currentHistory.selected_title_index ? 'success' : 'info'"
              size="small"
            >
              {{ index === currentHistory.selected_title_index ? '✓ ' : '' }}{{ title }}
            </el-tag>
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">内容</label>
          <div class="bg-gray-50 rounded-lg p-4 text-gray-800 whitespace-pre-wrap">
            {{ currentHistory.description }}
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">标签</label>
          <div class="flex flex-wrap gap-1">
            <el-tag
              v-for="(tag, index) in parseTags(currentHistory.tags)"
              :key="index"
              size="small"
              type="info"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>
        <div v-if="currentHistory.change_reason">
          <label class="block text-sm font-medium text-gray-700 mb-1">修改原因</label>
          <p class="text-gray-600">{{ currentHistory.change_reason }}</p>
        </div>
      </div>
      <template #footer>
        <el-button @click="viewDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="handleRestore(currentHistory)">恢复此版本</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Timer, Loading } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { getHistoryList, getHistoryDetail, restoreHistory, type ContentHistory } from '@/api/index'

const loading = ref(false)
const hasData = ref(false)
const viewDialogVisible = ref(false)
const currentHistory = ref<ContentHistory | null>(null)

const historyData = ref<ContentHistory[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 获取类型标签文本
const getTypeLabel = (type: string) => {
  const map: Record<string, string> = {
    create: '新建',
    edit: '修改',
    delete: '删除',
    publish: '发布'
  }
  return map[type] || type
}

// 获取类型标签样式类
const getTypeClass = (type: string) => {
  const map: Record<string, string> = {
    create: 'bg-green-100 text-green-700',
    edit: 'bg-blue-100 text-blue-700',
    delete: 'bg-red-100 text-red-700',
    publish: 'bg-purple-100 text-purple-700'
  }
  return map[type] || 'bg-gray-100 text-gray-700'
}

// 获取历史记录标题
const getHistoryTitle = (item: ContentHistory) => {
  const typeLabel = getTypeLabel(item.type)
  return `${typeLabel}了「${item.title}」`
}

// 解析标签
const parseTags = (tagsStr: string) => {
  try {
    const tags = JSON.parse(tagsStr)
    return Array.isArray(tags) ? tags.slice(0, 3) : []
  } catch {
    return []
  }
}

// 解析备选标题
const parseTitleOptions = (titleOptionsStr: string) => {
  try {
    const titles = JSON.parse(titleOptionsStr)
    return Array.isArray(titles) ? titles : []
  } catch {
    return []
  }
}

// 截断文本
const truncateText = (text: string, maxLength: number) => {
  if (!text) return ''
  return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
}

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return dayjs(dateStr).format('YYYY-MM-DD HH:mm')
}

// 加载历史记录列表
const loadHistoryList = async () => {
  loading.value = true
  try {
    const res = await getHistoryList({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    
    historyData.value = res.data.list
    pagination.total = res.data.total
    hasData.value = res.data.list.length > 0
  } catch (error) {
    console.error('加载历史记录失败:', error)
    ElMessage.error('加载历史记录失败')
  } finally {
    loading.value = false
  }
}

// 查看详情
const handleView = async (item: ContentHistory) => {
  try {
    const res = await getHistoryDetail(item.id)
    currentHistory.value = res.data
    viewDialogVisible.value = true
  } catch (error) {
    console.error('获取历史记录详情失败:', error)
    ElMessage.error('获取历史记录详情失败')
  }
}

// 恢复版本
const handleRestore = (item: ContentHistory) => {
  ElMessageBox.confirm('确定要恢复到此版本吗？当前内容将被覆盖。', '版本恢复', {
    confirmButtonText: '确定恢复',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await restoreHistory(item.id)
      ElMessage.success('版本恢复成功')
      viewDialogVisible.value = false
      loadHistoryList()
    } catch (error) {
      console.error('恢复版本失败:', error)
      ElMessage.error('恢复版本失败')
    }
  }).catch(() => {
  })
}

// 分页大小改变
const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadHistoryList()
}

// 页码改变
const handlePageChange = (page: number) => {
  pagination.page = page
  loadHistoryList()
}

onMounted(() => {
  loadHistoryList()
})
</script>
