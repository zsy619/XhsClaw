<template>
  <div class="creation-center-container" v-loading.fullscreen.lock="fullscreenLoading"
    element-loading-text="DeepSeek 正在为您火速创作中..."
    element-loading-background="rgba(255, 255, 255, 0.8)">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-xiaohongshu-dark flex items-center" style="gap: 2px;">
        <el-icon class="text-primary-500"><Edit /></el-icon>
        创作中心
      </h1>
      <p class="mt-1 text-sm text-gray-500">一站式文案生成、改写和图片渲染</p>
    </div>

    <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <!-- 左侧输入表单区 -->
      <div class="lg:col-span-1 space-y-6">
        <!-- 输入表单 -->
        <div class="rounded-xl bg-white p-6 shadow-xiaohongshu">
          <h2 class="mb-4 text-lg font-semibold text-xiaohongshu-dark">
            <el-icon class="mr-2 text-primary-500"><Document /></el-icon>
            输入信息
          </h2>

          <el-form
            :model="form"
            :rules="rules"
            ref="formRef"
            label-position="top"
            class="space-y-4"
          >
            <!-- 主题/内容输入框 -->
            <el-form-item label="主题/内容" prop="content">
              <el-input
                v-model="form.content"
                type="textarea"
                :rows="5"
                placeholder="请输入主题或要改写的内容..."
                maxlength="2000"
                show-word-limit
                class="w-full"
              />
            </el-form-item>

            <!-- 内容风格选择器 -->
            <el-form-item label="内容风格" prop="style">
              <el-select
                v-model="form.style"
                placeholder="请选择风格"
                class="w-full"
              >
                <el-option label="活泼可爱" value="cute">
                  <div class="flex items-center">
                    <span class="mr-2">😊</span>
                    <span>活泼可爱</span>
                  </div>
                </el-option>
                <el-option label="专业严谨" value="professional">
                  <div class="flex items-center">
                    <span class="mr-2">📚</span>
                    <span>专业严谨</span>
                  </div>
                </el-option>
                <el-option label="文艺清新" value="artistic">
                  <div class="flex items-center">
                    <span class="mr-2">🌸</span>
                    <span>文艺清新</span>
                  </div>
                </el-option>
                <el-option label="幽默风趣" value="humorous">
                  <div class="flex items-center">
                    <span class="mr-2">😂</span>
                    <span>幽默风趣</span>
                  </div>
                </el-option>
                <el-option label="干货分享" value="informative">
                  <div class="flex items-center">
                    <span class="mr-2">💡</span>
                    <span>干货分享</span>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>

            <!-- 自定义风格输入 -->
            <el-form-item label="自定义风格（可选）">
              <el-input
                v-model="form.customStyle"
                placeholder="输入您想要的风格描述..."
                maxlength="100"
              />
            </el-form-item>

            <!-- 目标受众选择 -->
            <el-form-item label="目标受众">
              <el-select
                v-model="form.audiences"
                multiple
                filterable
                allow-create
                placeholder="选择或输入目标受众"
                class="w-full"
              >
                <el-option label="18-25岁" value="18-25岁" />
                <el-option label="26-35岁" value="26-35岁" />
                <el-option label="36-45岁" value="36-45岁" />
                <el-option label="学生党" value="学生党" />
                <el-option label="职场新人" value="职场新人" />
                <el-option label="宝妈" value="宝妈" />
                <el-option label="健身爱好者" value="健身爱好者" />
                <el-option label="美食爱好者" value="美食爱好者" />
                <el-option label="旅游达人" value="旅游达人" />
                <el-option label="美妆博主" value="美妆博主" />
              </el-select>
            </el-form-item>

            <!-- 文案字数控制 -->
            <el-form-item label="文案字数">
              <div class="flex items-center gap-4">
                <el-slider
                  v-model="form.wordCount"
                  :min="50"
                  :max="1000"
                  :step="50"
                  show-stops
                  class="flex-1"
                />
                <el-input-number
                  v-model="form.wordCount"
                  :min="50"
                  :max="1000"
                  :step="50"
                  controls-position="right"
                  style="width: 120px"
                />
              </div>
              <p class="mt-1 text-xs text-gray-400">字数：{{ form.wordCount }} 字</p>
            </el-form-item>

            <!-- 操作按钮 -->
            <div class="flex flex-col gap-3 pt-2">
              <el-button-group class="flex w-full">
                <el-button
                  type="primary"
                  size="large"
                  :loading="generating"
                  @click="handleGenerate"
                  class="flex-1 h-12 text-base font-medium"
                >
                  <el-icon class="mr-2"><MagicStick /></el-icon>
                  生成文案
                </el-button>
                <el-button
                  type="warning"
                  size="large"
                  :loading="rewriting"
                  :disabled="!result?.content"
                  @click="handleRewrite"
                  class="flex-1 h-12 text-base font-medium"
                >
                  <el-icon class="mr-2"><Edit /></el-icon>
                  改写文案
                </el-button>
                <el-button
                  type="success"
                  size="large"
                  :loading="renderingImages"
                  :disabled="!result?.content"
                  @click="showImageRenderDialog = true"
                  class="flex-1 h-12 text-base font-medium"
                >
                  <el-icon class="mr-2"><Picture /></el-icon>
                  渲染图片
                </el-button>
              </el-button-group>
              <div class="flex gap-2">
                <el-button
                  v-if="hasHistory"
                  size="large"
                  @click="showHistoryDialog = true"
                  class="flex-1"
                >
                  <el-icon class="mr-1"><Timer /></el-icon>
                  历史记录
                </el-button>
                <el-button size="large" @click="handleReset" class="flex-1">
                  重置
                </el-button>
              </div>
            </div>
          </el-form>
        </div>

        <!-- 快捷提示 -->
        <div class="rounded-xl bg-primary-50 p-4">
          <h3 class="mb-2 text-sm font-medium text-primary-600">💡 小贴士</h3>
          <ul class="text-xs text-gray-600 space-y-1">
            <li>• 主题越具体，生成效果越好</li>
            <li>• 可以尝试不同风格找到最适合的</li>
            <li>• 生成后可在编辑器中进一步修改</li>
            <li>• 所有操作历史自动保存</li>
          </ul>
        </div>
      </div>

      <!-- 右侧预览区 -->
      <div class="lg:col-span-2">
        <div
          v-if="result"
          class="rounded-xl bg-white p-6 shadow-xiaohongshu min-h-[600px]"
        >
          <!-- 顶部操作栏 -->
          <div class="mb-4 flex flex-wrap items-center justify-between gap-4">
            <h2 class="text-lg font-semibold text-xiaohongshu-dark flex items-center">
              <el-icon class="mr-2 text-primary-500"><Star /></el-icon>
              创作结果
              <el-tag v-if="resultHistory.length > 0" type="info" size="small" class="ml-2">
                版本 {{ resultHistory.length }}
              </el-tag>
            </h2>
            <div class="flex flex-wrap gap-2">
              <el-button
                v-if="resultHistory.length > 1"
                size="small"
                @click="handleUndoRewrite"
                :disabled="currentHistoryIndex === 0"
              >
                <el-icon><RefreshLeft /></el-icon>
                撤销
              </el-button>
              <el-button size="small" @click="handleCopy">
                <el-icon><DocumentCopy /></el-icon>
                复制
              </el-button>
              <el-button
                v-if="generatedImages.length > 0"
                size="small"
                type="primary"
                @click="handleDownloadAllImages"
              >
                <el-icon><Download /></el-icon>
                下载全部
              </el-button>
              <el-button size="small" type="success" @click="handleSave">
                <el-icon><Download /></el-icon>
                保存
              </el-button>
            </div>
          </div>

          <!-- 标题备选方案与标签展示 -->
          <div v-if="result?.title || titleOptions.length > 0 || result?.tags" class="mb-6 space-y-4">
            
            <div v-if="result?.title" class="p-4 bg-primary-50 rounded-lg border-l-4 border-primary-500">
              <h3 class="text-lg font-bold text-gray-800 flex items-center">
                <el-icon class="mr-2 text-primary-500"><Document /></el-icon>
                {{ result.title }}
              </h3>
            </div>

            <div v-if="titleOptions.length > 0">
              <h3 class="mb-3 text-sm font-semibold text-gray-600">
                <el-icon class="mr-1"><Document /></el-icon>
                标题备选方案
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
                <div
                  v-for="(title, index) in titleOptions"
                  :key="index"
                  @click="selectTitle(index)"
                  :class="[
                    'cursor-pointer rounded-lg p-4 border-2 transition-all',
                    selectedTitleIndex === index
                      ? 'border-primary-500 bg-primary-50'
                      : 'border-gray-200 hover:border-primary-300'
                  ]"
                >
                  <div class="text-sm font-medium text-gray-800">{{ title }}</div>
                  <div class="mt-2 text-xs text-gray-400">点击选择</div>
                </div>
              </div>
            </div>

            <div v-if="result?.tags && result.tags.length > 0" class="flex flex-wrap gap-2">
              <el-tag
                v-for="(tag, index) in result.tags"
                :key="index"
                effect="light"
                round
                type="danger"
              >
                {{ tag }}
              </el-tag>
            </div>

          </div>

          <!-- 内容预览区 -->
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <!-- 图片预览 -->
            <div class="rounded-xl bg-xiaohongshu-bg p-4">
              <h3 class="mb-3 text-sm font-semibold text-gray-600 flex items-center">
                <el-icon class="mr-1"><Picture /></el-icon>
                图片预览
                <span v-if="generatedImages.length > 0" class="ml-2 text-xs text-gray-400">
                  ({{ currentImageIndex + 1 }}/{{ generatedImages.length }})
                </span>
              </h3>
              
              <!-- 图片切换按钮 -->
              <div v-if="generatedImages.length > 1" class="mb-3 flex items-center justify-center gap-2">
                <el-button
                  size="small"
                  :disabled="currentImageIndex === 0"
                  @click="currentImageIndex--"
                >
                  上一张
                </el-button>
                <el-button
                  size="small"
                  :disabled="currentImageIndex === generatedImages.length - 1"
                  @click="currentImageIndex++"
                >
                  下一张
                </el-button>
              </div>
              
              <div class="flex items-center justify-center rounded-lg bg-white p-8 min-h-[200px]">
                <img
                  v-if="generatedImages.length > 0"
                  :src="generatedImages[currentImageIndex]"
                  :alt="`生成的图片 ${currentImageIndex + 1}`"
                  class="max-h-64 max-w-full object-contain rounded-lg shadow-md"
                />
                <div v-else class="flex flex-col items-center text-gray-400">
                  <el-icon :size="40"><Picture /></el-icon>
                  <span class="mt-2">点击"渲染图片"生成</span>
                </div>
              </div>
              
              <!-- 下载单张按钮 -->
              <div v-if="generatedImages.length > 0" class="mt-3 flex justify-center">
                <el-button size="small" type="primary" @click="handleDownloadImage(currentImageIndex)">
                  <el-icon><Download /></el-icon>
                  下载这张
                </el-button>
              </div>
            </div>

            <!-- 文案预览 -->
            <div class="rounded-xl bg-primary-50/50 p-4">
              <h3 class="mb-3 text-sm font-semibold text-primary-600 flex items-center">
                <el-icon class="mr-1"><Document /></el-icon>
                文案预览
              </h3>
              <XiaohongshuEditor
                v-model="result.content"
                :preview="false"
                class="result-editor"
                @selection-change="handleSelectionChange"
              />
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div
          v-else
          class="flex min-h-[600px] flex-col items-center justify-center rounded-xl bg-white p-12 shadow-xiaohongshu"
        >
          <div class="mb-6 flex h-24 w-24 items-center justify-center rounded-full bg-primary-50">
            <el-icon :size="48" class="text-primary-400"><MagicStick /></el-icon>
          </div>
          <h3 class="mb-2 text-lg font-medium text-xiaohongshu-dark">开始您的创作</h3>
          <p class="text-center text-sm text-gray-500 max-w-md">
            在左侧输入信息，选择风格和受众，点击"生成文案"按钮<br />
            即可获得符合小红书平台风格的原创文案
          </p>

          <!-- 快捷示例 -->
          <div class="mt-8 w-full max-w-md">
            <p class="mb-3 text-sm font-medium text-gray-700">试试这些示例：</p>
            <div class="flex flex-wrap gap-2">
              <el-tag
                v-for="example in examples"
                :key="example"
                class="cursor-pointer hover:bg-primary-100"
                @click="fillExample(example)"
              >
                {{ example }}
              </el-tag>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 图片渲染弹窗 -->
    <el-dialog
      v-model="showImageRenderDialog"
      title="图片生成配置"
      width="600px"
      :close-on-click-modal="false"
    >
      <div class="space-y-4">
        <!-- 图片样式主题 -->
        <el-form-item label="图片样式主题">
          <el-select v-model="imageConfig.styleKey" class="w-full">
            <el-option label="简约灰" value="default" />
            <el-option label="小红书红" value="xiaohongshu" />
            <el-option label="活泼几何" value="playful-geometric" />
            <el-option label="新野兽派" value="neo-brutalism" />
            <el-option label="植物系" value="botanical" />
            <el-option label="专业商务" value="professional" />
            <el-option label="复古风格" value="retro" />
            <el-option label="终端风格" value="terminal" />
            <el-option label="手绘风格" value="sketch" />
          </el-select>
        </el-form-item>

        <!-- 智能分页 -->
        <el-form-item label="智能分页">
          <el-switch v-model="imageConfig.enableSmartPagination" />
          <span class="ml-2 text-sm text-gray-500">自动拆分长内容到多张卡片</span>
        </el-form-item>

        <!-- 卡片尺寸 -->
        <el-divider content-position="left">卡片尺寸配置</el-divider>
        
        <el-form-item label="卡片宽度">
          <el-input-number
            v-model="imageConfig.cardWidth"
            :min="720"
            :max="1440"
            :step="40"
            class="w-full"
          />
          <span class="ml-2 text-xs text-gray-400">px</span>
        </el-form-item>

        <el-form-item label="卡片高度">
          <el-input-number
            v-model="imageConfig.cardHeight"
            :min="960"
            :max="1920"
            :step="40"
            class="w-full"
          />
          <span class="ml-2 text-xs text-gray-400">px</span>
        </el-form-item>
      </div>

      <template #footer>
        <div class="flex justify-between">
          <el-button @click="showImageRenderDialog = false">取消</el-button>
          <div class="flex gap-2">
            <el-button
              v-if="imageRenderProgress > 0"
              @click="cancelImageRender"
            >
              取消生成
            </el-button>
            <el-button
              type="primary"
              :loading="renderingImages"
              @click="handleRenderImages"
            >
              开始生成
            </el-button>
          </div>
        </div>

        <!-- 进度条 -->
        <el-progress
          v-if="imageRenderProgress > 0"
          :percentage="imageRenderProgress"
          :stroke-width="8"
          class="mt-4"
        />
      </template>
    </el-dialog>

    <!-- 历史记录弹窗 -->
    <el-dialog
      v-model="showHistoryDialog"
      title="操作历史记录"
      width="800px"
    >
      <div class="max-h-96 overflow-y-auto">
        <el-timeline>
          <el-timeline-item
            v-for="(item, index) in resultHistory.slice().reverse()"
            :key="index"
            :timestamp="item.timestamp"
            placement="top"
          >
            <div class="flex items-center justify-between">
              <div>
                <el-tag :type="getHistoryItemType(item.type)" size="small">
                  {{ getHistoryItemLabel(item.type) }}
                </el-tag>
                <p class="mt-2 text-sm text-gray-600 line-clamp-2">
                  {{ item.preview }}
                </p>
              </div>
              <el-button size="small" @click="restoreHistory(item)">
                恢复
              </el-button>
            </div>
          </el-timeline-item>
        </el-timeline>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import {
  MagicStick,
  Edit,
  DocumentCopy,
  Download,
  Document,
  Star,
  Picture,
  Timer,
  RefreshLeft
} from '@element-plus/icons-vue'
import { http } from '@/api/request'
import { renderMarkdown, getRenderedImage } from '@/api/xiaohongshuRenderer'
import { useUserStore } from '@/stores/user'
import XiaohongshuEditor from '@/components/editor/XiaohongshuEditor.vue'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()

