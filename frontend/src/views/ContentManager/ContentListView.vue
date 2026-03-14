<template>
  <div>
    <!-- 页面头部 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center" style="gap: 2px;">
        <el-icon class="text-primary-500"><Document /></el-icon>
        内容管理
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
          <el-table-column label="操作" width="240" fixed="right">
            <template #default="{ row }">
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
  </div>
</template>

<script setup lang="ts">
import { Document } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)

const searchForm = reactive({
  keyword: '',
  status: ''
})

const tableData = ref<any[]>([
  {
    id: 1,
    title: 'Python 效率工具推荐',
    description: '今天给大家带来超棒的Python效率工具推荐，让你的工作效率翻倍！',
    tags: '["Python","效率","工具"]',
    status: 2,
    created_at: '2026-03-13T10:30:00Z'
  },
  {
    id: 2,
    title: '美食探店指南',
    description: '发现城市里的隐藏宝藏美食，带你探索不一样的味道！',
    tags: '["美食","探店","推荐"]',
    status: 1,
    created_at: '2026-03-13T09:20:00Z'
  },
  {
    id: 3,
    title: '旅游攻略',
    description: '省钱出行的小技巧，让你的旅行更划算更开心！',
    tags: '["旅游","攻略","省钱"]',
    status: 0,
    created_at: '2026-03-12T16:45:00Z'
  }
])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 3
})

const getStatusLabel = (status: number) => {
  const map: Record<number, string> = {
    0: '草稿',
    1: '待发布',
    2: '已发布',
    3: '发布失败'
  }
  return map[status] || '未知'
}

const getStatusTag = (status: number) => {
  const map: Record<number, string> = {
    0: 'info',
    1: 'warning',
    2: 'success',
    3: 'danger'
  }
  return map[status] || ''
}

const parseTags = (tagsStr: string) => {
  try {
    const tags = JSON.parse(tagsStr)
    return Array.isArray(tags) ? tags.slice(0, 3) : []
  } catch {
    return []
  }
}

const truncateText = (text: string, maxLength: number) => {
  if (!text) return ''
  return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return dayjs(dateStr).format('YYYY-MM-DD HH:mm')
}

const handleSearch = () => {
  console.log('搜索:', searchForm)
  ElMessage.success('搜索功能已触发')
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  ElMessage.info('搜索条件已重置')
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑功能开发中')
}

const handleView = (row: any) => {
  ElMessage.info('查看功能开发中')
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm('确定要删除该内容吗？删除后无法恢复！', '删除警告', {
    confirmButtonText: '确定删除',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success('内容删除成功')
  }).catch(() => {
  })
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
}

const handlePageChange = (page: number) => {
  pagination.page = page
}

onMounted(() => {
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
