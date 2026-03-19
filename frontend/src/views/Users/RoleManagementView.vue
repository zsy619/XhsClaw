<template>
  <div class="role-management-view">
    <!-- 页面头部 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center" style="gap: 2px;">
        <el-icon class="text-primary-500"><Key /></el-icon>
        权限设置
      </h1>
      <p class="text-gray-500 mt-1">管理系统角色和权限配置</p>
    </div>

    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>角色列表</span>
              <el-button type="primary" size="small" @click="showCreateDialog">新增角色</el-button>
            </div>
          </template>

          <el-menu
            :default-active="activeRole"
            class="role-menu"
            @select="handleRoleSelect"
          >
            <el-menu-item
              v-for="role in roleList"
              :key="role.id"
              :index="String(role.id)"
            >
              <el-tag v-if="role.is_system" size="small" type="info" class="mr-2">系统</el-tag>
              <span>{{ role.name }}</span>
            </el-menu-item>
          </el-menu>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>权限配置 - {{ currentRole?.name || '请选择角色' }}</span>
              <div class="flex gap-2">
                <el-button size="small" @click="handleReset">重置</el-button>
                <el-button type="primary" size="small" @click="handleSave" :loading="saveLoading">保存配置</el-button>
                <el-button
                  v-if="currentRole && !currentRole.is_system"
                  type="danger"
                  size="small"
                  @click="handleDelete"
                >
                  删除角色
                </el-button>
              </div>
            </div>
          </template>

          <div v-if="currentRole" class="permission-config">
            <el-form label-width="120px">
              <el-form-item label="角色名称">
                <el-input v-model="currentRole.name" :disabled="currentRole.is_system" />
              </el-form-item>
              <el-form-item label="角色描述">
                <el-input
                  v-model="currentRole.description"
                  type="textarea"
                  :rows="3"
                />
              </el-form-item>
              <el-form-item label="权限列表">
                <div class="space-y-4">
                  <div v-for="module in modules" :key="module.name" class="border rounded-lg p-4">
                    <h4 class="font-medium text-gray-800 mb-3 flex items-center gap-2">
                      <el-icon><component :is="module.icon" /></el-icon>
                      {{ module.label }}
                    </h4>
                    <el-checkbox-group v-model="currentRolePermissions">
                      <div class="permission-items">
                        <el-checkbox
                          v-for="perm in getPermissionsByModule(module.name)"
                          :key="perm.code"
                          :label="perm.code"
                        >
                          {{ perm.name }}
                          <el-tooltip :content="perm.description" placement="top">
                            <el-icon class="ml-1 text-gray-400"><QuestionFilled /></el-icon>
                          </el-tooltip>
                        </el-checkbox>
                      </div>
                    </el-checkbox-group>
                  </div>
                </div>
              </el-form-item>
            </el-form>
          </div>

          <el-empty v-else description="请从左侧选择一个角色进行配置" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 创建角色对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="新增角色"
      width="500px"
    >
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="角色名称">
          <el-input v-model="createForm.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色代码">
          <el-input v-model="createForm.code" placeholder="请输入角色代码（英文）" />
        </el-form-item>
        <el-form-item label="角色描述">
          <el-input
            v-model="createForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入角色描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="createLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { createRole, deleteRole, getAllRoles, getPermissions, getRole, updateRole, type Permission, type Role } from '@/api/index'
import { DataAnalysis, Edit, Folder, Key, QuestionFilled, Setting, Upload, User } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'

const activeRole = ref('')
const loading = ref(false)
const saveLoading = ref(false)
const createLoading = ref(false)
const createDialogVisible = ref(false)

const roleList = ref<Role[]>([])
const permissionList = ref<Permission[]>([])
const currentRole = ref<Role | null>(null)

const currentRolePermissions = computed({
  get: () => {
    if (!currentRole.value) return []
    try {
      return JSON.parse(currentRole.value.permissions)
    } catch {
      return []
    }
  },
  set: (val) => {
    if (currentRole.value) {
      currentRole.value.permissions = JSON.stringify(val)
    }
  }
})

