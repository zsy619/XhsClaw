<template>
  <div class="main-layout flex h-screen w-full overflow-hidden bg-xiaohongshu-bg">
    <!-- 移动端遮罩层 -->
    <div 
      v-if="mobileSidebarOpen" 
      class="fixed inset-0 bg-black/50 z-40 md:hidden"
      @click="mobileSidebarOpen = false"
    ></div>

    <!-- 侧边栏 - 响应式设计 -->
    <aside 
      class="sidebar fixed left-0 top-0 z-50 h-full w-56 bg-white shadow-xiaohongshu overflow-y-auto transition-transform duration-300 md:translate-x-0 md:block"
      :class="mobileSidebarOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <!-- Logo区域 -->
      <div class="flex h-16 items-center justify-center gap-3 border-b border-gray-100 px-4">
        <img src="/vite.svg" alt="Logo" class="h-8 w-8" />
        <span class="text-lg font-bold text-xiaohongshu-dark">小红书内容生成</span>
      </div>

      <!-- 菜单区域 -->
      <nav class="mt-4 px-2">
        <template v-for="item in menuItems" :key="item.path">
          <!-- 有子菜单的项 -->
          <div v-if="item.children" class="mb-2">
            <div
              @click="toggleSubmenu(item.path)"
              class="flex items-center justify-between rounded-lg px-3 py-2.5 text-sm cursor-pointer"
              :class="isParentActive(item) ? 'bg-primary-50 text-primary-500' : 'text-gray-600 hover:bg-primary-50 hover:text-primary-500'"
            >
              <div class="flex items-center">
                <el-icon class="text-xl">
                  <component :is="item.icon" />
                </el-icon>
                <span class="ml-3">{{ item.title }}</span>
              </div>
              <el-icon class="text-sm transition-transform" :class="{ 'rotate-180': expandedMenus.includes(item.path) }">
                <ArrowDown />
              </el-icon>
            </div>
            <!-- 子菜单 -->
            <div v-show="expandedMenus.includes(item.path)" class="ml-4 mt-1 space-y-1">
              <div
                v-for="child in item.children"
                :key="child.path"
                @click="navigateTo(child.path); mobileSidebarOpen = false"
                class="flex items-center rounded-lg px-3 py-2 text-sm cursor-pointer"
                :class="route.path === child.path ? 'bg-primary-100 text-primary-600 font-medium' : 'text-gray-500 hover:bg-gray-50 hover:text-gray-700'"
              >
                <el-icon class="text-base mr-2">
                  <component :is="child.icon" />
                </el-icon>
                <span>{{ child.title }}</span>
              </div>
            </div>
          </div>
          <!-- 无子菜单的项 -->
          <div
            v-else
            @click="navigateTo(item.path); mobileSidebarOpen = false"
            class="mb-1 flex items-center rounded-lg px-3 py-2.5 text-sm cursor-pointer"
            :class="route.path === item.path ? 'bg-primary-50 text-primary-500' : 'text-gray-600 hover:bg-primary-50 hover:text-primary-500'"
          >
            <el-icon class="text-xl">
              <component :is="item.icon" />
            </el-icon>
            <span class="ml-3">{{ item.title }}</span>
          </div>
        </template>
      </nav>
    </aside>

    <!-- 主内容区 - 响应式设计 -->
    <div class="main-content flex flex-1 flex-col overflow-hidden md:ml-56">
      <!-- 顶部导航栏 -->
      <header class="flex h-16 items-center justify-between border-b border-gray-100 bg-white px-4 md:px-6">
        <div class="flex items-center gap-4">
          <!-- 移动端菜单按钮 -->
          <button 
            @click="mobileSidebarOpen = true" 
            class="md:hidden p-2 rounded-lg hover:bg-gray-100 transition-colors"
          >
            <el-icon :size="24"><Menu /></el-icon>
          </button>
          <div class="flex items-center text-sm text-gray-500">
            <span>首页</span>
            <el-icon class="mx-2 text-xs"><ArrowRight /></el-icon>
            <span class="text-xiaohongshu-dark font-medium">{{ currentPageTitle }}</span>
          </div>
        </div>

        <!-- 用户信息 -->
        <div class="flex items-center gap-3">
          <el-dropdown @command="handleCommand">
            <div class="flex cursor-pointer items-center gap-2 rounded-lg px-3 py-1.5 hover:bg-gray-100">
              <el-avatar :size="32" class="border-2 border-primary-100">
                {{ userStore.username ? userStore.username.charAt(0).toUpperCase() : 'U' }}
              </el-avatar>
              <div class="text-left">
                <div class="text-sm font-medium text-xiaohongshu-dark">{{ userStore.username || '用户' }}</div>
                <div class="text-xs text-gray-500">{{ userStore.userRole === 'admin' ? '管理员' : '普通用户' }}</div>
              </div>
              <el-icon class="text-gray-400"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="settings">
                  <el-icon class="mr-2"><Setting /></el-icon>
                  系统设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon class="mr-2"><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <!-- 内容区域 -->
      <main class="flex-1 overflow-auto p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  ArrowDown,
  ArrowRight,
  DataAnalysis,
  Edit,
  Folder,
  Setting,
  SwitchButton,
  Upload,
  List,
  Timer,
  Connection,
  DocumentCopy,
  User,
  UserFilled,
  Key,
  Menu
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { logout } from '@/api/auth'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 移动端侧边栏状态
const mobileSidebarOpen = ref(false)

