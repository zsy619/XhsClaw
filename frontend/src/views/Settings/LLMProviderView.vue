<template>
  <div class="llm-provider-view">
    <div class="flex justify-between items-center mb-4">
      <div>
        <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center gap-2">
          <el-icon class="text-primary-500"><Monitor /></el-icon>
          大模型配置
        </h1>
        <p class="mt-1 text-sm text-gray-500">管理您的大模型服务商配置</p>
      </div>
      <el-button type="primary" @click="openDialog('create')">
        <el-icon><Plus /></el-icon>
        新增配置
      </el-button>
    </div>

    <!-- 搜索筛选 -->
    <el-card class="mb-4">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="服务商">
          <el-select v-model="searchForm.provider" placeholder="全部" clearable>
            <el-option
              v-for="item in providerOptions"
              :key="item.code"
              :label="item.name"
              :value="item.code"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="启用" value="enabled" />
            <el-option label="禁用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 配置列表 -->
    <el-card>
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="name" label="配置名称" min-width="120" />
        <el-table-column prop="provider" label="服务商" min-width="100">
          <template #default="{ row }">
            <el-tag type="info">{{ getProviderName(row.provider) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="base_url" label="API地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="model_name" label="模型" min-width="120" />
        <el-table-column prop="is_default" label="默认" width="80" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_default" type="success" size="small">默认</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_enabled" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_enabled ? 'success' : 'danger'" size="small">
              {{ row.is_enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDialog('edit', row)">编辑</el-button>
            <el-button link type="primary" @click="handleTest(row)">测试</el-button>
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
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogMode === 'create' ? '新增大模型配置' : '编辑大模型配置'"
      width="600px"
      @close="resetForm"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="配置名称" prop="name">
          <el-input v-model="form.name" placeholder="如：我的DeepSeek配置" />
        </el-form-item>
        <el-form-item label="服务商" prop="provider">
          <el-select v-model="form.provider" placeholder="选择服务商" style="width: 100%">
            <el-option
              v-for="item in providerOptions"
              :key="item.code"
              :label="item.name"
              :value="item.code"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="API地址" prop="base_url">
          <el-input v-model="form.base_url" placeholder="如：https://api.deepseek.com" />
        </el-form-item>
        <el-form-item label="API密钥" prop="api_key">
          <el-input v-model="form.api_key" type="password" show-password placeholder="输入API密钥" />
        </el-form-item>
        <el-form-item label="模型名称" prop="model_name">
          <el-select v-model="form.model_name" placeholder="选择或输入模型" filterable allow-create style="width: 100%">
            <el-option
              v-for="item in getModelOptions(form.provider)"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="超时时间">
          <el-input-number v-model="form.timeout" :min="10" :max="300" /> 秒
        </el-form-item>
        <el-form-item label="重试次数">
          <el-input-number v-model="form.retry_count" :min="0" :max="10" />
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
import { Plus } from '@element-plus/icons-vue'
import {
  getLLMProviders,
  createLLMProvider,
  updateLLMProvider,
  deleteLLMProvider,
  testLLMConnection,
  type LLMProvider,
  type LLMProviderRequest
} from '@/api/llm'
import { getLLMProviderDicts } from '@/api/systemDict'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit'>('create')
const currentId = ref<number | null>(null)
const providerOptions = ref<{ code: string; name: string; value: string }[]>([])

const tableData = ref<LLMProvider[]>([])
const searchForm = reactive({
  provider: '',
  status: ''
})
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const formRef = ref()
const form = reactive<LLMProviderRequest>({
  name: '',
  provider: '',
  base_url: '',
  api_key: '',
  model_name: '',
  is_default: false,
  is_enabled: true,
  timeout: 60,
  retry_count: 3,
  description: ''
})

const rules = {
  name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
  provider: [{ required: true, message: '请选择服务商', trigger: 'change' }]
}

// 服务商选项
const providerList = [
  { code: 'openai', name: 'OpenAI', models: ['gpt-4o', 'gpt-4o-mini', 'gpt-4-turbo', 'gpt-4', 'gpt-3.5-turbo'] },
  { code: 'deepseek', name: 'DeepSeek', models: ['deepseek-chat', 'deepseek-coder'] },
  { code: 'azure', name: 'Azure OpenAI', models: ['gpt-4o', 'gpt-4', 'gpt-35-turbo'] },
  { code: 'claude', name: 'Anthropic Claude', models: ['claude-3-5-sonnet', 'claude-3-opus', 'claude-3-haiku'] },
  { code: 'gemini', name: 'Google Gemini', models: ['gemini-1.5-pro', 'gemini-1.5-flash', 'gemini-pro'] },
  { code: 'qwen', name: '通义千问', models: ['qwen-turbo', 'qwen-plus', 'qwen-max'] },
  { code: 'glm', name: '智谱AI', models: ['glm-4', 'glm-4-flash', 'glm-3-turbo'] },
  { code: 'ollama', name: 'Ollama', models: [] },
  { code: 'custom', name: '自定义', models: [] }
]

const getProviderName = (code: string) => {
  return providerList.find(p => p.code === code)?.name || code
}

const getModelOptions = (provider: string) => {
  return providerList.find(p => p.code === provider)?.models || []
}

const loadProviders = async () => {
  try {
    const res = await getLLMProviderDicts()
    if (res.code === 0 && res.data) {
      providerOptions.value = res.data.map((d: any) => ({
        code: d.code,
        name: d.name,
        value: d.value
      }))
    }
  } catch (error) {
    providerOptions.value = providerList
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getLLMProviders({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res.code === 0) {
      tableData.value = res.data.items || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  searchForm.provider = ''
  searchForm.status = ''
  pagination.page = 1
  loadData()
}

const openDialog = (mode: 'create' | 'edit', row?: LLMProvider) => {
  dialogMode.value = mode
  if (mode === 'edit' && row) {
    currentId.value = row.id
    Object.assign(form, {
      name: row.name,
      provider: row.provider,
      base_url: row.base_url,
      api_key: row.api_key || '',
      model_name: row.model_name,
      is_default: row.is_default,
      is_enabled: row.is_enabled,
      timeout: row.timeout,
      retry_count: row.retry_count,
      description: row.description
    })
  }
  dialogVisible.value = true
}

const resetForm = () => {
  formRef.value?.resetFields()
  Object.assign(form, {
    name: '',
    provider: '',
    base_url: '',
    api_key: '',
    model_name: '',
    is_default: false,
    is_enabled: true,
    timeout: 60,
    retry_count: 3,
    description: ''
  })
  currentId.value = null
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    submitting.value = true
    if (dialogMode.value === 'create') {
      await createLLMProvider(form)
      ElMessage.success('创建成功')
    } else {
      await updateLLMProvider(currentId.value!, form)
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

const handleTest = async (row: LLMProvider) => {
  try {
    await ElMessageBox.confirm('确定要测试此配置吗？', '测试连接', {
      type: 'info'
    })
    const res = await testLLMConnection(row.id)
    if (res.code === 0 && res.data?.success) {
      ElMessage.success('连接测试成功！')
    } else {
      ElMessage.error(res.data?.message || '连接测试失败')
    }
  } catch (error) {
    console.error('测试失败:', error)
  }
}

const handleDelete = async (row: LLMProvider) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除配置"${row.name}"吗？`,
      '删除确认',
      { type: 'warning' }
    )
    await deleteLLMProvider(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

onMounted(() => {
  loadProviders()
  loadData()
})
</script>

<style scoped>
.llm-provider-view {
  padding: 20px;
}
</style>