// 状态变量
const fullscreenLoading = ref(false)
const generating = ref(false)
const rewriting = ref(false)
const renderingImages = ref(false)
const result = ref<any>(null)
const showImageRenderDialog = ref(false)
const showHistoryDialog = ref(false)
const generatedImages = ref<string[]>([])
const currentImageIndex = ref(0)
const imageRenderProgress = ref(0)

// 标题相关
const titleOptions = ref<string[]>([])
const selectedTitleIndex = ref(0)

// 历史记录
const resultHistory = ref<any[]>([])
const currentHistoryIndex = ref(0)
const hasHistory = computed(() => resultHistory.value.length > 0)

// 快捷示例
const examples = [
  '夏日穿搭分享',
  '美食探店推荐',
  '旅行攻略',
  '护肤品测评',
  '职场新人指南'
]

// 表单数据
const form = reactive({
  content: '',
  style: 'cute',
  customStyle: '',
  audiences: [] as string[],
  wordCount: 300
})

// 图片配置
const imageConfig = reactive({
  styleKey: 'default',
  enableSmartPagination: true,
  cardWidth: 1080,
  cardHeight: 1440
})

// 验证规则
const rules: FormRules = {
  content: [
    { required: true, message: '请输入主题或内容', trigger: 'blur' },
    { min: 2, max: 2000, message: '内容长度在 2 到 2000 个字符', trigger: 'blur' }
  ],
  style: [
    { required: true, message: '请选择内容风格', trigger: 'change' }
  ]
}

