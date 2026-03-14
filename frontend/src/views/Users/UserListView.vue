<template>
  <div>
    <!-- 页面头部 -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center" style="gap: 2px;">
          <el-icon class="text-primary-500"><User /></el-icon>
          用户管理
        </h1>
        <p class="text-gray-500 mt-1">管理系统用户和角色权限</p>
      </div>
      <div class="mt-4 md:mt-0">
        <el-button
          type="primary"
          class="bg-gradient-to-r from-xiaohongshu-red to-xiaohongshu-pink border-none hover:opacity-90"
        >
          <el-icon class="mr-1"><Plus /></el-icon>
          新增用户
        </el-button>
      </div>
    </div>

    <!-- 主内容卡片 -->
    <div class="bg-white rounded-xl shadow-xiaohongshu p-5">
      <!-- 搜索栏 -->
      <div class="bg-gray-50 rounded-xl p-4 mb-6">
        <el-form :inline="true" :model="searchForm" class="flex flex-wrap gap-4">
          <el-form-item label="关键词">
            <el-input
              v-model="searchForm.keyword"
              placeholder="搜索用户名或邮箱"
              clearable
              style="width: 220px"
            />
          </el-form-item>
          <el-form-item label="角色">
            <el-select v-model="searchForm.role" placeholder="全部角色" clearable style="width: 160px">
              <el-option label="超级管理员" value="超级管理员" />
              <el-option label="内容管理员" value="内容管理员" />
              <el-option label="普通用户" value="普通用户" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部状态" clearable style="width: 140px">
              <el-option label="启用" value="active" />
              <el-option label="禁用" value="inactive" />
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

      <!-- 用户列表 -->
      <div class="overflow-x-auto">
        <el-table
          :data="userList"
          style="width: 100%"
          class="xiaohongshu-table"
        >
          <el-table-column prop="id" label="用户ID" width="80" />
          <el-table-column prop="username" label="用户名" width="140" />
          <el-table-column prop="email" label="邮箱" min-width="200" />
          <el-table-column prop="role" label="角色" width="140">
            <template #default="{ row }">
              <el-tag :type="getRoleTagType(row.role)" size="small">
                {{ row.role }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag
                :type="row.status === 'active' ? 'success' : 'danger'"
                size="small"
              >
                {{ row.status === 'active' ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" width="180" />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" size="small" @click="handleEdit(row)">
                编辑
              </el-button>
              <el-button link type="danger" size="small" @click="handleDelete(row)">
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

    <!-- 小贴士 -->
    <div class="mt-6 bg-blue-50 border border-blue-100 rounded-xl p-4">
      <div class="flex items-start gap-3">
        <span class="text-xl">💡</span>
        <div>
          <h4 class="font-medium text-gray-800 mb-1">用户管理说明</h4>
          <ul class="text-sm text-gray-600 space-y-1">
            <li>• 超级管理员：拥有系统所有权限，可管理用户和角色</li>
            <li>• 内容管理员：可创建、编辑、发布内容，管理内容库</li>
            <li>• 普通用户：只能创建和管理自己的内容</li>
            <li>• 禁用用户无法登录系统，所有权限被收回</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Plus, User } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { reactive, ref } from 'vue'

const searchForm = reactive({
  keyword: '',
  role: '',
  status: ''
})

const userList = ref([
  {
    id: 1,
    username: 'admin',
    email: 'admin@example.com',
    role: '超级管理员',
    status: 'active',
    createdAt: '2024-01-01 00:00:00'
  },
  {
    id: 2,
    username: 'editor',
    email: 'editor@example.com',
    role: '内容管理员',
    status: 'active',
    createdAt: '2024-01-15 10:30:00'
  },
  {
    id: 3,
    username: 'user',
    email: 'user@example.com',
    role: '普通用户',
    status: 'active',
    createdAt: '2024-02-01 14:20:00'
  }
])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 3
})

const getRoleTagType = (role: string) => {
  const typeMap: Record<string, any> = {
    '超级管理员': 'danger',
    '内容管理员': 'warning',
    '普通用户': 'info'
  }
  return typeMap[role] || 'info'
}

const handleSearch = () => {
  ElMessage.success('搜索功能已触发')
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.role = ''
  searchForm.status = ''
  ElMessage.info('搜索条件已重置')
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑功能开发中')
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm('确定要删除该用户吗？删除后无法恢复！', '删除警告', {
    confirmButtonText: '确定删除',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success('用户删除成功')
  }).catch(() => {
  })
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
}

const handlePageChange = (page: number) => {
  pagination.page = page
}
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