// 展开的菜单
const expandedMenus = ref<string[]>(['/content', '/publish', '/users'])

// 固定的菜单项列表 - 与路由配置完全一致
const menuItems = [
  {
    path: '/dashboard',
    title: '数据概览',
    icon: DataAnalysis
  },
  {
    path: '/creation',
    title: '创作中心',
    icon: Edit
  },
  {
    path: '/content',
    title: '内容管理',
    icon: Folder,
    children: [
      {
        path: '/content/list',
        title: '我的笔记',
        icon: List
      },
      {
        path: '/content/history',
        title: '创作记录',
        icon: Timer
      }
    ]
  },
  {
    path: '/publish',
    title: '发布管理',
    icon: Upload,
    children: [
      {
        path: '/publish/test',
        title: '发布测试',
        icon: Connection
      },
      {
        path: '/publish/history',
        title: '发布历史',
        icon: DocumentCopy
      }
    ]
  },
  {
    path: '/users',
    title: '用户管理',
    icon: User,
    children: [
      {
        path: '/users/list',
        title: '用户列表',
        icon: UserFilled
      },
      {
        path: '/users/roles',
        title: '权限设置',
        icon: Key
      }
    ]
  },
  {
    path: '/settings',
    title: '系统设置',
    icon: Setting
  }
]

// 切换子菜单展开/收起
const toggleSubmenu = (path: string) => {
  const index = expandedMenus.value.indexOf(path)
  if (index > -1) {
    expandedMenus.value.splice(index, 1)
  } else {
    expandedMenus.value.push(path)
  }
}

// 判断父菜单是否激活
const isParentActive = (item: any) => {
  if (!item.children) return route.path === item.path
  return item.children.some((child: any) => route.path === child.path)
}

// 获取当前页面标题
const currentPageTitle = computed(() => {
  // 先检查是否是子菜单
  for (const item of menuItems) {
    if (item.children) {
      const child = item.children.find((c: any) => c.path === route.path)
      if (child) {
        return child.title
      }
    }
    if (item.path === route.path) {
      return item.title
    }
  }
  return '首页'
})

// 导航到指定路径
const navigateTo = (path: string) => {
  console.log('导航到:', path)
  router.push(path)
}

// 处理下拉菜单
const handleCommand = async (command: string) => {
  switch (command) {
    case 'settings':
      router.push('/settings')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        })
        // 调用后端退出登录接口
        await logout()
        // 清除本地用户信息和 Token
        userStore.clearUser()
        // 跳转到登录页
        router.push('/login')
        ElMessage.success('已退出登录')
      } catch (error) {
        // 如果是用户取消退出，不做任何处理
        // 如果是接口调用失败，仍然清除本地信息并跳转
        if (error !== 'cancel') {
          userStore.clearUser()
          router.push('/login')
          ElMessage.success('已退出登录')
        }
      }
      break
  }
}
</script>

<style scoped lang="scss">
.main-layout {
  .sidebar {
    // 滚动条样式
    &::-webkit-scrollbar {
      width: 4px;
    }
    &::-webkit-scrollbar-thumb {
      background: #e0e0e0;
      border-radius: 2px;
    }
  }
}
</style>
