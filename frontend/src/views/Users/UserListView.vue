<template>
  <div>
    <!-- 页面头部 -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center" style="gap: 2px;">
          <el-icon class="text-primary-500"><User /></el-icon>
          用户列表
        </h1>
        <p class="text-gray-500 mt-1">管理系统用户和角色权限</p>
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
            <el-select v-model="searchForm.role_id" placeholder="全部角色" clearable style="width: 160px">
              <el-option
                v-for="role in roleList"
                :key="role.id"
                :label="role.name"
                :value="role.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部状态" clearable style="width: 140px">
              <el-option label="启用" :value="1" />
              <el-option label="禁用" :value="0" />
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
          v-loading="loading"
          style="width: 100%"
          class="xiaohongshu-table"
        >
          <el-table-column prop="id" label="用户ID" width="80" />
          <el-table-column prop="username" label="用户名" width="140" />
          <el-table-column prop="nickname" label="昵称" width="140" />
          <el-table-column prop="email" label="邮箱" min-width="200" />
          <el-table-column prop="role" label="角色" width="140">
            <template #default="{ row }">
              <el-tag :type="getRoleTagType(row.role?.code)" size="small">
                {{ row.role?.name || '-' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag
                :type="row.status === 1 ? 'success' : 'danger'"
                size="small"
              >
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180">
            <template #default="{ row }">
              <span class="text-gray-500 text-sm">{{ formatDate(row.created_at) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="280" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" size="small" @click="handleEditRole(row)">
                设置角色
              </el-button>
              <el-button
                v-if="row.status === 1"
                link type="warning"
                size="small"
                @click="handleToggleStatus(row, 0)"
              >
                禁用
              </el-button>
              <el-button
                v-else
                link type="success"
                size="small"
                @click="handleToggleStatus(row, 1)"
              >
                启用
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

    <!-- 设置角色对话框 -->
    <el-dialog
      v-model="roleDialogVisible"
      title="设置用户角色"
      width="500px"
    >
      <el-form :model="roleForm" label-width="100px">
        <el-form-item label="用户名">
          <span>{{ currentUser?.username }}</span>
        </el-form-item>
        <el-form-item label="当前角色">
          <el-tag :type="getRoleTagType(currentUser?.role?.code)" size="small">
            {{ currentUser?.role?.name || '-' }}
          </el-tag>
        </el-form-item>
        <el-form-item label="选择角色">
          <el-select v-model="roleForm.role_id" placeholder="请选择角色" style="width: 100%">
            <el-option
              v-for="role in roleList"
              :key="role.id"
              :label="role.name"
              :value="role.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveRole" :loading="saveRoleLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Plus, User } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { getUserList, getAllRoles, updateUserRole, updateUserStatus, type User, type Role } from '@/api/index'

const loading = ref(false)
const saveRoleLoading = ref(false)
const roleDialogVisible = ref(false)
const currentUser = ref<User | null>(null)

const searchForm = reactive({
  keyword: '',
  role_id: '',
  status: '' as string | number
})

const userList = ref<User[]>([])
const roleList = ref<Role[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const roleForm = reactive({
  role_id: 0
})

// 获取角色标签类型
const getRoleTagType = (roleCode?: string) => {
  const typeMap: Record<string, any> = {
    'super_admin': 'danger',
    'content_manager': 'warning',
    'user': 'info'
  }
  return typeMap[roleCode || ''] || 'info'
}

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return dayjs(dateStr).format('YYYY-MM-DD HH:mm')
}

// 加载用户列表
const loadUserList = async () => {
  loading.value = true
  try {
    const res = await getUserList({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    
    let list = res.data.list
    // 前端过滤
    if (searchForm.keyword) {
      const keyword = searchForm.keyword.toLowerCase()
      list = list.filter((item: User) => 
        item.username.toLowerCase().includes(keyword) || 
        (item.email && item.email.toLowerCase().includes(keyword))
      )
    }
    if (searchForm.role_id) {
      list = list.filter((item: User) => item.role_id === Number(searchForm.role_id))
    }
    if (searchForm.status !== '' && searchForm.status !== null) {
      list = list.filter((item: User) => item.status === Number(searchForm.status))
    }
    
    userList.value = list
    pagination.total = res.data.total
  } catch (error) {
    console.error('加载用户列表失败:', error)
    ElMessage.error('加载用户列表失败')
  } finally {
    loading.value = false
  }
}

// 加载角色列表
const loadRoleList = async () => {
  try {
    const res = await getAllRoles()
    roleList.value = res.data
  } catch (error) {
    console.error('加载角色列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadUserList()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.role_id = ''
  searchForm.status = ''
  pagination.page = 1
  loadUserList()
}

// 编辑角色
const handleEditRole = (row: User) => {
  currentUser.value = row
  roleForm.role_id = row.role_id
  roleDialogVisible.value = true
}

// 保存角色
const handleSaveRole = async () => {
  if (!currentUser.value || !roleForm.role_id) {
    ElMessage.warning('请选择角色')
    return
  }
  
  saveRoleLoading.value = true
  try {
    await updateUserRole(currentUser.value.id, {
      role_id: roleForm.role_id
    })
    ElMessage.success('角色设置成功')
    roleDialogVisible.value = false
    loadUserList()
  } catch (error) {
    console.error('设置角色失败:', error)
    ElMessage.error('设置角色失败')
  } finally {
    saveRoleLoading.value = false
  }
}

// 切换用户状态
const handleToggleStatus = (row: User, status: number) => {
  const statusText = status === 1 ? '启用' : '禁用'
  ElMessageBox.confirm(`确定要${statusText}该用户吗？`, '确认操作', {
    confirmButtonText: `确定${statusText}`,
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await updateUserStatus(row.id, { status })
      ElMessage.success(`用户${statusText}成功`)
      loadUserList()
    } catch (error) {
      console.error('操作失败:', error)
      ElMessage.error('操作失败')
    }
  }).catch(() => {
  })
}

// 分页大小改变
const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadUserList()
}

// 页码改变
const handlePageChange = (page: number) => {
  pagination.page = page
  loadUserList()
}

onMounted(() => {
  loadUserList()
  loadRoleList()
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