// 检查登录状态
onMounted(() => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  // 从本地缓存恢复
  loadFromLocalCache()
})

// 自动保存到本地缓存
watch(
  () => result.value,
  (newVal) => {
    if (newVal) {
      saveToLocalCache()
    }
  },
  { deep: true }
)

// 保存到本地缓存
const saveToLocalCache = () => {
  const cacheData = {
    form: { ...form },
    result: result.value,
    resultHistory: resultHistory.value,
    timestamp: Date.now()
  }
  localStorage.setItem('creation_center_cache', JSON.stringify(cacheData))
}

// 从本地缓存加载
const loadFromLocalCache = () => {
  const cached = localStorage.getItem('creation_center_cache')
  if (cached) {
    try {
      const data = JSON.parse(cached)
      // 检查是否在24小时内
      if (Date.now() - data.timestamp < 24 * 60 * 60 * 1000) {
        Object.assign(form, data.form)
        result.value = data.result
        resultHistory.value = data.resultHistory || []
        currentHistoryIndex.value = resultHistory.value.length - 1
      }
    } catch (e) {
      console.error('加载缓存失败:', e)
    }
  }
}

// 填充示例
const fillExample = (example: string) => {
  form.content = example
  form.style = 'cute'
  form.audiences = []
  form.wordCount = 300
}

