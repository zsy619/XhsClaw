<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-xiaohongshu-red/10 via-xiaohongshu-pink/5 to-white p-4 relative overflow-hidden">
    <!-- 背景装饰元素 -->
    <div class="absolute top-0 left-0 w-96 h-96 bg-xiaohongshu-red/5 rounded-full -translate-x-1/2 -translate-y-1/2 animate-blob"></div>
    <div class="absolute bottom-0 right-0 w-96 h-96 bg-xiaohongshu-pink/5 rounded-full translate-x-1/2 translate-y-1/2 animate-blob animation-delay-2000"></div>
    
    <div class="w-full max-w-md relative z-10">
      <!-- Logo和标题 -->
      <div class="text-center mb-8 animate-fade-in-down">
        <div class="w-20 h-20 mx-auto mb-4 bg-gradient-to-br from-xiaohongshu-red to-xiaohongshu-pink rounded-2xl flex items-center justify-center shadow-xiaohongshu hover:scale-105 transition-transform duration-300">
          <svg class="w-10 h-10 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5z" />
            <path d="M2 17l10 5 10-5M2 12l10 5 10-5" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-gray-800">小红书内容生成与管理系统</h1>
        <p class="text-gray-500 mt-2">欢迎回来，登录您的账号</p>
      </div>

      <!-- 登录卡片 -->
      <div class="bg-white rounded-3xl shadow-xiaohongshu-xl p-8 animate-fade-in-up">
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          class="space-y-6"
          @keyup.enter="handleLogin"
        >
          <!-- 用户名 -->
          <el-form-item prop="username" class="mb-0">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名"
              size="large"
              clearable
              class="w-full"
            >
              <template #prefix>
                <User class="text-gray-400 w-5 h-5" />
              </template>
            </el-input>
          </el-form-item>

          <!-- 密码 -->
          <el-form-item prop="password" class="mb-0">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              show-password
              class="w-full"
            >
              <template #prefix>
                <Lock class="text-gray-400 w-5 h-5" />
              </template>
            </el-input>
          </el-form-item>

          <!-- 记住密码和忘记密码 -->
          <div class="flex items-center justify-between">
            <el-checkbox v-model="loginForm.rememberMe" class="text-sm text-gray-600">
              记住我
            </el-checkbox>
            <a href="#" class="text-sm text-xiaohongshu-red hover:underline">
              忘记密码？
            </a>
          </div>

          <!-- 登录按钮 -->
          <el-form-item class="mb-0">
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              class="w-full h-12 text-base font-medium rounded-xl bg-gradient-to-r from-xiaohongshu-red to-xiaohongshu-pink border-none hover:opacity-90 transition-all hover:shadow-lg hover:-translate-y-0.5 active:scale-95"
              @click="handleLogin"
            >
              <span v-if="!loading" class="flex items-center justify-center gap-2">
                <span>登 录</span>
              </span>
              <span v-else class="flex items-center justify-center gap-2">
                <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>登录中...</span>
              </span>
            </el-button>
          </el-form-item>
        </el-form>

        <!-- 分割线 -->
        <div class="flex items-center my-6">
          <div class="flex-1 h-px bg-gray-200"></div>
          <span class="px-4 text-sm text-gray-400">或</span>
          <div class="flex-1 h-px bg-gray-200"></div>
        </div>

        <!-- 社交媒体登录（占位） -->
        <div class="flex justify-center gap-4">
          <button class="w-12 h-12 rounded-full bg-gray-50 flex items-center justify-center hover:bg-gray-100 transition-colors">
            <svg class="w-6 h-6 text-gray-600" viewBox="0 0 24 24" fill="currentColor">
              <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
            </svg>
          </button>
          <button class="w-12 h-12 rounded-full bg-gray-50 flex items-center justify-center hover:bg-gray-100 transition-colors">
            <svg class="w-6 h-6 text-gray-600" viewBox="0 0 24 24" fill="currentColor">
              <path d="M22 12c0-5.523-4.477-10-10-10S2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.988C18.343 21.128 22 16.991 22 12z"/>
            </svg>
          </button>
          <button class="w-12 h-12 rounded-full bg-gray-50 flex items-center justify-center hover:bg-gray-100 transition-colors">
            <svg class="w-6 h-6 text-gray-600" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2.25c-5.385 0-9.75 4.365-9.75 9.75s4.365 9.75 9.75 9.75 9.75-4.365 9.75-9.75S17.385 2.25 12 2.25zm4.28 15.28a.75.75 0 001.06-1.06l-4.28-4.28a.75.75 0 00-1.28.53v3.45a.75.75 0 00.75.75h3.45a.75.75 0 00.53-1.28l-2.03-2.03 2.03-2.03z"/>
            </svg>
          </button>
        </div>

        <!-- 底部链接 -->
        <div class="text-center mt-6">
          <span class="text-gray-500">还没有账号？</span>
          <router-link
            to="/register"
            class="ml-1 text-xiaohongshu-red font-medium hover:underline transition-colors"
          >
            立即注册
          </router-link>
        </div>
      </div>

      <!-- 小贴士 -->
      <div class="mt-6 text-center animate-fade-in" style="animation-delay: 0.5s;">
        <p class="text-gray-400 text-sm">
          💡 登录即表示您同意我们的用户协议和隐私政策
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { login } from '@/api/auth'
import { useUserStore } from '@/stores/user'
import { Lock, User } from '@element-plus/icons-vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref<FormInstance>()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: '',
  rememberMe: false
})

