import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Auth/LoginView.vue'),
    meta: { title: '登录', requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Auth/RegisterView.vue'),
    meta: { title: '注册', requiresAuth: false }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/components/layout/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard/IndexView.vue'),
        meta: { title: '数据概览', icon: 'DataAnalysis' }
      },
      {
        path: 'creation',
        name: 'Creation',
        component: () => import('@/views/ContentGenerator/CreationCenterView.vue'),
        meta: { title: '创作中心', icon: 'Edit', requiresAuth: false }
      },
      {
        path: 'content',
        name: 'Content',
        redirect: '/content/list',
        meta: { title: '内容管理', icon: 'Folder' },
        children: [
          {
            path: 'list',
            name: 'ContentList',
            component: () => import('@/views/ContentManager/ContentListView.vue'),
            meta: { title: '我的笔记', icon: 'List' }
          },
          {
            path: 'history',
            name: 'ContentHistory',
            component: () => import('@/views/ContentManager/HistoryView.vue'),
            meta: { title: '创作记录', icon: 'Timer' }
          }
        ]
      },
      {
        path: 'publish',
        name: 'Publish',
        redirect: '/publish/test',
        meta: { title: '发布管理', icon: 'Upload' },
        children: [
          {
            path: 'test',
            name: 'PublishTest',
            component: () => import('@/views/Publish/PublishTestView.vue'),
            meta: { title: '发布测试', icon: 'Connection' }
          },
          {
            path: 'history',
            name: 'PublishHistory',
            component: () => import('@/views/Publish/PublishHistoryView.vue'),
            meta: { title: '发布历史', icon: 'DocumentCopy' }
          }
        ]
      },
      {
        path: 'users',
        name: 'Users',
        redirect: '/users/list',
        meta: { title: '用户管理', icon: 'User' },
        children: [
          {
            path: 'list',
            name: 'UserList',
            component: () => import('@/views/Users/UserListView.vue'),
            meta: { title: '用户列表', icon: 'UserFilled' }
          },
          {
            path: 'roles',
            name: 'RoleManagement',
            component: () => import('@/views/Users/RoleManagementView.vue'),
            meta: { title: '权限设置', icon: 'Key' }
          }
        ]
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/Settings/IndexView.vue'),
        meta: { title: '系统设置', icon: 'Setting' }
      },
      {
        path: 'token-usage',
        name: 'TokenUsage',
        component: () => import('@/views/TokenUsage/TokenUsageView.vue'),
        meta: { title: 'Token使用统计', icon: 'Coin' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/components/common/NotFound.vue'),
    meta: { title: '404' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 小红书内容生成与管理系统`
  }
  
  // 检查是否需要登录
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth !== false && !token && to.path !== '/login') {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