// 生成文案
const handleGenerate = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    fullscreenLoading.value = true

    try {
      const audiencesText = form.audiences.length > 0 ? form.audiences.join(', ') : ''
      const styleText = form.customStyle || form.style
      
      const res = await http.post('/generation/theme', {
        keywords: form.content,
        style_preference: styleText,
        target_audience: audiencesText,
        length: form.wordCount
      }, { timeout: 120000 })  // DeepSeek 生成可能较慢，延长超时到 120 秒

      const content = res.data?.generated_content || ''
      const title = res.data?.generated_title || ''
      const tags = res.data?.generated_tags || []
      
      // 生成标题备选方案
      titleOptions.value = [
        `${form.content}超全攻略！`,
        `分享我的${form.content}心得`,
        `必看！${form.content}干货整理`
      ]
      
      result.value = { content, title, tags }
      
      // 保存到历史记录
      addToHistory('generate', content)
      
      ElMessage.success('生成成功')
    } catch (error) {
      console.error('生成失败:', error)
      ElMessage.error('生成失败，请检查网络或后端的 DeepSeek 配置')
    } finally {
      fullscreenLoading.value = false
    }
  })
}

// 改写文案
const handleRewrite = async () => {
  if (!result.value?.content) {
    ElMessage.warning('请先生成文案')
    return
  }

  fullscreenLoading.value = true

  try {
    const styleText = form.customStyle || form.style
    
    const res = await http.post('/generation/rewrite', {
      content: result.value.content,
      style_preference: styleText,
      preserve_key_info: true,
      length: form.wordCount
    }, { timeout: 120000 })  // DeepSeek 改写可能较慢，延长超时到 120 秒

    const newContent = res.data?.generated_content || result.value.content
    const title = res.data?.generated_title || ''
    const tags = res.data?.generated_tags || []
    
    // 保存原文案用于撤销
    const oldContent = result.value.content
    
    result.value = { content: newContent, title, tags }
    addToHistory('rewrite', newContent, oldContent)
    
    ElMessage.success('改写成功')
  } catch (error) {
    console.error('改写失败:', error)
    ElMessage.error('改写失败，请检查网络或后端的 DeepSeek 配置')
  } finally {
    fullscreenLoading.value = false
  }
}

