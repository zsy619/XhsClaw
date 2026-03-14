<template>
  <div>
    <!-- 页面头部 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center" style="gap: 2px;">
        <el-icon class="text-primary-500"><Timer /></el-icon>
        内容历史记录
      </h1>
      <p class="text-gray-500 mt-1">查看内容的版本历史和修改记录</p>
    </div>

    <!-- 内容卡片 -->
    <div class="bg-white rounded-xl shadow-xiaohongshu p-5">
      <!-- 空状态 -->
      <div v-if="!hasData" class="text-center py-16">
        <div class="w-24 h-24 mx-auto mb-6 bg-gray-100 rounded-full flex items-center justify-center">
          <svg class="w-12 h-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-600 mb-2">暂无历史记录</h3>
        <p class="text-gray-400">当您创建或修改内容后，历史记录会显示在这里</p>
      </div>

      <!-- 历史记录列表 -->
      <div v-else class="space-y-4">
        <!-- 示例记录项 -->
        <div
          v-for="item in historyData"
          :key="item.id"
          class="border border-gray-100 rounded-xl p-4 hover:border-xiaohongshu-red/30 hover:bg-xiaohongshu-red/5 transition-all"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <span class="px-3 py-1 rounded-full text-xs font-medium" :class="getTypeClass(item.type)">
                  {{ getTypeLabel(item.type) }}
                </span>
                <span class="text-sm text-gray-500">{{ item.created_at }}</span>
              </div>
              <h4 class="font-medium text-gray-800 mb-1">{{ item.title }}</h4>
              <p class="text-sm text-gray-500 line-clamp-2">{{ item.description }}</p>
            </div>
            <div class="flex items-center gap-2 ml-4">
              <el-button size="small" type="primary" link @click="handleView(item)">
                查看详情
              </el-button>
              <el-button size="small" type="success" link @click="handleRestore(item)">
                恢复此版本
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="hasData" class="mt-6 flex justify-end">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 小贴士 -->
    <div class="mt-6 bg-xiaohongshu-red/5 border border-xiaohongshu-red/20 rounded-xl p-4">
      <div class="flex items-start gap-3">
        <span class="text-xl">💡</span>
        <div>
          <h4 class="font-medium text-gray-800 mb-1">历史记录小贴士</h4>
          <ul class="text-sm text-gray-600 space-y-1">
            <li>• 系统会自动保存内容的每次修改记录</li>
            <li>• 您可以随时查看和恢复历史版本</li>
            <li>• 历史记录会保存原始内容、修改时间和操作人信息</li>
            <li>• 删除的内容也可以从历史记录中恢复</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Timer } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { reactive, ref } from 'vue'

const hasData = ref(false)

const historyData = ref<any[]>([
  {
    id: 1,
    type: 'edit',
    title: '修改了「Python 效率工具推荐」',
    description: '更新了推荐工具列表，新增了3个实用工具，优化了文案表达',
    created_at: '2026-03-13 14:30:00'
  },
  {
    id: 2,
    type: 'create',
    title: '创建了「美食探店指南」',
    description: '使用方案3（主题+图片）生成了新内容，包含3张配套图片',
    created_at: '2026-03-13 09:20:00'
  }
])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 2
})

const getTypeLabel = (type: string) => {
  const map: Record<string, string> = {
    create: '新建',
    edit: '修改',
    delete: '删除',
    publish: '发布'
  }
  return map[type] || type
}

const getTypeClass = (type: string) => {
  const map: Record<string, string> = {
    create: 'bg-green-100 text-green-700',
    edit: 'bg-blue-100 text-blue-700',
    delete: 'bg-red-100 text-red-700',
    publish: 'bg-purple-100 text-purple-700'
  }
  return map[type] || 'bg-gray-100 text-gray-700'
}

const handleView = (item: any) => {
  ElMessage.info('查看详情功能开发中')
}

const handleRestore = (item: any) => {
  ElMessageBox.confirm('确定要恢复到此版本吗？', '版本恢复', {
    confirmButtonText: '确定恢复',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success('版本恢复成功')
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
