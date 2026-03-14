<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-xiaohongshu-red/10 via-xiaohongshu-pink/5 to-white p-4">
    <div class="w-full max-w-md">
      <!-- Logo和标题 -->
      <div class="text-center mb-8">
        <div class="w-16 h-16 mx-auto mb-4 bg-gradient-to-br from-xiaohongshu-red to-xiaohongshu-pink rounded-2xl flex items-center justify-center shadow-xiaohongshu">
          <svg class="w-8 h-8 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5z" />
            <path d="M2 17l10 5 10-5M2 12l10 5 10-5" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-gray-800">注册新账号</h1>
        <p class="text-gray-500 mt-2">创建您的小红书内容管理账号</p>
      </div>

      <!-- 注册卡片 -->
      <div class="bg-white rounded-3xl shadow-xiaohongshu-xl p-8">
        <el-form
          ref="registerFormRef"
          :model="registerForm"
          :rules="registerRules"
          class="space-y-5"
        >
          <!-- 用户名 -->
          <el-form-item prop="username">
            <el-input
              v-model="registerForm.username"
              placeholder="请输入用户名（6-20位，字母数字组合）"
              size="large"
            >
              <template #prefix>
                <User class="text-gray-400 w-5 h-5" />
              </template>
            </el-input>
          </el-form-item>

          <!-- 用户昵称 -->
          <el-form-item prop="nickname">
            <el-input
              v-model="registerForm.nickname"
              placeholder="请输入用户昵称（不少于3个字符）"
              size="large"
            >
              <template #prefix>
                <User class="text-gray-400 w-5 h-5" />
              </template>
            </el-input>
          </el-form-item>

          <!-- 邮箱 -->
          <el-form-item prop="email">
            <el-input
              v-model="registerForm.email"
              placeholder="请输入邮箱地址"
              size="large"
            >
              <template #prefix>
                <Message class="text-gray-400 w-5 h-5" />
              </template>
            </el-input>
          </el-form-item>

          <!-- 密码 -->
          <el-form-item prop="password">
            <el-input
              v-model="registerForm.password"
              type="password"
              placeholder="请输入密码（8位以上，含大小写字母+数字+特殊字符）"
              size="large"
              show-password
            >
              <template #prefix>
                <Lock class="text-gray-400 w-5 h-5" />
              </template>
            </el-input>
          </el-form-item>

          <!-- 确认密码 -->
          <el-form-item prop="confirmPassword">
            <el-input
              v-model="registerForm.confirmPassword"
              type="password"
              placeholder="请再次输入密码确认"
              size="large"
              show-password
            >
              <template #prefix>
                <Lock class="text-gray-400 w-5 h-5" />
              </template>
            </el-input>
          </el-form-item>

          <!-- 注册按钮 -->
          <el-form-item class="pt-2">
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              class="w-full h-12 text-base font-medium rounded-xl bg-gradient-to-r from-xiaohongshu-red to-xiaohongshu-pink border-none hover:opacity-90 transition-all hover:shadow-lg hover:-translate-y-0.5"
              @click="handleRegister"
            >
              {{ loading ? '注册中...' : '立即注册' }}
            </el-button>
          </el-form-item>
        </el-form>

        <!-- 底部链接 -->
        <div class="text-center">
          <span class="text-gray-500">已有账号？</span>
          <router-link
            to="/login"
            class="ml-1 text-xiaohongshu-red font-medium hover:underline"
          >
            立即登录
          </router-link>
        </div>
      </div>

      <!-- 小贴士 -->
      <div class="mt-6 text-center">
        <p class="text-gray-400 text-sm">
          💡 注册即表示您同意我们的用户协议和隐私政策
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { http } from '@/api/request'
import { Lock, Message, User } from '@element-plus/icons-vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const registerFormRef = ref<FormInstance>()
const loading = ref(false)

const registerForm = reactive({
  username: '',
  nickname: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule: any, value: string, callback: any) => {
  if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

// 验证用户名：6-20位，字母数字组合
const validateUsername = (rule: any, value: string, callback: any) => {
  if (!value) {
    callback(new Error('请输入用户名'))
  } else if (value.length < 6 || value.length > 20) {
    callback(new Error('用户名长度必须在6-20位之间'))
  } else if (!/^[a-zA-Z0-9]+$/.test(value)) {
    callback(new Error('用户名只能包含字母和数字'))
  } else {
    callback()
  }
}

// 验证昵称：不少于3个字符
const validateNickname = (rule: any, value: string, callback: any) => {
  if (!value) {
    callback(new Error('请输入用户昵称'))
  } else if (value.length < 3) {
    callback(new Error('用户昵称不能少于3个字符'))
  } else if (value.length > 50) {
    callback(new Error('用户昵称不能超过50个字符'))
  } else {
    callback()
  }
}

// 验证密码：8位以上，包含大小写字母+数字+特殊字符
const validatePassword = (rule: any, value: string, callback: any) => {
  if (!value) {
    callback(new Error('请输入密码'))
  } else if (value.length < 8) {
    callback(new Error('密码长度不能少于8位'))
  } else if (!/[a-z]/.test(value)) {
    callback(new Error('密码必须包含小写字母'))
  } else if (!/[A-Z]/.test(value)) {
    callback(new Error('密码必须包含大写字母'))
  } else if (!/[0-9]/.test(value)) {
    callback(new Error('密码必须包含数字'))
  } else if (!/[!@#$%^&*(),.?":{}|<>]/.test(value)) {
    callback(new Error('密码必须包含特殊字符'))
  } else {
    callback()
  }
}

const registerRules: FormRules = {
  username: [
    { validator: validateUsername, trigger: 'blur' }
  ],
  nickname: [
    { validator: validateNickname, trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { validator: validatePassword, trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  await registerFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    
    try {
      // 直接使用 http 实例，不经过统一拦截器处理
      await http.post('/auth/register', {
        username: registerForm.username,
        nickname: registerForm.nickname,
        email: registerForm.email,
        password: registerForm.password
      }, {
        // 跳过统一响应拦截器的 code 检查
        skipInterceptor: true
      })
      
      ElMessage.success('注册成功，请登录')
      router.push('/login')
    } catch (error: any) {
      console.error('注册失败:', error)
      // 如果是 HTTP 错误，显示具体错误信息
      if (error.response) {
        const status = error.response.status
        if (status === 400) {
          ElMessage.error(error.response.data?.detail || '请求参数错误')
        } else {
          ElMessage.error(`注册失败：${status}`)
        }
      } else {
        ElMessage.error('网络错误，请稍后重试')
      }
    } finally {
      loading.value = false
    }
  })
}
</script>