// 撤销改写
const handleUndoRewrite = () => {
  if (currentHistoryIndex.value > 0) {
    currentHistoryIndex.value--
    const item = resultHistory.value[currentHistoryIndex.value]
    result.value = { content: item.content }
  }
}

// 选中文字改写
const handleSelectionChange = (selection: any) => {
  // 可以在这里实现选中文字的针对性改写
}

// 选择标题
const selectTitle = (index: number) => {
  selectedTitleIndex.value = index
  // 可以在这里实现标题替换功能
  ElMessage.info(`已选择标题：${titleOptions.value[index]}`)
}

// 渲染图片
const handleRenderImages = async () => {
  if (!result.value?.content) {
    ElMessage.warning('请先生成文案')
    return
  }

  renderingImages.value = true
  imageRenderProgress.value = 10
  generatedImages.value = []
  
  try {
    imageRenderProgress.value = 30
    
    const response = await renderMarkdown({
      markdown_content: result.value.content,
      style_key: imageConfig.styleKey,
      enable_smart_pagination: imageConfig.enableSmartPagination,
      card_width: imageConfig.cardWidth,
      card_height: imageConfig.cardHeight,
      max_content_height: imageConfig.cardHeight - 340
    })

    imageRenderProgress.value = 70

    console.log('=== 调试信息 ===')
    console.log('完整响应:', response)
    console.log('response.data:', response.data)
    
    // 响应拦截器已经处理，response 就是 { code, message, data }
    // data 里面是 { success, message, images }
    const renderData = response.data
    console.log('renderData:', renderData)
    
    if (renderData) {
      const images = renderData.images || []
      console.log('图片路径列表:', images)
      
      generatedImages.value = images.map((path: string) => {
        const imageUrl = getRenderedImage(path)
        console.log('生成的图片URL:', imageUrl)
        return imageUrl
      })
      
      currentImageIndex.value = 0
      imageRenderProgress.value = 100
      
      setTimeout(() => {
        showImageRenderDialog.value = false
        imageRenderProgress.value = 0
        ElMessage.success(`成功生成 ${generatedImages.value.length} 张图片`)
      }, 500)
    } else {
      throw new Error(renderData?.message || response.message || '渲染失败')
    }
  } catch (error) {
    console.error('图片渲染失败:', error)
    imageRenderProgress.value = 0
    ElMessage.error('图片渲染失败，请稍后重试')
  } finally {
    renderingImages.value = false
  }
}

