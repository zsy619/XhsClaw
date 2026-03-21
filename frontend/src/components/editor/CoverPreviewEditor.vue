<template>
  <div class="cover-preview-editor">
    <!-- 封面编辑器区域 -->
    <div class="editor-section mb-4">
      <div class="flex items-center justify-between mb-3">
        <div class="flex items-center gap-3">
          <h3 class="text-base font-semibold text-gray-800 flex items-center gap-2">
            <el-icon class="text-primary-500"><Edit /></el-icon>
            封面文案编辑
          </h3>
          <el-tag v-if="isModified" type="warning" size="small" effect="dark">
            已修改
          </el-tag>
        </div>
        <div class="flex items-center gap-3">
          <!-- 字数统计 -->
          <div class="text-xs text-gray-500">
            {{ wordCount }} 字
          </div>
          <!-- 编辑模式切换 -->
          <el-segmented
            v-model="editMode"
            :options="editModeOptions"
            size="small"
          />
        </div>
      </div>

      <!-- Markdown编辑器 -->
      <div v-show="editMode === 'edit'" class="editor-wrapper">
        <MdEditor
          v-model="localCoverSuggestion"
          :style="{ height: editorHeight }"
          :toolbars="(editorToolbars as any)"
          :language="language"
          :placeholder="placeholder"
          :theme="editorTheme"
          @on-change="handleChange"
          @on-save="handleSave"
        />
      </div>

      <!-- Markdown预览 -->
      <div v-show="editMode === 'preview'" class="preview-wrapper rounded-xl border-2 border-gray-200 p-6 bg-white min-h-[350px]">
        <div class="cover-preview-display">
          <div v-if="localCoverSuggestion" class="markdown-body" v-html="renderedHtml"></div>
          <div v-else class="flex flex-col items-center justify-center py-12 text-gray-400">
            <el-icon :size="48"><Document /></el-icon>
            <span class="mt-2">暂无封面文案</span>
            <span class="mt-1 text-xs">切换到编辑模式添加内容</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 工具提示 -->
    <div class="tips-section mb-4 p-3 bg-blue-50 rounded-lg border border-blue-200">
      <div class="flex items-start gap-2">
        <el-icon class="text-blue-500 mt-0.5" :size="16"><InfoFilled /></el-icon>
        <div class="text-xs text-blue-700">
          <p class="mb-1"><strong>💡 小贴士：</strong></p>
          <ul class="list-disc list-inside space-y-1">
            <li>使用 Markdown 语法可以快速格式化文案</li>
            <li>点击"快速模板"可以快速插入常用格式</li>
            <li>封面文案将用于生成小红书封面图片</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- 封面样式选择 -->
    <div class="style-section mb-4">
      <h3 class="text-sm font-semibold text-gray-700 flex items-center gap-2 mb-3">
        <el-icon class="text-primary-500"><Brush /></el-icon>
        封面样式
        <el-tooltip content="选择适合你内容的封面风格" placement="top">
          <el-icon class="text-gray-400 cursor-help"><QuestionFilled /></el-icon>
        </el-tooltip>
      </h3>
      <div class="grid grid-cols-2 gap-3">
        <div
          v-for="style in coverStyles"
          :key="style.value"
          @click="selectCoverStyle(style.value)"
          class="style-option p-4 rounded-xl border-2 cursor-pointer transition-all hover:shadow-md"
          :class="[
            selectedStyle === style.value
              ? 'border-primary-500 bg-primary-50 shadow-sm'
              : 'border-gray-200 hover:border-primary-300 bg-white'
          ]"
        >
          <div class="flex items-center gap-3">
            <div class="style-icon w-12 h-12 rounded-lg flex items-center justify-center text-2xl"
                 :class="style.iconBg">
              {{ style.icon }}
            </div>
            <div class="flex-1 min-w-0">
              <div class="text-sm font-semibold text-gray-800">{{ style.label }}</div>
              <div class="text-xs text-gray-500 mt-1">{{ style.desc }}</div>
            </div>
            <el-icon v-if="selectedStyle === style.value" class="text-primary-500" :size="20">
              <CircleCheckFilled />
            </el-icon>
          </div>
        </div>
      </div>
    </div>

    <!-- 快速模板 -->
    <div class="template-section">
      <h3 class="text-sm font-semibold text-gray-700 flex items-center gap-2 mb-3">
        <el-icon class="text-primary-500"><DocumentCopy /></el-icon>
        快速模板
        <el-tag type="success" size="small" effect="plain">
          点击应用
        </el-tag>
      </h3>
      <div class="flex flex-wrap gap-2">
        <el-tag
          v-for="template in quickTemplates"
          :key="template.label"
          class="template-tag cursor-pointer transition-all"
          effect="light"
          type="primary"
          @click="applyTemplate(template.content)"
        >
          <span class="mr-1">{{ template.icon }}</span>
          {{ template.label }}
        </el-tag>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import {
  Edit,
  Brush,
  Document,
  DocumentCopy,
  InfoFilled,
  QuestionFilled,
  CircleCheckFilled
} from '@element-plus/icons-vue'

