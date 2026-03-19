<template>
  <div>
    <!-- 页面头部 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-800 flex items-center" style="gap: 2px;">
        <el-icon class="text-primary-500"><Setting /></el-icon>
        系统设置
      </h1>
      <p class="text-gray-500 mt-1">配置系统参数和个人偏好设置</p>
    </div>

    <!-- 主内容卡片 -->
    <div class="bg-white rounded-xl shadow-xiaohongshu overflow-hidden">
      <el-tabs v-model="activeTab" class="settings-tabs">
        <!-- 大模型配置 -->
        <el-tab-pane label="大模型配置" name="llm">
          <div class="p-6">
            <h3 class="text-lg font-semibold text-gray-800 mb-6 flex items-center gap-2">
              <span class="w-2 h-2 bg-xiaohongshu-red rounded-full"></span>
              大模型配置
            </h3>
            <div class="mb-4 p-4 bg-blue-50 rounded-lg border border-blue-100">
              <p class="text-sm text-blue-700 flex items-center gap-2">
                <el-icon><InfoFilled /></el-icon>
                大模型配置是生成文案和渲染图片的必要参数，请确保填写正确
              </p>
            </div>
            <el-form :model="llmConfig" label-width="160px" class="max-w-2xl">
              <el-form-item label="API Key" required>
                <el-input
                  v-model="llmConfig.llm_api_key"
                  type="password"
                  show-password
                  size="large"
                  placeholder="请输入您的大模型API Key"
                />
                <div class="mt-1 text-xs text-gray-500">
                  例如：sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
                </div>
              </el-form-item>
              <el-form-item label="Base URL" required>
                <el-input
                  v-model="llmConfig.llm_base_url"
                  size="large"
                  placeholder="请输入大模型API地址，例如：https://api.deepseek.com"
                />
                <div class="mt-1 text-xs text-gray-500">
                  请确保包含完整的协议和域名
                </div>
              </el-form-item>
              <el-form-item label="模型名称" required>
                <el-input
                  v-model="llmConfig.llm_model"
                  size="large"
                  placeholder="请输入模型名称，例如：deepseek-chat"
                />
                <div class="mt-1 text-xs text-gray-500">
                  请根据您使用的大模型服务提供商的要求填写
                </div>
              </el-form-item>
              <el-form-item class="pt-4">
                <el-button
                  type="primary"
                  size="large"
                  :loading="savingLLM"
                  @click="handleSaveLLMConfig"
                  class="px-8 bg-gradient-to-r from-xiaohongshu-red to-xiaohongshu-pink border-none hover:opacity-90"
                >
                  保存配置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <!-- 小红书配置 -->
        <el-tab-pane label="小红书配置" name="xiaohongshu">
          <div class="p-6">
            <h3 class="text-lg font-semibold text-gray-800 mb-6 flex items-center gap-2">
              <span class="w-2 h-2 bg-green-500 rounded-full"></span>
              小红书配置
            </h3>
            <el-form :model="xiaohongshuConfig" label-width="160px" class="max-w-2xl">
              <el-form-item label="Cookie">
                <el-input
                  v-model="xiaohongshuConfig.xiaohongshu_cookie"
                  type="textarea"
                  :rows="3"
                  size="large"
                  placeholder="请输入您的小红书Cookie"
                />
              </el-form-item>
              <el-form-item label="User ID">
                <el-input
                  v-model="xiaohongshuConfig.xiaohongshu_user_id"
                  size="large"
                  placeholder="请输入您的小红书用户ID"
                />
              </el-form-item>
              <el-form-item label="Token">
                <el-input
                  v-model="xiaohongshuConfig.xiaohongshu_token"
                  type="password"
                  show-password
                  size="large"
                  placeholder="请输入您的小红书Token"
                />
              </el-form-item>
              <el-form-item label="默认发布时间">
                <el-time-picker
                  v-model="defaultPublishTime"
                  size="large"
                  format="HH:mm"
                  value-format="HH:mm"
                  placeholder="选择默认发布时间"
                />
              </el-form-item>
              <el-form-item label="自动发布">
                <el-switch v-model="xiaohongshuConfig.auto_publish_enabled" />
                <span class="ml-2 text-sm text-gray-500">启用后，内容将在设定的时间自动发布</span>
              </el-form-item>
              <el-form-item class="pt-4">
                <el-button
                  type="primary"
                  size="large"
                  :loading="savingXiaohongshu"
                  @click="handleSaveXiaohongshuConfig"
                  class="px-8 bg-gradient-to-r from-xiaohongshu-red to-xiaohongshu-pink border-none hover:opacity-90"
                >
                  保存配置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

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
                />
              </el-form-item>
              <el-form-item label="新密码">
                <el-input
                  v-model="securityForm.newPassword"
                  type="password"
                  size="large"
                  placeholder="请输入新密码（至少6位）"
                />
              </el-form-item>
              <el-form-item label="确认新密码">
                <el-input
                  v-model="securityForm.confirmPassword"
                  type="password"
                  size="large"
                  placeholder="请再次输入新密码确认"
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

    <!-- 小贴士 -->
    <div class="mt-6 bg-xiaohongshu-red/5 border border-xiaohongshu-red/20 rounded-xl p-4">
      <div class="flex items-start gap-3">
        <span class="text-xl">💡</span>
        <div>
          <h4 class="font-medium text-gray-800 mb-1">设置小贴士</h4>
          <ul class="text-sm text-gray-600 space-y-1">
            <li>• 请妥善保管您的API Key和Cookie，不要泄露给他人</li>
            <li>• 建议定期修改密码，确保账号安全</li>
            <li>• 每个用户可以独立配置自己的大模型和小红书账号</li>
            <li>• 如有问题，请联系系统管理员获取帮助</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getUserConfig, updateUserConfig, type UserConfigRequest } from '@/api/userConfig'
