<template>
  <div class="xiaohongshu-editor">
    <MdEditor
      v-model="content"
      :toolbars="toolbars"
      :preview="preview"
      :theme="theme"
      :read-only="readOnly"
      :editor-style="editorStyle"
      :preview-style="previewStyle"
      :screen-full="false"
      class="editor"
      @on-change="handleChange"
      @on-fullscreen="handleFullscreen"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

interface Props {
  modelValue?: string
  preview?: boolean
  readOnly?: boolean
  theme?: 'light' | 'dark'
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  preview: false,
  readOnly: false,
  theme: 'light'
})

const emit = defineEmits(['update:modelValue', 'change', 'fullscreen'])

const content = ref(props.modelValue)

// 编辑器工具栏配置（小红书风格）
const toolbars = [
  'bold',
  'underline',
  'italic',
  'strikeThrough',
  '-',
  'title',
  'sub',
  'sup',
  'quote',
  'unorderedList',
  'orderedList',
  '-',
  'codeRow',
  'code',
  'link',
  'image',
  'table',
  '-',
  'revoke',
  'next',
  '=',
  'pageFullscreen',
  'fullscreen',
  'preview',
  'htmlPreview',
  'catalog'
]

// 小红书风格编辑器样式
const editorStyle = computed(() => ({
  height: '100%',
  minHeight: '400px',
  background: props.theme === 'dark' ? '#1e1e1e' : '#ffffff',
  borderRadius: '8px'
}))

// 预览样式（小红书风格）
const previewStyle = computed(() => ({
  color: props.theme === 'dark' ? '#d4d4d4' : '#333333',
  fontSize: '15px',
  lineHeight: '1.8',
  fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
}))

watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal !== content.value) {
      content.value = newVal
    }
  }
)

watch(content, (newVal) => {
  emit('update:modelValue', newVal)
})

const handleChange = (value: string) => {
  emit('change', value)
}

const handleFullscreen = (status: boolean) => {
  emit('fullscreen', status)
}

// 暴露方法给父组件
defineExpose({
  getContent: () => content.value,
  setContent: (value: string) => {
    content.value = value
  }
})
</script>

<style scoped lang="scss">
.xiaohongshu-editor {
  width: 100%;
  height: 100%;
  
  .editor {
    height: 100%;
    min-height: 400px;
  }
}

// 小红书风格预览样式
:deep(.md-editor-preview-wrapper) {
  h1, h2, h3, h4, h5, h6 {
    font-weight: 600;
    margin: 20px 0 10px 0;
    color: #ff2442;
  }
  
  p {
    margin: 10px 0;
    line-height: 1.8;
  }
  
  strong {
    color: #ff2442;
    font-weight: 600;
  }
  
  blockquote {
    border-left: 4px solid #ff2442;
    padding: 10px 15px;
    margin: 15px 0;
    background: rgba(255, 36, 66, 0.05);
    border-radius: 0 8px 8px 0;
  }
  
  img {
    max-width: 100%;
    border-radius: 8px;
    margin: 15px 0;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
  
  code {
    background: #f5f7fa;
    padding: 2px 8px;
    border-radius: 4px;
    color: #ff2442;
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 14px;
  }
  
  pre {
    background: #1e1e1e;
    border-radius: 8px;
    padding: 15px;
    margin: 15px 0;
    overflow-x: auto;
    
    code {
      background: transparent;
      color: #d4d4d4;
      padding: 0;
    }
  }
  
  ul, ol {
    margin: 10px 0;
    padding-left: 25px;
  }
  
  li {
    margin: 5px 0;
    line-height: 1.8;
  }
  
  a {
    color: #ff2442;
    text-decoration: none;
    
    &:hover {
      text-decoration: underline;
    }
  }
}
</style>