// Props定义
interface Props {
  modelValue: string
  height?: string
  theme?: 'light' | 'dark'
}

const props = withDefaults(defineProps<Props>(), {
  height: '350px',
  theme: 'light'
})

// Emits定义
const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
  (e: 'save', value: string): void
}>()

// 编辑器配置
const editorHeight = computed(() => props.height)
const editorTheme = computed(() => props.theme)
const language = 'zh-CN'
const placeholder = '请输入封面文案...'

// 本地封面文案
const localCoverSuggestion = ref(props.modelValue)

// 编辑模式
const editMode = ref<'edit' | 'preview'>('edit')
const editModeOptions = [
  { label: '编辑', value: 'edit' },
  { label: '预览', value: 'preview' }
]

// 工具栏配置
const editorToolbars = [
  'bold',
  'italic',
  '-',
  'title',
  'quote',
  '-',
  'link',
  'image',
  'code',
  '=',
  'preview',
  'fullscreen'
]

// 封面样式选项
const coverStyles = [
  {
    value: 'text-only',
    label: '纯文字',
    icon: '📝',
    iconBg: 'bg-gray-100',
    desc: '简洁大气的纯文字风格'
  },
  {
    value: 'emoji-title',
    label: 'emoji+标题',
    icon: '🎨',
    iconBg: 'bg-pink-100',
    desc: 'emoji图标+主标题组合'
  },
  {
    value: 'contrast',
    label: '对比型',
    icon: '⚡',
    iconBg: 'bg-yellow-100',
    desc: '使用对比展示效果'
  },
  {
    value: 'data-highlight',
    label: '数据强调',
    icon: '📊',
    iconBg: 'bg-blue-100',
    desc: '突出数字和数据'
  }
]

// 选中的样式
const selectedStyle = ref('text-only')

// 快速模板
const quickTemplates = [
  {
    label: '痛点型',
    icon: '💔',
    content: `## 核心痛点
⚠️ 你是否也遇到这些问题？
- 问题1：XXX
- 问题2：XXX
- 问题3：XXX

## 解决方案
💡 别担心，有了这个方法...`
  },
  {
    label: '数字型',
    icon: '📈',
    content: `## 效果数据
🔥 3天见效！7天逆袭！
- 第1天：XXX
- 第3天：XXX
- 第7天：XXX

## 秘诀揭秘
💡 关键在于...`
  },
  {
    label: '对比型',
    icon: '⚡',
    content: `## 前后对比
❌ 之前：XXX
✅ 之后：XXX

## 变化过程
🌟 第一步：XXX
🌟 第二步：XXX
🌟 第三步：XXX`
  },
  {
    label: '悬念型',
    icon: '🤔',
    content: `## 揭秘真相
👀 你不知道的XXX...
❌ 误区一：XXX
❌ 误区二：XXX
❌ 误区三：XXX

## 正确方法
✅ 正确做法是...`
  },
  {
    label: '干货型',
    icon: '📚',
    content: `## 知识要点
📝 本期重点：
1. XXX
2. XXX
3. XXX

## 详细讲解
💡 XXX原理：
- 要点1
- 要点2
- 要点3`
  },
  {
    label: '种草型',
    icon: '🌱',
    content: `## 好物推荐
✨ 今天分享一款XXX
- 优点1：XXX
- 优点2：XXX
- 优点3：XXX

## 使用体验
💕 我的感受是...`
  }
]

// Markdown渲染HTML
const renderedHtml = computed(() => {
  if (!localCoverSuggestion.value) return ''
  // 简单的Markdown转HTML转换
  return markdownToHtml(localCoverSuggestion.value)
})

// 监听外部值变化
watch(() => props.modelValue, (newVal) => {
  if (newVal !== localCoverSuggestion.value) {
    localCoverSuggestion.value = newVal
  }
})

// 内容变化处理
const handleChange = (value: string) => {
  localCoverSuggestion.value = value
  emit('update:modelValue', value)
  emit('change', value)
}

// 保存处理
const handleSave = (value: string) => {
  emit('save', value)
}

// 选择封面样式
const selectCoverStyle = (style: string) => {
  selectedStyle.value = style
}

// 应用模板
const applyTemplate = (content: string) => {
  localCoverSuggestion.value = content
  emit('update:modelValue', content)
  emit('change', content)
  ElMessage.success('模板已应用')
}