// 取消图片渲染
const cancelImageRender = () => {
  // 这里可以实现取消逻辑
  imageRenderProgress.value = 0
  ElMessage.info('已取消生成')
}

// 下载图片
const handleDownloadImage = async (index: number) => {
  const imageUrl = generatedImages.value[index]
  if (!imageUrl) {
    ElMessage.warning('图片地址无效')
    return
  }

  try {
    // 使用 fetch 下载图片，支持跨域
    const response = await fetch(imageUrl)
    if (!response.ok) {
      throw new Error('图片下载失败')
    }

    const blob = await response.blob()
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `xiaohongshu-${Date.now()}-${index + 1}.png`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success('下载已开始')
  } catch (error) {
    console.error('下载图片失败:', error)
    // 降级方案：直接打开链接
    window.open(imageUrl, '_blank')
    ElMessage.warning('直接打开图片，请手动保存')
  }
}

// 下载所有图片
const handleDownloadAllImages = () => {
  generatedImages.value.forEach((_, index) => {
    setTimeout(() => handleDownloadImage(index), index * 500)
  })
}

// 添加到历史记录
const addToHistory = (type: string, content: string, oldContent?: string) => {
  const historyItem = {
    type,
    content,
    oldContent,
    preview: content.substring(0, 100) + (content.length > 100 ? '...' : ''),
    timestamp: new Date().toLocaleString('zh-CN')
  }
  
  resultHistory.value.push(historyItem)
  currentHistoryIndex.value = resultHistory.value.length - 1
}

