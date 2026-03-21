<template>
  <div>
    <!-- 页面头部 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-800 flex items-center" style="gap: 2px;">
        <el-icon class="text-primary-500"><Setting /></el-icon>
        系统设置
      </h1>
      <p class="text-gray-500 mt-1">管理您的账号安全和个人信息</p>
    </div>

    <!-- 主内容卡片 -->
    <div class="bg-white rounded-xl shadow-xiaohongshu overflow-hidden">
      <el-tabs v-model="activeTab" class="settings-tabs">
        <!-- 账号安全 -->
        <el-tab-pane label="账号安全" name="security">
          <div class="p-6">
            <h3 class="text-lg font-semibold text-gray-800 mb-6 flex items-center gap-2">
              <span class="w-2 h-2 bg-blue-500 rounded-full"></span>
              密码修改
            </h3>
            <el-form :model="securityForm" label-width="160px" class="max-w-2xl">
              <el-form-item label="当前密码">
                <el-input
                  v-model="securityForm.oldPassword"
                  type="password"
                  size="large"
                  placeholder="请输入当前密码"
                  show-password
                />
              </el-form-item>
              <el-form-item label="新密码">
                <el-input
                  v-model="securityForm.newPassword"
                  type="password"
                  size="large"
                  placeholder="请输入新密码（至少6位）"
                  show-password
                />
              </el-form-item>
              <el-form-item label="确认新密码">
                <el-input
                  v-model="securityForm.confirmPassword"
                  type="password"
                  size="large"
                  placeholder="请再次输入新密码确认"
                  show-password
                />
              </el-form-item>
              <el-form-item class="pt-4">
                <el-button
                  type="primary"
                  size="large"
                  @click="handleChangePassword"
                  class="px-8 bg-gradient-to-r from-xiaohongshu-red to-xiaohongshu-pink border-none hover:opacity-90"
                >
                  修改密码
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <!-- 关于系统 -->
        <el-tab-pane label="关于系统" name="about">
          <div class="p-6">
            <div class="max-w-2xl">
              <div class="text-center mb-8">
                <div class="w-24 h-24 mx-auto mb-4 bg-gradient-to-br from-xiaohongshu-red to-xiaohongshu-pink rounded-2xl flex items-center justify-center shadow-xiaohongshu">
                  <svg class="w-12 h-12 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 2L2 7l10 5 10-5-10-5z" />
                    <path d="M2 17l10 5 10-5M2 12l10 5 10-5" />
                  </svg>
                </div>
                <h3 class="text-2xl font-bold text-gray-800 mb-2">小红书内容生成与管理系统</h3>
                <p class="text-gray-500">高效创作，轻松管理</p>
              </div>

              <div class="bg-gray-50 rounded-xl p-6 space-y-4">
                <div class="flex justify-between items-center">
                  <span class="text-gray-600">版本号</span>
                  <span class="font-medium text-gray-800">v2.0.0</span>
                </div>
                <div class="flex justify-between items-center">
                  <span class="text-gray-600">构建时间</span>
                  <span class="font-medium text-gray-800">2026-03-14</span>
                </div>
                <div class="pt-4 border-t border-gray-200">
                  <p class="text-gray-600 mb-3">技术栈</p>
                  <div class="grid grid-cols-2 gap-3">
                    <div class="bg-white rounded-lg p-3 text-center">
                      <p class="text-sm font-medium text-gray-800">前端</p>
                      <p class="text-xs text-gray-500">Vue 3 + TypeScript</p>
                    </div>
                    <div class="bg-white rounded-lg p-3 text-center">
                      <p class="text-sm font-medium text-gray-800">UI框架</p>
                      <p class="text-xs text-gray-500">Element Plus + Tailwind</p>
                    </div>
                    <div class="bg-white rounded-lg p-3 text-center">
                      <p class="text-sm font-medium text-gray-800">后端</p>
                      <p class="text-xs text-gray-500">Go + Hertz + GORM</p>
                    </div>
                    <div class="bg-white rounded-lg p-3 text-center">
                      <p class="text-sm font-medium text-gray-800">数据库</p>
                      <p class="text-xs text-gray-500">PostgreSQL</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 安全提示 -->
    <div class="mt-6 bg-blue-50 border border-blue-200 rounded-xl p-4">
      <div class="flex items-start gap-3">
        <span class="text-xl">🔒</span>
        <div>
          <h4 class="font-medium text-gray-800 mb-1">账号安全提示</h4>
          <ul class="text-sm text-gray-600 space-y-1">
            <li>• 建议使用强密码，包含字母、数字和特殊字符</li>
            <li>• 建议定期修改密码，确保账号安全</li>
            <li>• 不要在多个网站使用相同的密码</li>
            <li>• 如发现账号异常，请立即修改密码并联系管理员</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Setting } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { reactive, ref } from 'vue'

// 当前激活的标签页
const activeTab = ref('security')

// 密码修改表单
const securityForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 修改密码处理
const handleChangePassword = () => {
  // 验证当前密码
  if (!securityForm.oldPassword) {
    ElMessage.error('请输入当前密码')
    return
  }

  // 验证新密码长度
  if (!securityForm.newPassword || securityForm.newPassword.length < 6) {
    ElMessage.error('新密码长度不能少于6个字符')
    return
  }

  // 验证两次密码一致
  if (securityForm.newPassword !== securityForm.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }

  // 提示成功（实际项目中需要调用后端API）
  ElMessage.success({
    message: '密码修改成功',
    type: 'success',
    duration: 3000,
    showClose: true
  })

  // 清空表单
  securityForm.oldPassword = ''
  securityForm.newPassword = ''
  securityForm.confirmPassword = ''
}
</script>

<style scoped>
.settings-tabs :deep(.el-tabs__header) {
  background-color: #f8f8f8;
  margin: 0;
  padding: 0 20px;
}

.settings-tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.settings-tabs :deep(.el-tabs__item) {
  padding: 0 24px;
  height: 56px;
  line-height: 56px;
  font-size: 15px;
}

.settings-tabs :deep(.el-tabs__item.is-active) {
  color: #ff2442;
  font-weight: 600;
}

.settings-tabs :deep(.el-tabs__active-bar) {
  background-color: #ff2442;
  height: 3px;
}
</style>
