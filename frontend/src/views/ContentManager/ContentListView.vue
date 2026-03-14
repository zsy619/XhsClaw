<template>
  <div>
    <!-- 页面头部 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center" style="gap: 2px;">
        <el-icon class="text-primary-500"><Document /></el-icon>
        我的笔记
      </h1>
      <p class="mt-1 text-sm text-gray-500">管理和查看所有生成的内容</p>
    </div>

    <!-- 内容卡片 -->
    <div class="bg-white rounded-xl p-5 shadow-xiaohongshu">
      <!-- 搜索栏 -->
      <div class="bg-gray-50 rounded-xl p-4 mb-6">
        <el-form :inline="true" :model="searchForm" class="flex flex-wrap gap-4">
          <el-form-item label="关键词">
            <el-input
              v-model="searchForm.keyword"
              placeholder="搜索标题或内容"
              clearable
              style="width: 220px"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部状态" clearable style="width: 140px">
              <el-option label="草稿" :value="0" />
              <el-option label="待发布" :value="1" />
              <el-option label="已发布" :value="2" />
              <el-option label="发布失败" :value="3" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch" class="bg-xiaohongshu-red border-none">
              搜索
            </el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格 -->
      <div class="overflow-x-auto">
        <el-table
          :data="tableData"
          v-loading="loading"
          style="width: 100%"
          class="xiaohongshu-table"
        >
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
          <el-table-column prop="description" label="内容预览" min-width="250" show-overflow-tooltip>
            <template #default="{ row }">
              <span class="text-gray-500 text-sm">{{ truncateText(row.description, 50) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="标签" min-width="180">
            <template #default="{ row }">
              <div class="flex flex-wrap gap-1">
                <el-tag
                  v-for="(tag, index) in parseTags(row.tags)"
                  :key="index"
                  size="small"
                  type="info"
                  class="mr-1 mb-1"
                >
                  {{ tag }}
                </el-tag>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusTag(row.status)" size="small">
                {{ getStatusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="180">
            <template #default="{ row }">
              <span class="text-gray-500 text-sm">{{ formatDate(row.created_at) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="320" fixed="right">
            <template #default="{ row }">
              <el-button size="small" type="success" @click="handlePublish(row)" class="mr-1" :disabled="row.status === 1 || row.status === 2">
                发布
              </el-button>
              <el-button size="small" type="primary" @click="handleEdit(row)" class="mr-1">
                编辑
              </el-button>
              <el-button size="small" type="info" @click="handleView(row)" class="mr-1">
                查看
              </el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="mt-6 flex justify-end">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 发布对话框 -->
    <el-dialog
      v-model="publishDialogVisible"
      title="发布内容"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="publishForm" label-position="top">
        <el-form-item label="发布方式">
          <el-radio-group v-model="publishForm.publishType">
            <el-radio value="now">立即发布</el-radio>
            <el-radio value="schedule">定时发布</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="publishForm.publishType === 'schedule'" label="发布时间" required>
          <el-date-picker
            v-model="publishForm.publishTime"
            type="datetime"
            placeholder="选择发布时间"
            style="width: 100%"
            :disabled-date="disabledDate"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="publishDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="publishLoading" @click="handleConfirmPublish" class="bg-xiaohongshu-red border-none">
            确认发布
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 查看详情对话框 -->
    <el-dialog
      v-model="viewDialogVisible"
      title="内容详情"
      width="700px"
      :close-on-click-modal="false"
    >
      <div v-if="currentContent" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">标题</label>
          <p class="text-gray-800">{{ currentContent.title }}</p>
        </div>
        <div v-if="currentContent.title_options && parseTitleOptions(currentContent.title_options).length > 0">
          <label class="block text-sm font-medium text-gray-700 mb-1">备选标题</label>
          <div class="space-y-1">
            <el-tag
              v-for="(title, index) in parseTitleOptions(currentContent.title_options)"
              :key="index"
              :type="index === currentContent.selected_title_index ? 'success' : 'info'"
              size="small"
            >
              {{ index === currentContent.selected_title_index ? '✓ ' : '' }}{{ title }}
            </el-tag>
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">内容</label>
          <div class="bg-gray-50 rounded-lg p-4 text-gray-800 whitespace-pre-wrap">
            {{ currentContent.description }}
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">标签</label>
          <div class="flex flex-wrap gap-1">
            <el-tag
              v-for="(tag, index) in parseTags(currentContent.tags)"
              :key="index"
              size="small"
              type="info"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">状态</label>
          <el-tag :type="getStatusTag(currentContent.status)" size="small">
            {{ getStatusLabel(currentContent.status) }}
          </el-tag>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">创建时间</label>
            <p class="text-gray-500 text-sm">{{ formatDate(currentContent.created_at) }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">更新时间</label>
            <p class="text-gray-500 text-sm">{{ formatDate(currentContent.updated_at) }}</p>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="viewDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="handleEdit(currentContent)">编辑</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Document } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getContentList, deleteContent, getContent, type Content } from '@/api/index'
import { publishNow, schedulePublish } from '@/api/publish'

const router = useRouter()
const loading = ref(false)
const viewDialogVisible = ref(false)
const currentContent = ref<Content | null>(null)
const publishDialogVisible = ref(false)
const publishLoading = ref(false)
const publishForm = reactive({
  publishType: 'now', // 'now' 或 'schedule'
  publishTime: ''
})

const searchForm = reactive({
  keyword: '',
  status: '' as string | number
})

const tableData = ref<Content[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 获取状态标签文本
const getStatusLabel = (status: number) => {
  const map: Record<number, string> = {
    0: '草稿',
    1: '待发布',
    2: '已发布',
    3: '发布失败'
  }
  return map[status] || '未知'
}

// 获取状态标签类型
const getStatusTag = (status: number) => {
  const map: Record<number, string> = {
    0: 'info',
    1: 'warning',
    2: 'success',
    3: 'danger'
  }
  return map[status] || ''
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

// 加载内容列表
const loadContentList = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    
    if (searchForm.status !== '' && searchForm.status !== null) {
      params.status = Number(searchForm.status)
    }

    const res = await getContentList(params)
    
    // 处理关键词过滤（后端暂不支持，先在前端过滤）
    let list = res.data.list
    if (searchForm.keyword) {
      const keyword = searchForm.keyword.toLowerCase()
      list = list.filter((item: Content) => 
        item.title.toLowerCase().includes(keyword) || 
        item.description.toLowerCase().includes(keyword)
      )
    }
    
    tableData.value = list
    pagination.total = res.data.total
  } catch (error) {
    console.error('加载内容列表失败:', error)
    ElMessage.error('加载内容列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadContentList()
}

// 重置搜索条件
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  pagination.page = 1
  loadContentList()
}

// 编辑内容
const handleEdit = (row: Content) => {
  // 跳转到创作中心并传递内容ID
  router.push({
    name: 'CreationCenter',
    query: { contentId: row.id }
  })
}

// 查看内容详情
const handleView = async (row: Content) => {
  try {
    const res = await getContent(row.id)
    currentContent.value = res.data
    viewDialogVisible.value = true
  } catch (error) {
    console.error('获取内容详情失败:', error)
    ElMessage.error('获取内容详情失败')
  }
}

// 处理发布
const handlePublish = (row: Content) => {
  currentContent.value = row
  publishForm.publishType = 'now'
  publishForm.publishTime = ''
  publishDialogVisible.value = true
}

// 禁用过去的日期
const disabledDate = (time: Date) => {
  return time.getTime() < Date.now() - 8.64e7
}

// 确认发布
const handleConfirmPublish = async () => {
  if (!currentContent.value) return

  if (publishForm.publishType === 'schedule' && !publishForm.publishTime) {
    ElMessage.warning('请选择发布时间')
    return
  }

  publishLoading.value = true
  try {
    if (publishForm.publishType === 'now') {
      await publishNow({ content_id: currentContent.value.id })
      ElMessage.success('发布成功！')
    } else {
      await schedulePublish({
        content_id: currentContent.value.id,
        publish_time: publishForm.publishTime
      })
      ElMessage.success('定时发布设置成功！')
    }
    publishDialogVisible.value = false
    loadContentList()
  } catch (error) {
    console.error('发布失败:', error)
    ElMessage.error('发布失败，请稍后重试')
  } finally {
    publishLoading.value = false
  }
}

// 删除内容
const handleDelete = (row: Content) => {
  ElMessageBox.confirm('确定要删除该内容吗？删除后无法恢复！', '删除警告', {
    confirmButtonText: '确定删除',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteContent(row.id)
      ElMessage.success('内容删除成功')
      loadContentList()
    } catch (error) {
      console.error('删除内容失败:', error)
      ElMessage.error('删除内容失败')
    }
  }).catch(() => {
  })
}

// 分页大小改变
const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadContentList()
}

// 页码改变
const handlePageChange = (page: number) => {
  pagination.page = page
  loadContentList()
}

onMounted(() => {
  loadContentList()
})
</script>

<style scoped>
.xiaohongshu-table :deep(.el-table__header-wrapper th) {
  background-color: #f8f8f8;
  font-weight: 600;
  color: #333;
}

.xiaohongshu-table :deep(.el-table__row:hover td) {
  background-color: #fff5f6 !important;
}
</style>