const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' }
  ]
}

/**
 * 保存登录信息到本地存储
 */
const saveLoginInfo = () => {
  if (loginForm.rememberMe) {
    localStorage.setItem('rememberedUsername', loginForm.username)
    localStorage.setItem('rememberedPassword', loginForm.password)
    localStorage.setItem('rememberMe', 'true')
  } else {
    localStorage.removeItem('rememberedUsername')
    localStorage.removeItem('rememberedPassword')
    localStorage.removeItem('rememberMe')
  }
}

/**
 * 从本地存储恢复登录信息
 */
const restoreLoginInfo = () => {
  const rememberedUsername = localStorage.getItem('rememberedUsername')
  const rememberedPassword = localStorage.getItem('rememberedPassword')
  const rememberMe = localStorage.getItem('rememberMe')
  
  if (rememberedUsername && rememberedPassword && rememberMe === 'true') {
    loginForm.username = rememberedUsername
    loginForm.password = rememberedPassword
    loginForm.rememberMe = true
  }
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  // 使用 Promise 封装 validate 方法
  const valid = await new Promise<boolean>((resolve) => {
    loginFormRef.value!.validate((valid) => {
      resolve(valid)
    })
  })
  
  if (!valid) return
  
  loading.value = true
  
  try {
    // 使用 URLSearchParams 格式发送请求
    const params = new URLSearchParams()
    params.append('username', loginForm.username)
    params.append('password', loginForm.password)
    
    console.log('登录请求数据:', { username: loginForm.username, password: '***' })
    console.log('URLSearchParams:', params.toString())
    
    // 使用封装的 login 函数
    const res = await login(params)
    console.log('登录响应:', res)
    
    // 检查响应是否成功
    // 后端返回格式：{ code, message, data }
    if (res.code !== undefined && res.code !== 0) {
      // 有错误码，显示错误信息
      ElMessage.error(res.message || '登录失败')
      return
    }
    
    // 获取 data（真正的登录数据在 data 里面）
    const data = res.data || res
    
    // 检查是否有 access_token
    const accessToken = data.access_token
    const userInfo = data.user
    
    if (!accessToken) {
      ElMessage.error('登录响应格式错误，缺少 token')
      return
    }
    
    // 登录成功后保存记住密码信息
    saveLoginInfo()
    
    // 保存 token 和用户信息
    userStore.setToken(accessToken)
    if (userInfo) {
      userStore.setUserInfo(userInfo)
    }
    ElMessage.success('登录成功')
    
    // 延迟跳转，等待 token 保存
    setTimeout(() => {
      router.push('/dashboard')
    }, 500)
  } catch (error: any) {
    console.error('登录失败:', error)
    // 如果是 HTTP 错误，显示具体错误信息
    if (error.response) {
      const status = error.response.status
      if (status === 401) {
        ElMessage.error('用户名或密码错误')
      } else if (status === 400) {
        ElMessage.error('请求参数错误')
      } else {
        ElMessage.error(`登录失败：${status}`)
      }
    } else if (error.message) {
      // 如果是拦截器返回的错误消息
      ElMessage.error(error.message)
    } else {
      ElMessage.error('网络错误，请稍后重试')
    }
  } finally {
    loading.value = false
  }
}

// 组件挂载时恢复登录信息
onMounted(() => {
  restoreLoginInfo()
})
</script>

<style scoped>
/* 动画定义 */
@keyframes blob {
  0%, 100% {
    transform: translate(-50%, -50%) scale(1);
  }
  50% {
    transform: translate(-50%, -45%) scale(1.05);
  }
}

@keyframes fade-in-down {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

/* 动画类 */
.animate-blob {
  animation: blob 8s ease-in-out infinite;
}

.animation-delay-2000 {
  animation-delay: 2s;
}

.animate-fade-in-down {
  animation: fade-in-down 0.6s ease-out;
}

.animate-fade-in-up {
  animation: fade-in-up 0.6s ease-out 0.2s both;
}

.animate-fade-in {
  animation: fade-in 0.6s ease-out;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* 确保表单元素宽度一致 */
:deep(.el-form-item) {
  margin-bottom: 0 !important;
}

:deep(.el-input__wrapper) {
  box-shadow: 0 0 0 1px var(--el-border-color) inset !important;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color-hover)) inset !important;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--el-input-focus-border-color) inset !important;
}

:deep(.el-button) {
  padding: 0 !important;
}
</style>