// 模块配置
const modules = [
  { name: 'dashboard', label: '数据概览', icon: DataAnalysis },
  { name: 'creation', label: '创作中心', icon: Edit },
  { name: 'content', label: '内容管理', icon: Folder },
  { name: 'publish', label: '发布管理', icon: Upload },
  { name: 'user', label: '用户管理', icon: User },
  { name: 'role', label: '权限设置', icon: Key },
  { name: 'settings', label: '系统设置', icon: Setting }
]

// 根据模块获取权限
const getPermissionsByModule = (module: string) => {
  return permissionList.value.filter(p => p.module === module)
}

// 加载角色列表
const loadRoleList = async () => {
  loading.value = true
  try {
    const res = await getAllRoles()
    roleList.value = res.data
    if (roleList.value.length > 0 && !activeRole.value) {
      activeRole.value = String(roleList.value[0].id)
      handleRoleSelect(activeRole.value)
    }
  } catch (error) {
    console.error('加载角色列表失败:', error)
    ElMessage.error('加载角色列表失败')
  } finally {
    loading.value = false
  }
}

// 加载权限列表
const loadPermissionList = async () => {
  try {
    const res = await getPermissions()
    permissionList.value = res.data
  } catch (error) {
    console.error('加载权限列表失败:', error)
  }
}

// 选择角色
const handleRoleSelect = async (index: string) => {
  activeRole.value = index
  try {
    const res = await getRole(Number(index))
    currentRole.value = res.data
  } catch (error) {
    console.error('获取角色详情失败:', error)
    ElMessage.error('获取角色详情失败')
  }
}

// 保存角色
const handleSave = async () => {
  if (!currentRole.value) return
  
  saveLoading.value = true
  try {
    await updateRole(currentRole.value.id, {
      name: currentRole.value.name,
      description: currentRole.value.description,
      permissions: currentRolePermissions.value
    })
    ElMessage.success('保存成功')
    loadRoleList()
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  } finally {
    saveLoading.value = false
  }
}

// 重置
const handleReset = () => {
  if (currentRole.value) {
    handleRoleSelect(String(currentRole.value.id))
  }
}

// 显示创建对话框
const showCreateDialog = () => {
  createForm.name = ''
  createForm.code = ''
  createForm.description = ''
  createDialogVisible.value = true
}

const createForm = ref({
  name: '',
  code: '',
  description: ''
})

// 创建角色
const handleCreate = async () => {
  if (!createForm.value.name || !createForm.value.code) {
    ElMessage.warning('请填写角色名称和代码')
    return
  }
  
  createLoading.value = true
  try {
    await createRole({
      name: createForm.value.name,
      code: createForm.value.code,
      description: createForm.value.description,
      permissions: []
    })
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    loadRoleList()
  } catch (error) {
    console.error('创建失败:', error)
    ElMessage.error('创建失败')
  } finally {
    createLoading.value = false
  }
}

// 删除角色
const handleDelete = () => {
  if (!currentRole.value) return
  
  ElMessageBox.confirm('确定要删除该角色吗？删除后无法恢复！', '删除警告', {
    confirmButtonText: '确定删除',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteRole(currentRole.value!.id)
      ElMessage.success('删除成功')
      currentRole.value = null
      activeRole.value = ''
      loadRoleList()
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {
  })
}

onMounted(() => {
  loadRoleList()
  loadPermissionList()
})
</script>

<style scoped lang="scss">
.role-management-view {
  .box-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .role-menu {
      border-right: none;
    }

    .permission-config {
      padding: 20px 0;
    }

    .permission-items {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;
      width: 100%;
      
      :deep(.el-checkbox) {
        margin-right: 0;
        margin-bottom: 8px;
        white-space: nowrap;
      }
    }
  }
}
</style>