// 获取历史记录项类型
const getHistoryItemType = (type: string) => {
  const types: Record<string, any> = {
    generate: 'success',
    rewrite: 'warning'
  }
  return types[type] || 'info'
}

// 获取历史记录项标签
const getHistoryItemLabel = (type: string) => {
  const labels: Record<string, string> = {
    generate: '生成文案',
    rewrite: '改写文案'
  }
  return labels[type] || '操作'
}

// 恢复历史记录
const restoreHistory = (item: any) => {
  result.value = { content: item.content }
  showHistoryDialog.value = false
  ElMessage.success('已恢复')
}

// 重置
const handleReset = () => {
  form.content = ''
  form.style = 'cute'
  form.customStyle = ''
  form.audiences = []
  form.wordCount = 300
  result.value = null
  titleOptions.value = []
  selectedTitleIndex.value = 0
  generatedImages.value = []
  currentImageIndex.value = 0
  resultHistory.value = []
  currentHistoryIndex.value = 0
  
  // 清除本地缓存
  localStorage.removeItem('creation_center_cache')
  
  ElMessage.info('已重置')
}

// 复制
const handleCopy = () => {
  if (result.value?.content) {
    navigator.clipboard.writeText(result.value.content)
    ElMessage.success('复制成功')
  }
}

// 保存
const handleSave = async () => {
  if (!result.value?.content) return

  try {
    const title = result.value.title || form.content || '未命名内容'
    const description = result.value.content
    const tags = result.value.tags || []

    await http.post('/content/save', {
      title: title,
      title_options: titleOptions.value,
      selected_title_index: selectedTitleIndex.value,
      description: description,
      tags: tags,
      content_attributes: {
        content_style: form.style,
        custom_style: form.customStyle,
        target_audience: form.audiences
      },
      render_attributes: {
        image_style_theme: imageConfig.styleKey,
        enable_smart_pagination: imageConfig.enableSmartPagination,
        card_width: imageConfig.cardWidth,
        card_height: imageConfig.cardHeight
      }
    })
    ElMessage.success('保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败，请稍后重试')
  }
}
</script>

<style scoped lang="scss">
.creation-center-container {
  .result-editor {
    min-height: 400px;
  }
}
</style>