import { Setting } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

const activeTab = ref('llm')
const savingLLM = ref(false)
const savingXiaohongshu = ref(false)
const defaultPublishTime = ref('')

const llmConfig = reactive<UserConfigRequest>({
  llm_api_key: '',
  llm_base_url: '',
  llm_model: ''
})

const xiaohongshuConfig = reactive<UserConfigRequest>({
  xiaohongshu_cookie: '',
  xiaohongshu_user_id: '',
  xiaohongshu_token: '',
  default_publish_time: '',
  auto_publish_enabled: false
})

const securityForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 加载用户配置
const loadUserConfig = async () => {
  try {
    const res = await getUserConfig()
    const config = res.data
    
    // 填充大模型配置
    if (config.llm_api_key) llmConfig.llm_api_key = config.llm_api_key
    if (config.llm_base_url) llmConfig.llm_base_url = config.llm_base_url
    if (config.llm_model) llmConfig.llm_model = config.llm_model
    
    // 填充小红书配置
    if (config.xiaohongshu_cookie) xiaohongshuConfig.xiaohongshu_cookie = config.xiaohongshu_cookie
    if (config.xiaohongshu_user_id) xiaohongshuConfig.xiaohongshu_user_id = config.xiaohongshu_user_id
    if (config.xiaohongshu_token) xiaohongshuConfig.xiaohongshu_token = config.xiaohongshu_token
    if (config.default_publish_time) {
      xiaohongshuConfig.default_publish_time = config.default_publish_time
      defaultPublishTime.value = config.default_publish_time
    }
    if (config.auto_publish_enabled !== undefined) {
      xiaohongshuConfig.auto_publish_enabled = config.auto_publish_enabled
    }
  } catch (error) {
    console.error('加载用户配置失败:', error)
  }
}

// 保存大模型配置
const handleSaveLLMConfig = async () => {
  // 验证配置
  if (!llmConfig.llm_api_key) {
    ElMessage.warning('请输入 API Key')
    return
  }
  if (!llmConfig.llm_base_url) {
    ElMessage.warning('请输入 Base URL')
    return
  }
  if (!llmConfig.llm_model) {
    ElMessage.warning('请输入模型名称')
    return
  }

  savingLLM.value = true
  try {
    await updateUserConfig(llmConfig)
    ElMessage.success({
      message: '大模型配置已保存',
      type: 'success',
      duration: 3000,
      showClose: true
    })
  } catch (error) {
    console.error('保存大模型配置失败:', error)
    ElMessage.error('保存失败，请稍后重试')
  } finally {
    savingLLM.value = false
  }
}

// 保存小红书配置
const handleSaveXiaohongshuConfig = async () => {
  // 验证配置
  if (!xiaohongshuConfig.xiaohongshu_cookie) {
    ElMessage.warning('请输入小红书 Cookie')
    return
  }

  savingXiaohongshu.value = true
  try {
    if (defaultPublishTime.value) {
      xiaohongshuConfig.default_publish_time = defaultPublishTime.value
    }
    await updateUserConfig(xiaohongshuConfig)
    ElMessage.success({
      message: '小红书配置已保存',
      type: 'success',
      duration: 3000,
      showClose: true
    })
  } catch (error) {
    console.error('保存小红书配置失败:', error)
    ElMessage.error('保存失败，请稍后重试')
  } finally {
    savingXiaohongshu.value = false
  }
}

const handleChangePassword = () => {
  if (!securityForm.oldPassword) {
    ElMessage.error('请输入当前密码')
    return
  }
  if (!securityForm.newPassword || securityForm.newPassword.length < 6) {
    ElMessage.error('新密码长度不能少于6个字符')
    return
  }
  if (securityForm.newPassword !== securityForm.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }
  ElMessage.success('密码修改成功')
  securityForm.oldPassword = ''
  securityForm.newPassword = ''
  securityForm.confirmPassword = ''
}

// 组件挂载时加载配置
onMounted(() => {
  loadUserConfig()
})
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
