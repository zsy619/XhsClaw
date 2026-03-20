<template>
  <div class="system-dict-view">
    <div class="flex justify-between items-center mb-4">
      <div>
        <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center gap-2">
          <el-icon class="text-primary-500"><Document /></el-icon>
          系统字典
        </h1>
        <p class="mt-1 text-sm text-gray-500">管理系统字典数据</p>
      </div>
      <el-button type="primary" @click="openDialog('create')">
        <el-icon><Plus /></el-icon>
        新增字典
      </el-button>
    </div>

    <!-- 分类筛选 -->
    <el-card class="mb-4">
      <el-form :inline="true">
        <el-form-item label="字典分类">
          <el-select v-model="searchCategory" placeholder="选择分类" clearable @change="loadData">
            <el-option
              v-for="item in categoryOptions"
              :key="item"
              :label="getCategoryName(item)"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="searchKeyword" placeholder="搜索名称/编码" clearable @input="handleSearch" />
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 字典列表 -->
    <el-card>
      <el-table :data="filteredData" v-loading="loading" stripe>
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag type="info">{{ getCategoryName(row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="code" label="编码" width="150" />
        <el-table-column prop="name" label="名称" min-width="120" />
        <el-table-column prop="value" label="值" min-width="200" show-overflow-tooltip />
        <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
        <el-table-column prop="sort_order" label="排序" width="80" align="center" />
        <el-table-column prop="enabled" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'danger'" size="small">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDialog('edit', row)">编辑</el-button>
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
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogMode === 'create' ? '新增字典' : '编辑字典'"
      width="600px"
      @close="resetForm"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="分类" prop="category">
          <el-select v-model="form.category" placeholder="选择分类" style="width: 100%">
            <el-option
              v-for="item in categoryOptions"
              :key="item"
              :label="getCategoryName(item)"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="编码" prop="code">
          <el-input v-model="form.code" placeholder="字典编码" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="字典名称" />
        </el-form-item>
        <el-form-item label="值">
          <el-input v-model="form.value" type="textarea" :rows="2" placeholder="字典值" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" placeholder="字典描述" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Document } from '@element-plus/icons-vue'
import {
  getAllDicts,
  createDict,
  updateDict,
  deleteDict,
  getDictCategories,
  type SystemDict,
  type SystemDictRequest
} from '@/api/systemDict'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit'>('create')
const currentId = ref<number | null>(null)

const tableData = ref<SystemDict[]>([])
const searchCategory = ref('')
const searchKeyword = ref('')
const categoryOptions = ref<string[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 50,
  total: 0
})

// 分类名称映射
const categoryNames: Record<string, string> = {
  llm_provider: '大模型服务商',
  llm_model: '大模型',
  xhs_status: '小红书配置状态'
}

const getCategoryName = (category: string) => {
  return categoryNames[category] || category
}

const filteredData = computed(() => {
  if (!searchKeyword.value) return tableData.value
  const keyword = searchKeyword.value.toLowerCase()
  return tableData.value.filter(
    (item) =>
      item.name.toLowerCase().includes(keyword) ||
      item.code.toLowerCase().includes(keyword) ||
      item.description?.toLowerCase().includes(keyword)
  )
})

const formRef = ref()
const form = reactive<SystemDictRequest>({
  category: '',
  code: '',
  name: '',
  value: '',
  description: '',
  sort_order: 0,
  enabled: true
})

const rules = {
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  code: [{ required: true, message: '请输入编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }]
}

const loadCategories = async () => {
  try {
    const res = await getDictCategories()
    if (res.code === 0) {
      categoryOptions.value = res.data || []
    }
  } catch (error) {
    console.error('加载分类失败:', error)
    categoryOptions.value = ['llm_provider', 'llm_model', 'xhs_status']
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getAllDicts({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    if (res.code === 0) {
      let data = res.data.items || []
      if (searchCategory.value) {
        data = data.filter((item: SystemDict) => item.category === searchCategory.value)
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

const handleSearch = () => {
  // 搜索是实时的，通过 computed 属性 filteredData 实现
}

const openDialog = (mode: 'create' | 'edit', row?: SystemDict) => {
  dialogMode.value = mode
  if (mode === 'edit' && row) {
    currentId.value = row.id
    Object.assign(form, {
      category: row.category,
      code: row.code,
      name: row.name,
      value: row.value || '',
      description: row.description || '',
      sort_order: row.sort_order || 0,
      enabled: row.enabled
    })
  } else {
    form.category = searchCategory.value || ''
  }
  dialogVisible.value = true
}

const resetForm = () => {
  formRef.value?.resetFields()
  Object.assign(form, {
    category: '',
    code: '',
    name: '',
    value: '',
    description: '',
    sort_order: 0,
    enabled: true
  })
  currentId.value = null
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    submitting.value = true
    if (dialogMode.value === 'create') {
      await createDict(form)
      ElMessage.success('创建成功')
    } else {
      await updateDict(currentId.value!, form)
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

const handleDelete = async (row: SystemDict) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除字典"${row.name}"吗？`,
      '删除确认',
      { type: 'warning' }
    )
    await deleteDict(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

onMounted(() => {
  loadCategories()
  loadData()
})
</script>

<style scoped>
.system-dict-view {
  padding: 20px;
}
</style>
