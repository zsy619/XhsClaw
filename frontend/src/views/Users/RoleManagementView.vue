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
              <el-button type="primary" size="small">新增角色</el-button>
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
              <el-icon><component :is="role.icon" /></el-icon>
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
            </div>
          </template>

          <div v-if="currentRole" class="permission-config">
            <el-form label-width="120px">
              <el-form-item label="角色名称">
                <el-input v-model="currentRole.name" disabled />
              </el-form-item>
              <el-form-item label="角色描述">
                <el-input
                  v-model="currentRole.description"
                  type="textarea"
                  :rows="3"
                />
              </el-form-item>
              <el-form-item label="权限列表">
                <el-checkbox-group v-model="currentRole.permissions">
                  <el-checkbox label="content:create">内容创建</el-checkbox>
                  <el-checkbox label="content:edit">内容编辑</el-checkbox>
                  <el-checkbox label="content:delete">内容删除</el-checkbox>
                  <el-checkbox label="content:publish">内容发布</el-checkbox>
                  <el-checkbox label="content:audit">内容审核</el-checkbox>
                  <el-checkbox label="data:view">数据查看</el-checkbox>
                  <el-checkbox label="user:manage">用户管理</el-checkbox>
                  <el-checkbox label="system:config">系统配置</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item>
                <el-button type="primary">保存配置</el-button>
                <el-button>重置</el-button>
              </el-form-item>
            </el-form>
          </div>

          <el-empty v-else description="请从左侧选择一个角色进行配置" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Key } from '@element-plus/icons-vue'

// 当前选中的角色ID
const activeRole = ref('1')

// 角色列表
const roleList = ref([
  {
    id: 1,
    name: '超级管理员',
    icon: 'UserFilled',
    description: '拥有系统所有权限',
    permissions: [
      'content:create',
      'content:edit',
      'content:delete',
      'content:publish',
      'content:audit',
      'data:view',
      'user:manage',
      'system:config'
    ]
  },
  {
    id: 2,
    name: '内容管理员',
    icon: 'EditPen',
    description: '负责内容的创作、编辑和发布',
    permissions: [
      'content:create',
      'content:edit',
      'content:delete',
      'content:publish',
      'data:view'
    ]
  },
  {
    id: 3,
    name: '普通用户',
    icon: 'User',
    description: '普通用户权限，只能查看和创建自己的内容',
    permissions: ['content:create', 'content:edit', 'data:view']
  }
])

// 当前角色
const currentRole = ref(roleList.value[0])

// 处理角色选择
const handleRoleSelect = (index: string) => {
  const role = roleList.value.find(r => r.id === Number(index))
  if (role) {
    currentRole.value = role
  }
}
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
  }
}
</style>