// 简单的Markdown转HTML函数
const markdownToHtml = (markdown: string): string => {
  if (!markdown) return ''

  let html = markdown

  // 标题处理
  html = html.replace(/^### (.*$)/gim, '<h3 class="text-lg font-bold mt-4 mb-2">$1</h3>')
  html = html.replace(/^## (.*$)/gim, '<h2 class="text-xl font-bold mt-4 mb-2">$1</h2>')
  html = html.replace(/^# (.*$)/gim, '<h1 class="text-2xl font-bold mt-4 mb-2">$1</h1>')

  // 粗体和斜体
  html = html.replace(/\*\*\*(.*?)\*\*\*/gim, '<strong class="font-bold">$1</strong>')
  html = html.replace(/\*\*(.*?)\*\*/gim, '<strong class="font-bold">$1</strong>')
  html = html.replace(/\*(.*?)\*/gim, '<em class="italic">$1</em>')

  // 列表处理
  html = html.replace(/^- (.*$)/gim, '<li class="ml-4">$1</li>')
  html = html.replace(/^(\d+)\. (.*$)/gim, '<li class="ml-4 list-decimal">$2</li>')

  // 换行处理
  html = html.replace(/\n\n/gim, '</p><p class="mb-2">')
  html = html.replace(/\n/gim, '<br>')

  // 包裹在p标签中
  html = `<p class="mb-2">${html}</p>`

  // 清理空标签
  html = html.replace(/<p class="mb-2"><\/p>/gim, '')

  return html
}

// 引入ElMessage
import { ElMessage } from 'element-plus'

// 计算属性
const isModified = computed(() => {
  return localCoverSuggestion.value !== props.modelValue
})

const wordCount = computed(() => {
  return localCoverSuggestion.value.trim().length
})

// 暴露方法给父组件
defineExpose({
  getContent: () => localCoverSuggestion.value,
  setContent: (content: string) => {
    localCoverSuggestion.value = content
    emit('update:modelValue', content)
  },
  getSelectedStyle: () => selectedStyle.value
})
</script>

<style scoped lang="scss">
.cover-preview-editor {
  :deep(.md-editor) {
    border-radius: 12px;
    overflow: hidden;
    border: 2px solid #e5e7eb;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);

    &:focus-within {
      border-color: #6366f1;
      box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
    }
  }

  .editor-wrapper {
    position: relative;
    z-index: 1;
  }

  .preview-wrapper {
    min-height: 350px;
    max-height: 500px;
    overflow-y: auto;
    transition: all 0.3s ease;

    .markdown-body {
      font-size: 15px;
      line-height: 1.8;
      color: #374151;

      :deep(h1) {
        font-size: 24px;
        font-weight: 700;
        margin-top: 20px;
        margin-bottom: 12px;
        color: #1f2937;
      }

      :deep(h2) {
        font-size: 20px;
        font-weight: 600;
        margin-top: 16px;
        margin-bottom: 10px;
        color: #374151;
      }

      :deep(h3) {
        font-size: 18px;
        font-weight: 600;
        margin-top: 14px;
        margin-bottom: 8px;
        color: #4b5563;
      }

      :deep(p) {
        margin-bottom: 14px;
        line-height: 1.8;
      }

      :deep(strong) {
        color: #1f2937;
        font-weight: 600;
      }

      :deep(em) {
        color: #6b7280;
        font-style: italic;
      }

      :deep(li) {
        margin-bottom: 8px;
        line-height: 1.6;
      }

      :deep(ul),
      :deep(ol) {
        padding-left: 20px;
        margin-bottom: 14px;
      }
    }
  }

  .style-option {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(135deg, rgba(99, 102, 241, 0.05) 0%, rgba(99, 102, 241, 0.1) 100%);
      opacity: 0;
      transition: opacity 0.3s ease;
    }

    &:hover {
      transform: translateY(-3px);
      box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);

      &::before {
        opacity: 1;
      }
    }

    .style-icon {
      transition: transform 0.3s ease;
    }

    &:hover .style-icon {
      transform: scale(1.1);
    }
  }

  .template-tag {
    transition: all 0.2s ease;

    &:hover {
      transform: scale(1.05);
      box-shadow: 0 2px 8px rgba(99, 102, 241, 0.2);
    }
  }

  .tips-section {
    transition: all 0.3s ease;

    &:hover {
      background-color: #eff6ff;
    }
  }
}

// 响应式样式
@media (max-width: 768px) {
  .cover-preview-editor {
    :deep(.md-editor) {
      border-radius: 8px;
    }

    .preview-wrapper {
      min-height: 300px;
    }

    .style-option {
      padding: 12px;

      .style-icon {
        width: 40px !important;
        height: 40px !important;
        font-size: 18px !important;
      }
    }
  }
}
</style>
