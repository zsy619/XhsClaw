<template>
  <div>
    <!-- 页面头部 -->
    <div class="mb-6">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-800 flex items-center gap-2">
            <el-icon class="text-primary-500"><Connection /></el-icon>
            小红书配置
          </h1>
          <p class="mt-1 text-sm text-gray-500">管理您的小红书账号配置</p>
        </div>
        <el-button type="primary" @click="openDialog('create')">
          <el-icon><Plus /></el-icon>
          新增配置
        </el-button>
      </div>
    </div>

    <!-- 状态筛选 -->
    <el-card class="mb-6">
      <el-form :inline="true">
        <el-form-item label="状态">
          <el-select v-model="searchStatus" placeholder="全部" clearable @change="loadData">
            <el-option label="全部" value="" />
            <el-option label="正常" value="active" />
            <el-option label="待验证" value="pending" />
            <el-option label="已过期" value="expired" />
            <el-option label="异常" value="error" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 配置列表 -->
    <el-card>
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="name" label="配置名称" min-width="150" />
        <el-table-column prop="xhs_user_id" label="小红书用户ID" min-width="150">
          <template #default="{ row }">
            <span v-if="row.xhs_user_id">{{ row.xhs_user_id }}</span>
            <span v-else class="text-gray-400">未设置</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_default" label="默认" width="80" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_default" type="success" size="small">默认</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_enabled" label="启用" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_enabled ? 'success' : 'danger'" size="small">
              {{ row.is_enabled ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_login_at" label="最后验证时间" min-width="160">
          <template #default="{ row }">
            <span v-if="row.last_login_at">{{ formatTime(row.last_login_at) }}</span>
            <span v-else class="text-gray-400">未验证</span>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDialog('edit', row)">编辑</el-button>
            <el-button link type="warning" @click="handleVerify(row)">验证</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="mt-4 flex justify-end">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogMode === 'create' ? '新增小红书配置' : '编辑小红书配置'"
      width="600px"
      @close="resetForm"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="110px">
        <el-form-item label="配置名称" prop="name">
          <el-input v-model="form.name" placeholder="如：我的小红书账号" />
        </el-form-item>
        <el-form-item label="Cookie" prop="cookie">
          <el-input
            v-model="form.cookie"
            type="textarea"
            :rows="3"
            placeholder="粘贴小红书Cookie"
          />
        </el-form-item>
        <el-form-item label="小红书用户ID">
          <el-input v-model="form.xhs_user_id" placeholder="小红书用户ID（可选）" />
        </el-form-item>
        <el-form-item label="Token">
          <el-input v-model="form.token" type="password" show-password placeholder="Token（可选）" />
        </el-form-item>
        <el-form-item label="设备ID">
          <el-input v-model="form.device_id" placeholder="设备ID（可选）" />
        </el-form-item>
        <el-form-item label="设为默认">
          <el-switch v-model="form.is_default" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.is_enabled" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Connection } from '@element-plus/icons-vue'
import {
  getXHSConfigs,
  createXHSConfig,
  updateXHSConfig,
  deleteXHSConfig,
  verifyXHSConfig,
  type XHSConfig,
  type XHSConfigRequest
} from '@/api/xhs'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit'>('create')
const currentId = ref<number | null>(null)
const searchStatus = ref('')

const tableData = ref<XHSConfig[]>([])
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const formRef = ref()
const form = reactive<XHSConfigRequest>({
  name: '',
  cookie: '',
  xhs_user_id: '',
  token: '',
  device_id: '',
  is_default: false,
  is_enabled: true,
  description: ''
})

const rules = {
  name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
  cookie: [{ required: true, message: '请输入Cookie', trigger: 'blur' }]
}

const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    active: 'success',
    pending: 'warning',
    expired: 'info',
    error: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '正常',
    pending: '待验证',
    expired: '已过期',
    error: '异常'
  }
  return map[status] || status
}

const formatTime = (time: string) => {
  return new Date(time).toLocaleString('zh-CN')
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getXHSConfigs({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res.code === 0) {
      let data = res.data.items || []
      if (searchStatus.value) {
        data = data.filter((item: XHSConfig) => item.status === searchStatus.value)
      }
      tableData.value = data
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

const openDialog = (mode: 'create' | 'edit', row?: XHSConfig) => {
  dialogMode.value = mode
  if (mode === 'edit' && row) {
    currentId.value = row.id
    Object.assign(form, {
      name: row.name,
      cookie: '',
      xhs_user_id: row.xhs_user_id || '',
      token: '',
      device_id: '',
      is_default: row.is_default,
      is_enabled: row.is_enabled,
      description: row.description
    })
  }
  dialogVisible.value = true
}

const resetForm = () => {
  formRef.value?.resetFields()
  Object.assign(form, {
    name: '',
    cookie: '',
    xhs_user_id: '',
    token: '',
    device_id: '',
    is_default: false,
    is_enabled: true,
    description: ''
  })
  currentId.value = null
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    submitting.value = true
    if (dialogMode.value === 'create') {
      await createXHSConfig(form)
      ElMessage.success('创建成功')
    } else {
      await updateXHSConfig(currentId.value!, form)
      ElMessage.success('更新成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    submitting.value = false
  }
}

const handleVerify = async (row: XHSConfig) => {
  try {
    ElMessage.info('正在验证配置，请稍候...')
    const res = await verifyXHSConfig(row.id)
    if (res.code === 0 && res.data?.success) {
      ElMessage.success('验证成功！Cookie有效')
      if (res.data.user_id) {
        ElMessage.info(`用户ID: ${res.data.user_id}`)
      }
    } else {
      ElMessage.error(res.data?.message || '验证失败')
    }
    loadData()
  } catch (error) {
    console.error('验证失败:', error)
    ElMessage.error('验证失败，请检查Cookie是否正确')
  }
}

const handleDelete = async (row: XHSConfig) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除配置"${row.name}"吗？`,
      '删除确认',
      { type: 'warning' }
    )
    await deleteXHSConfig(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.xhs-config-view {
  padding: 20px;
}
</style>
