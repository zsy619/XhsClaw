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
                  class="max-h-64 max-w-full object-contain rounded-lg shadow-md cursor-pointer hover:opacity-80 transition-opacity"
                  @click="openImageViewer(currentImageIndex)"
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
      width="680px"
      :close-on-click-modal="false"
      class="image-config-dialog"
    >
      <div class="image-config-content">
        <!-- 图片样式主题 -->
        <div class="config-section">
          <div class="section-header mb-4">
            <h3 class="text-base font-semibold text-gray-800 flex items-center gap-2">
              <el-icon class="text-primary-500"><Brush /></el-icon>
              图片样式主题
            </h3>
            <p class="text-xs text-gray-500">选择适合您内容的视觉风格</p>
          </div>
          
          <!-- 自定义下拉选择器 -->
          <div class="custom-select">
            <div 
              class="select-trigger" 
              @click="toggleStyleDropdown"
              :class="{ 'active': showStyleDropdown }"
            >
              <span class="selected-icon">{{ getCurrentStyle?.icon || '🎨' }}</span>
              <span class="selected-label">{{ getCurrentStyle?.label || '请选择样式主题' }}</span>
              <el-icon class="arrow-icon" :class="{ 'rotate': showStyleDropdown }">
                <ArrowDown />
              </el-icon>
            </div>
            
            <div 
              v-if="showStyleDropdown" 
              class="select-dropdown"
            >
              <div 
                v-for="style in styleOptions" 
                :key="style.value"
                class="select-option"
                :class="{ 'selected': imageConfig.styleKey === style.value }"
                @click="selectStyle(style.value)"
              >
                <span class="option-icon">{{ style.icon }}</span>
                <div class="option-content">
                  <span class="option-label">{{ style.label }}</span>
                  <span class="option-desc">{{ style.description }}</span>
                </div>
                <el-icon v-if="imageConfig.styleKey === style.value" class="check-icon">
                  <CircleCheck />
                </el-icon>
              </div>
            </div>
          </div>
          
          <!-- 选中样式预览 -->
          <div v-if="false" class="selected-style-preview mt-4 p-4 rounded-xl border-2 border-gray-200 bg-gray-50">
            <div class="flex items-center gap-3">
              <div 
                class="style-preview w-12 h-12 rounded-xl flex items-center justify-center text-xl"
                :class="getCurrentStyle?.previewClass || 'bg-gray-200'"
              >
                {{ getCurrentStyle?.icon || '🎨' }}
              </div>
              <div class="flex-1">
                <div class="font-semibold text-gray-800 text-lg">{{ getCurrentStyle?.label || '请选择' }}</div>
                <div class="text-sm text-gray-500">{{ getCurrentStyle?.description || '选择一种样式主题' }}</div>
              </div>
              <el-icon class="text-primary-500" :size="24">
                <CircleCheck v-if="imageConfig.styleKey" />
              </el-icon>
            </div>
          </div>
        </div>

        <!-- 分割线 -->
        <div class="my-6 border-t border-gray-100"></div>

        <!-- 智能分页 -->
        <div class="config-section">
          <div class="flex items-center justify-between p-4 bg-gray-50 rounded-xl">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center">
                <el-icon class="text-blue-500"><MagicStick /></el-icon>
              </div>
              <div>
                <div class="font-medium text-gray-800">智能分页</div>
                <div class="text-xs text-gray-500">自动拆分长内容到多张卡片</div>
              </div>
            </div>
            <el-switch 
              v-model="imageConfig.enableSmartPagination"
              size="large"
              active-color="#10b981"
            />
          </div>
        </div>

        <!-- 卡片尺寸配置 -->
        <div class="config-section mt-6">
          <div class="section-header mb-4">
            <h3 class="text-base font-semibold text-gray-800 flex items-center gap-2">
              <el-icon class="text-primary-500"><Edit /></el-icon>
              卡片尺寸配置
            </h3>
            <p class="text-xs text-gray-500">设置生成图片的尺寸，支持常见比例</p>
          </div>
          
          <!-- 预设尺寸 -->
          <div class="mb-4">
            <div class="text-xs font-medium text-gray-600 mb-2">快速选择</div>
            <div class="flex flex-wrap gap-2">
              <el-button
                v-for="preset in sizePresets"
                :key="preset.label"
                size="small"
                @click="imageConfig.cardWidth = preset.width; imageConfig.cardHeight = preset.height"
                :type="imageConfig.cardWidth === preset.width && imageConfig.cardHeight === preset.height ? 'primary' : 'default'"
                class="preset-btn"
              >
                {{ preset.label }}
              </el-button>
            </div>
          </div>

          <!-- 自定义尺寸 -->
          <div class="grid grid-cols-2 gap-4">
            <div class="p-4 bg-gray-50 rounded-xl">
              <label class="block text-sm font-medium text-gray-700 mb-2">卡片宽度</label>
              <div class="flex items-center gap-2">
                <el-button
                  size="small"
                  @click="imageConfig.cardWidth = Math.max(720, imageConfig.cardWidth - 40)"
                  :disabled="imageConfig.cardWidth <= 720"
                  class="flex-shrink-0"
                >
                  <el-icon><Minus /></el-icon>
                </el-button>
                <el-input-number
                  v-model="imageConfig.cardWidth"
                  :min="720"
                  :max="1440"
                  :step="40"
                  size="large"
                  class="flex-1"
                  controls-position="right"
                />
                <el-button
                  size="small"
                  @click="imageConfig.cardWidth = Math.min(1440, imageConfig.cardWidth + 40)"
                  :disabled="imageConfig.cardWidth >= 1440"
                  class="flex-shrink-0"
                >
                  <el-icon><Plus /></el-icon>
                </el-button>
              </div>
              <div class="text-xs text-gray-500 mt-1">范围: 720 - 1440 px</div>
            </div>
            
            <div class="p-4 bg-gray-50 rounded-xl">
              <label class="block text-sm font-medium text-gray-700 mb-2">卡片高度</label>
              <div class="flex items-center gap-2">
                <el-button
                  size="small"
                  @click="imageConfig.cardHeight = Math.max(960, imageConfig.cardHeight - 40)"
                  :disabled="imageConfig.cardHeight <= 960"
                  class="flex-shrink-0"
                >
                  <el-icon><Minus /></el-icon>
                </el-button>
                <el-input-number
                  v-model="imageConfig.cardHeight"
                  :min="960"
                  :max="1920"
                  :step="40"
                  size="large"
                  class="flex-1"
                  controls-position="right"
                />
                <el-button
                  size="small"
                  @click="imageConfig.cardHeight = Math.min(1920, imageConfig.cardHeight + 40)"
                  :disabled="imageConfig.cardHeight >= 1920"
                  class="flex-shrink-0"
                >
                  <el-icon><Plus /></el-icon>
                </el-button>
              </div>
              <div class="text-xs text-gray-500 mt-1">范围: 960 - 1920 px</div>
            </div>
          </div>
          
          <!-- 比例提示 -->
          <div class="mt-3 p-3 bg-primary-50 rounded-lg">
            <div class="text-xs text-primary-700 flex items-center gap-1">
              <el-icon :size="14"><InfoFilled /></el-icon>
              当前比例: {{ getAspectRatio() }}
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="image-config-footer">
          <!-- 进度条 -->
          <div v-if="imageRenderProgress > 0" class="progress-section mb-4">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-medium text-gray-700">正在生成...</span>
              <span class="text-sm text-primary-600 font-semibold">{{ imageRenderProgress }}%</span>
            </div>
            <el-progress
              :percentage="imageRenderProgress"
              :stroke-width="10"
              :show-text="false"
              striped
              striped-flow
            />
          </div>
          
          <div class="flex items-center justify-between">
            <el-button 
              size="large"
              @click="showImageRenderDialog = false"
              class="cancel-btn px-8"
            >
              取消
            </el-button>
            <div class="flex items-center gap-3">
              <el-button
                v-if="imageRenderProgress > 0"
                size="large"
                type="danger"
                @click="cancelImageRender"
                class="px-6"
              >
                <el-icon class="mr-1"><Close /></el-icon>
                取消生成
              </el-button>
              <el-button
                type="primary"
                size="large"
                :loading="renderingImages"
                @click="handleRenderImages"
                class="generate-btn px-8 h-12 text-base font-medium"
              >
                <el-icon class="mr-2"><Picture /></el-icon>
                {{ renderingImages ? '生成中...' : '开始生成' }}
              </el-button>
            </div>
          </div>
        </div>
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

    <!-- 图片查看器弹窗 -->
    <div
      v-if="showImageViewer"
      class="image-viewer-mask"
      @click.self="closeImageViewer"
      @keydown="handleViewerKeydown"
      tabindex="0"
      ref="imageViewerMask"
    >
      <!-- 关闭按钮 -->
      <div class="image-viewer-close" @click="closeImageViewer">
        <el-icon :size="28"><Close /></el-icon>
      </div>

      <!-- 上一张按钮 -->
      <div
        class="image-viewer-nav image-viewer-nav-left"
        :class="{ 'is-disabled': imageViewerIndex === 0 }"
        @click="prevImage"
      >
        <el-icon :size="36"><ArrowLeft /></el-icon>
      </div>

      <!-- 图片容器 -->
      <div class="image-viewer-content">
        <img
          :src="generatedImages[imageViewerIndex]"
          :alt="`图片 ${imageViewerIndex + 1}`"
          class="image-viewer-image"
        />
      </div>

      <!-- 下一张按钮 -->
      <div
        class="image-viewer-nav image-viewer-nav-right"
        :class="{ 'is-disabled': imageViewerIndex === generatedImages.length - 1 }"
        @click="nextImage"
      >
        <el-icon :size="36"><ArrowRight /></el-icon>
      </div>

      <!-- 底部工具栏 -->
      <div class="image-viewer-toolbar">
        <div class="image-viewer-info">
          <span class="image-viewer-count">{{ imageViewerIndex + 1 }} / {{ generatedImages.length }}</span>
        </div>
        <div class="image-viewer-actions">
          <el-button
            size="default"
            class="toolbar-btn"
            @click="handleDownloadImage(imageViewerIndex)"
          >
            <el-icon><Download /></el-icon>
            下载
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { http } from '@/api/request'
import { getRenderedImage, renderMarkdown } from '@/api/xiaohongshuRenderer'
import XiaohongshuEditor from '@/components/editor/XiaohongshuEditor.vue'
import { useUserStore } from '@/stores/user'
import {
    ArrowLeft,
    ArrowRight,
    Brush,
    CircleCheck,
    Close,
    Document,
    DocumentCopy,
    Download,
    Edit,
    InfoFilled,
    MagicStick,
    Minus,
    Picture,
    Plus,
    RefreshLeft,
    Star,
    Timer
} from '@element-plus/icons-vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const imageViewerMask = ref<HTMLElement | null>(null)

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
const showImageViewer = ref(false)
const imageViewerIndex = ref(0)

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
  '职场新人指南',
  '美妆教程',
  '健身打卡',
  '生活好物分享'
]

// 样式主题选项
const styleOptions = [
  {
    value: 'default',
    label: '简约灰',
    icon: '⚪',
    description: '简洁大方的通用风格',
    previewClass: 'bg-gray-100'
  },
  {
    value: 'xiaohongshu',
    label: '小红书红',
    icon: '🔴',
    description: '小红书平台经典红色主题',
    previewClass: 'bg-red-100'
  },
  {
    value: 'purple',
    label: '紫韵',
    icon: '💜',
    description: '梦幻紫色，优雅浪漫',
    previewClass: 'bg-purple-100'
  },
  {
    value: 'mint',
    label: '清新薄荷',
    icon: '🌿',
    description: '清新自然，清爽舒适',
    previewClass: 'bg-green-100'
  },
  {
    value: 'sunset',
    label: '日落橙',
    icon: '🌅',
    description: '温暖夕阳，活力满满',
    previewClass: 'bg-orange-100'
  },
  {
    value: 'ocean',
    label: '深海蓝',
    icon: '🌊',
    description: '深邃宁静，理性优雅',
    previewClass: 'bg-blue-100'
  },
  {
    value: 'elegant',
    label: '优雅白',
    icon: '⚪',
    description: '简约纯粹，经典永恒',
    previewClass: 'bg-gray-50'
  },
  {
    value: 'dark',
    label: '暗黑模式',
    icon: '🌙',
    description: '深邃神秘，炫酷潮流',
    previewClass: 'bg-gray-900'
  },
  {
    value: 'playful-geometric',
    label: '活泼几何',
    icon: '🔷',
    description: '充满活力的几何图形设计',
    previewClass: 'bg-indigo-100'
  },
  {
    value: 'neo-brutalism',
    label: '新野兽派',
    icon: '🟡',
    description: '粗旷有力的现代美学',
    previewClass: 'bg-yellow-100'
  },
  {
    value: 'botanical',
    label: '植物系',
    icon: '🌱',
    description: '清新自然的植物风格',
    previewClass: 'bg-emerald-100'
  },
  {
    value: 'professional',
    label: '专业商务',
    icon: '💼',
    description: '稳重专业的商务风格',
    previewClass: 'bg-slate-100'
  },
  {
    value: 'retro',
    label: '复古风格',
    icon: '📺',
    description: '怀旧复古的设计风格',
    previewClass: 'bg-amber-100'
  },
  {
    value: 'terminal',
    label: '终端风格',
    icon: '💻',
    description: '程序员专属终端风格',
    previewClass: 'bg-zinc-900'
  },
  {
    value: 'sketch',
    label: '手绘风格',
    icon: '✏️',
    description: '温馨可爱的手绘风格',
    previewClass: 'bg-pink-100'
  },
  {
    value: 'pink-cream',
    label: '粉色奶油',
    icon: '💕',
    description: '甜美温柔，少女心爆棚',
    previewClass: 'bg-pink-200'
  },
  {
    value: 'coral',
    label: '珊瑚粉',
    icon: '🪸',
    description: '活力清新，元气满满',
    previewClass: 'bg-rose-200'
  },
  {
    value: 'lavender',
    label: '薰衣草紫',
    icon: '💐',
    description: '清新淡雅，温柔浪漫',
    previewClass: 'bg-violet-200'
  },
  {
    value: 'cream',
    label: '奶黄包',
    icon: '🧁',
    description: '温暖甜蜜，清新可爱',
    previewClass: 'bg-amber-200'
  },
  {
    value: 'nordic',
    label: '北欧风格',
    icon: '🏔️',
    description: '清新冷淡，简约时尚',
    previewClass: 'bg-slate-200'
  },
  {
    value: 'peach',
    label: '蜜桃粉',
    icon: '🍑',
    description: '水嫩甜美，元气少女',
    previewClass: 'bg-orange-200'
  },
  {
    value: 'matcha',
    label: '抹茶绿',
    icon: '🍵',
    description: '日式清新，自然健康',
    previewClass: 'bg-green-200'
  },
  {
    value: 'cherry',
    label: '樱花浪漫',
    icon: '🌸',
    description: '樱花浪漫，甜美梦幻',
    previewClass: 'bg-pink-200'
  },
  {
    value: 'strawberry',
    label: '草莓甜心',
    icon: '🍓',
    description: '甜美可爱，元气少女',
    previewClass: 'bg-red-200'
  },
  {
    value: 'blueberry',
    label: '蓝莓之夜',
    icon: '🫐',
    description: '深邃静谧，优雅神秘',
    previewClass: 'bg-blue-300'
  },
  {
    value: 'grape',
    label: '葡萄紫',
    icon: '🍇',
    description: '优雅高贵，神秘魅惑',
    previewClass: 'bg-violet-300'
  },
  {
    value: 'lemon',
    label: '柠檬黄',
    icon: '🍋',
    description: '清新酸甜，元气满满',
    previewClass: 'bg-yellow-200'
  },
  {
    value: 'lavender-gray',
    label: '薰衣草灰',
    icon: '💜',
    description: '低调优雅，温柔气质',
    previewClass: 'bg-purple-200'
  },
  {
    value: 'rose',
    label: '玫瑰金',
    icon: '🌹',
    description: '优雅浪漫，温柔女人味',
    previewClass: 'bg-rose-300'
  },
  {
    value: 'sky-blue',
    label: '天空蓝',
    icon: '☁️',
    description: '清新明亮，心旷神怡',
    previewClass: 'bg-sky-200'
  },
  {
    value: 'candy-pink',
    label: '糖果粉',
    icon: '🍭',
    description: '甜美可爱，少女心爆棚',
    previewClass: 'bg-pink-100'
  },
  {
    value: 'peach-blossom',
    label: '桃花粉',
    icon: '🌸',
    description: '温柔浪漫，春天的气息',
    previewClass: 'bg-orange-100'
  },
  {
    value: 'mint-green',
    label: '薄荷绿',
    icon: '🌿',
    description: '清新自然，治愈系风格',
    previewClass: 'bg-green-100'
  },
  {
    value: 'lemon-yellow',
    label: '柠檬黄',
    icon: '🍋',
    description: '明亮活泼，充满阳光',
    previewClass: 'bg-yellow-100'
  },
  {
    value: 'strawberry-red',
    label: '草莓红',
    icon: '🍓',
    description: '甜美可爱，少女心十足',
    previewClass: 'bg-red-100'
  },
  {
    value: 'ocean-blue',
    label: '海洋蓝',
    icon: '🌊',
    description: '深邃宁静，理性优雅',
    previewClass: 'bg-blue-100'
  },
  {
    value: 'forest-green',
    label: '森林绿',
    icon: '🌲',
    description: '自然清新，充满生机',
    previewClass: 'bg-green-100'
  },
  {
    value: 'sunset-orange',
    label: '日落橙',
    icon: '🌅',
    description: '温暖夕阳，活力满满',
    previewClass: 'bg-orange-100'
  },
  {
    value: 'neon-pink',
    label: '霓虹粉',
    icon: '💖',
    description: '炫酷时尚，活力四射',
    previewClass: 'bg-pink-100'
  },
  {
    value: 'crystal-purple',
    label: '水晶紫',
    icon: '🔮',
    description: '优雅梦幻，神秘高贵',
    previewClass: 'bg-purple-100'
  },
  {
    value: 'ice-blue',
    label: '冰蓝色',
    icon: '❄️',
    description: '清新冷淡，简约时尚',
    previewClass: 'bg-blue-100'
  },
  {
    value: 'rose-gold',
    label: '玫瑰金',
    icon: '🌹',
    description: '优雅浪漫，温柔女人味',
    previewClass: 'bg-pink-100'
  },
  {
    value: 'sand-brown',
    label: '沙滩棕',
    icon: '🏖️',
    description: '温暖自然，度假风情',
    previewClass: 'bg-amber-100'
  },
  {
    value: 'kiwi-green',
    label: '猕猴桃绿',
    icon: '🥝',
    description: '清新自然，健康活力',
    previewClass: 'bg-green-100'
  },
  {
    value: 'blueberry-night',
    label: '蓝莓之夜',
    icon: '🫐',
    description: '深邃静谧，优雅神秘',
    previewClass: 'bg-blue-100'
  },
  {
    value: 'lavender-gray',
    label: '薰衣草灰',
    icon: '💜',
    description: '低调优雅，温柔气质',
    previewClass: 'bg-gray-100'
  },
  // ========== 新增 20 个小红书风格样式 ==========
  {
    value: 'lavender-purple',
    label: '薰衣草紫调',
    icon: '💜',
    description: '梦幻薰衣草紫色，优雅浪漫',
    previewClass: 'bg-purple-200'
  },
  {
    value: 'mint-breeze',
    label: '薄荷绿风',
    icon: '🌿',
    description: '清新薄荷绿色，自然治愈',
    previewClass: 'bg-green-100'
  },
  {
    value: 'sakura-blossom',
    label: '樱花粉樱',
    icon: '🌸',
    description: '浪漫樱花粉色，温柔少女心',
    previewClass: 'bg-pink-200'
  },
  {
    value: 'deep-ocean',
    label: '深海蓝调',
    icon: '🌊',
    description: '深邃海洋蓝色，宁静优雅',
    previewClass: 'bg-blue-200'
  },
  {
    value: 'sunset-glow',
    label: '日落橙光',
    icon: '🌅',
    description: '温暖日落橙色，活力满满',
    previewClass: 'bg-orange-200'
  },
  {
    value: 'peach-blossom',
    label: '蜜桃粉桃',
    icon: '🍑',
    description: '甜美蜜桃粉色，温柔可爱',
    previewClass: 'bg-orange-100'
  },
  {
    value: 'matcha-green',
    label: '抹茶绿调',
    icon: '🍵',
    description: '清新抹茶绿色，自然健康',
    previewClass: 'bg-green-200'
  },
  {
    value: 'blueberry-night',
    label: '蓝莓之夜',
    icon: '🫐',
    description: '深邃蓝莓紫色，神秘优雅',
    previewClass: 'bg-indigo-200'
  },
  {
    value: 'cherry-blossom',
    label: '樱花浪漫',
    icon: '🌸',
    description: '浪漫樱花粉色，甜美梦幻',
    previewClass: 'bg-pink-200'
  },
  {
    value: 'strawberry-sweet',
    label: '草莓甜心',
    icon: '🍓',
    description: '甜美草莓红色，少女心十足',
    previewClass: 'bg-red-200'
  },
  {
    value: 'lemon-yellow',
    label: '柠檬黄调',
    icon: '🍋',
    description: '明亮柠檬黄色，清新活力',
    previewClass: 'bg-yellow-200'
  },
  {
    value: 'lavender-gray',
    label: '薰衣草灰',
    icon: '💜',
    description: '低调薰衣草灰色，优雅气质',
    previewClass: 'bg-gray-100'
  },
  {
    value: 'rose-gold',
    label: '玫瑰金',
    icon: '🌹',
    description: '优雅玫瑰金色，温柔女人味',
    previewClass: 'bg-pink-100'
  },
  {
    value: 'sky-blue',
    label: '天空蓝',
    icon: '☁️',
    description: '清新天空蓝色，心旷神怡',
    previewClass: 'bg-blue-100'
  },
  {
    value: 'candy-pink',
    label: '糖果粉',
    icon: '🍭',
    description: '甜美糖果粉色，少女心爆棚',
    previewClass: 'bg-pink-100'
  },
  {
    value: 'mint-green',
    label: '薄荷绿',
    icon: '🌿',
    description: '清新薄荷绿色，治愈系风格',
    previewClass: 'bg-green-100'
  },
  {
    value: 'strawberry-red',
    label: '草莓红',
    icon: '🍓',
    description: '甜美草莓红色，少女心十足',
    previewClass: 'bg-red-100'
  },
  {
    value: 'ocean-blue',
    label: '海洋蓝',
    icon: '🌊',
    description: '深邃海洋蓝色，宁静优雅',
    previewClass: 'bg-blue-100'
  },
  {
    value: 'forest-green',
    label: '森林绿',
    icon: '🌲',
    description: '自然森林绿色，充满生机',
    previewClass: 'bg-green-100'
  },
  {
    value: 'sunset-orange',
    label: '日落橙',
    icon: '🌅',
    description: '温暖日落橙色，活力满满',
    previewClass: 'bg-orange-100'
  }
]

// 尺寸预设
const sizePresets = [
  { label: '小红书 3:4', width: 1080, height: 1440 },
  { label: '正方形 1:1', width: 1080, height: 1080 },
  { label: '横屏 4:3', width: 1440, height: 1080 },
  { label: '横屏 16:9', width: 1920, height: 1080 },
  { label: '手机壁纸 9:16', width: 1080, height: 1920 }
]

// 计算宽高比
const getAspectRatio = () => {
  const ratio = imageConfig.cardWidth / imageConfig.cardHeight
  if (Math.abs(ratio - 3/4) < 0.01) return '3:4 (小红书标准)'
  if (Math.abs(ratio - 1) < 0.01) return '1:1 (正方形)'
  if (Math.abs(ratio - 4/3) < 0.01) return '4:3 (横屏)'
  if (Math.abs(ratio - 16/9) < 0.01) return '16:9 (宽屏)'
  if (Math.abs(ratio - 9/16) < 0.01) return '9:16 (手机壁纸)'
  return `${ratio.toFixed(2)}:1`
}

// 获取当前选中的样式
const getCurrentStyle = computed(() => {
  return styleOptions.find(s => s.value === imageConfig.styleKey)
})

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

// 下拉菜单状态
const showStyleDropdown = ref(false)

// 切换下拉菜单
const toggleStyleDropdown = () => {
  showStyleDropdown.value = !showStyleDropdown.value
}

// 选择样式
const selectStyle = (value: string) => {
  imageConfig.styleKey = value
  showStyleDropdown.value = false
}

// 点击外部关闭下拉菜单
onMounted(() => {
  document.addEventListener('click', (e) => {
    const dropdown = document.querySelector('.custom-select')
    if (dropdown && !dropdown.contains(e.target as Node)) {
      showStyleDropdown.value = false
    }
  })
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
    // 步骤 1: 准备所有参数
    imageRenderProgress.value = 20
    
    // 获取当前选中的标题
    const selectedTitle = titleOptions.value.length > 0 
      ? titleOptions.value[selectedTitleIndex.value] 
      : (result.value.title || form.content)
    
    // 构建标签字符串
    const tagsText = result.value.tags ? result.value.tags.join(',') : ''
    
    // 步骤 2: 根据是否启用智能分页来决定是否调用 DeepSeek 重新生成内容
    let finalContent = result.value.content
    let finalTitle = selectedTitle
    let finalTags = result.value.tags || []
    
    if (imageConfig.enableSmartPagination) {
      imageRenderProgress.value = 40
      ElMessage.info('正在使用 AI 优化内容以适应分页...')
      
      try {
        // 调用 DeepSeek API 重新生成符合分页要求的内容
        const audiencesText = form.audiences.length > 0 ? form.audiences.join(', ') : ''
        const styleText = form.customStyle || form.style
        
        const aiResponse = await http.post('/generation/theme', {
          keywords: form.content,
          style_preference: styleText,
          target_audience: audiencesText,
          length: form.wordCount,
          enable_pagination: true,  // 启用分页模式
          pagination_format: 'markdown'  // 使用 Markdown 格式便于分页
        }, { timeout: 120000 })
        
        if (aiResponse.data?.generated_content) {
          finalContent = aiResponse.data.generated_content
          if (aiResponse.data.generated_title) {
            finalTitle = aiResponse.data.generated_title
          }
          if (aiResponse.data.generated_tags && aiResponse.data.generated_tags.length > 0) {
            finalTags = aiResponse.data.generated_tags
          }
        }
      } catch (aiError) {
        console.warn('AI 内容生成失败，使用原始内容:', aiError)
        ElMessage.warning('AI 优化失败，将使用原始内容')
      }
    }
    
    imageRenderProgress.value = 50
    
    // 步骤 3: 构建完整的 Markdown 内容（包含 YAML 头部）
    const markdownContent = buildMarkdownContent(finalTitle, finalContent, finalTags)
    
    // 步骤 4: 生成封面图片
    imageRenderProgress.value = 60
    ElMessage.info('正在生成封面...')
    
    try {
      const coverResponse = await http.post('/xiaohongshu-renderer/cover', {
        title: finalTitle,
        subtitle: finalTags.length > 0 ? finalTags.slice(0, 3).join(' ') : '',
        style_key: imageConfig.styleKey,
        output_prefix: `cover_${Date.now()}`,
        width: imageConfig.cardWidth,
        height: imageConfig.cardHeight
      })
      
      console.log('封面生成响应:', coverResponse)
      
      // 响应格式：{ code: 0, message: '成功', data: { success: true, message: '封面生成成功', image: '...' } }
      let coverUrl = ''
      if (coverResponse.data?.data?.image) {
        coverUrl = getRenderedImage(coverResponse.data.data.image)
      } else if (coverResponse.data?.image) {
        coverUrl = getRenderedImage(coverResponse.data.image)
      }
      
      if (coverUrl) {
        generatedImages.value.push(coverUrl)
        console.log('封面图片 URL:', coverUrl)
      } else {
        console.warn('封面图片 URL 为空')
      }
    } catch (coverError) {
      console.warn('封面生成失败:', coverError)
      // 封面生成失败不影响后续流程
    }
    
    // 步骤 5: 生成内容图片（带分页）
    imageRenderProgress.value = 75
    ElMessage.info('正在生成内容图片...')
    
    const response = await renderMarkdown({
      markdown_content: markdownContent,
      style_key: imageConfig.styleKey,
      output_prefix: `content_${Date.now()}`,
      enable_smart_pagination: imageConfig.enableSmartPagination,
      card_width: imageConfig.cardWidth,
      card_height: imageConfig.cardHeight,
      max_content_height: imageConfig.cardHeight * 3  // 允许最大高度为 3 倍卡片高度
    })

    imageRenderProgress.value = 90

    const renderData = response.data
    if (renderData && renderData.images && renderData.images.length > 0) {
      const contentImages = renderData.images.map((path: string) => {
        return getRenderedImage(path)
      })
      
      // 将内容图片添加到封面图片后面
      generatedImages.value = [...generatedImages.value, ...contentImages]
      
      currentImageIndex.value = 0
      imageRenderProgress.value = 100
      
      setTimeout(() => {
        showImageRenderDialog.value = false
        imageRenderProgress.value = 0
        ElMessage.success(`成功生成 ${generatedImages.value.length} 张图片（包含封面）`)
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

// 构建 Markdown 内容（包含 YAML 头部）
const buildMarkdownContent = (title: string, content: string, tags: string[]) => {
  let markdown = '---\n'
  markdown += `title: ${title}\n`
  if (tags && tags.length > 0) {
    markdown += `tags: [${tags.join(', ')}]\n`
  }
  markdown += '---\n\n'
  markdown += content
  return markdown
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

// 打开图片查看器
const openImageViewer = (index: number) => {
  imageViewerIndex.value = index
  showImageViewer.value = true
  // 等待 DOM 更新后聚焦到遮罩层，以便键盘事件生效
  setTimeout(() => {
    if (imageViewerMask.value) {
      imageViewerMask.value.focus()
    }
  }, 100)
}

// 关闭图片查看器
const closeImageViewer = () => {
  showImageViewer.value = false
}

// 查看上一张图片
const prevImage = () => {
  if (imageViewerIndex.value > 0) {
    imageViewerIndex.value--
  }
}

// 查看下一张图片
const nextImage = () => {
  if (imageViewerIndex.value < generatedImages.value.length - 1) {
    imageViewerIndex.value++
  }
}

// 键盘导航
const handleViewerKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    closeImageViewer()
  } else if (e.key === 'ArrowLeft') {
    prevImage()
  } else if (e.key === 'ArrowRight') {
    nextImage()
  }
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

    // 获取生成的图片路径列表
    const images = generatedImages.value.map((url: string) => {
      // 从完整URL中提取后端返回的图片路径
      // URL格式: http://localhost:8000/api/v1/xiaohongshu-renderer/image/xxx.png
      // 需要提取: /xiaohongshu-renderer/image/xxx.png
      const match = url.match(/(\/xiaohongshu-renderer\/image\/.*\.png)/)
      return match ? match[1] : ''
    }).filter((path: string) => path !== '')

    await http.post('/content/save', {
      title: title,
      title_options: titleOptions.value,
      selected_title_index: selectedTitleIndex.value,
      description: description,
      tags: tags,
      images: images,
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

// 图片配置弹窗样式
.image-config-dialog {
  :deep(.el-dialog__header) {
    padding: 24px 24px 16px;
    border-bottom: 1px solid #f3f4f6;
  }

  :deep(.el-dialog__title) {
    font-size: 20px;
    font-weight: 600;
    color: #1f2937;
  }

  :deep(.el-dialog__body) {
    padding: 24px;
  }

  :deep(.el-dialog__footer) {
    padding: 16px 24px 24px;
    border-top: 1px solid #f3f4f6;
  }
}

.image-config-content {
  .config-section {
    .section-header {
      h3 {
        margin: 0;
      }

      p {
        margin: 4px 0 0;
      }
    }
  }

  .style-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;

    .style-option {
      transition: all 0.2s ease;

      &:hover {
        transform: translateY(-2px);
      }

      .style-preview {
        flex-shrink: 0;
      }
    }
  }

  .preset-btn {
    transition: all 0.2s ease;

    &:hover {
      transform: translateY(-1px);
    }
  }
}

// 自定义下拉选择器样式
.custom-select {
  position: relative;
  width: 100%;
  
  .select-trigger {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px 20px;
    border: 2px solid #e5e7eb;
    border-radius: 12px;
    background: white;
    cursor: pointer;
    transition: all 0.2s ease;
    
    &:hover {
      border-color: #6366f1;
      box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
    }
    
    &.active {
      border-color: #6366f1;
      box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
    }
    
    .selected-icon {
      font-size: 20px;
      flex-shrink: 0;
    }
    
    .selected-label {
      flex: 1;
      font-size: 15px;
      font-weight: 500;
      color: #1f2937;
    }
    
    .arrow-icon {
      font-size: 16px;
      color: #6b7280;
      transition: transform 0.2s ease;
      
      &.rotate {
        transform: rotate(180deg);
      }
    }
  }
  
  .select-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    margin-top: 8px;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
    max-height: 480px;
    overflow-y: auto;
    z-index: 1000;
    
    // 使用 Grid 布局，每行两列
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
    padding: 8px;
    
    .select-option {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 12px;
      cursor: pointer;
      transition: all 0.2s ease;
      border-radius: 8px;
      border: 1px solid transparent;
      
      &:hover {
        background-color: #f9fafb;
        border-color: #e5e7eb;
      }
      
      &.selected {
        background-color: #f5f3ff;
        border-color: #6366f1;
      }
      
      .option-icon {
        font-size: 24px;
        flex-shrink: 0;
      }
      
      .option-content {
        flex: 1;
        min-width: 0; // 允许内容收缩
        
        .option-label {
          display: block;
          font-weight: 500;
          color: #1f2937;
          font-size: 14px;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
        
        .option-desc {
          display: block;
          color: #9ca3af;
          font-size: 12px;
          margin-top: 2px;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
      }
      
      .check-icon {
        color: #6366f1;
        font-size: 18px;
        flex-shrink: 0;
      }
    }
  }
}

// 选中样式预览
.selected-style-preview {
  transition: all 0.2s ease;
  
  &:hover {
    border-color: #6366f1;
  }
  
  .style-preview {
    flex-shrink: 0;
  }
}

.image-config-footer {
  .progress-section {
    .el-progress {
      :deep(.el-progress__text) {
        font-weight: 600;
      }
    }
  }

  .cancel-btn {
    border-color: #d1d5db;
    color: #4b5563;
    transition: all 0.2s ease;

    &:hover {
      border-color: #9ca3af;
      color: #374151;
      background-color: #f9fafb;
    }
  }

  .generate-btn {
    background: linear-gradient(135deg, #ff2442 0%, #ff6b81 100%);
    border: none;
    box-shadow: 0 4px 12px rgba(255, 36, 66, 0.3);
    transition: all 0.2s ease;

    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 6px 16px rgba(255, 36, 66, 0.4);
    }

    &:active {
      transform: translateY(0);
    }

    &.is-loading {
      opacity: 0.8;
    }
  }
}

// 响应式优化
@media (max-width: 768px) {
  .image-config-dialog {
    :deep(.el-dialog) {
      width: 95% !important;
      margin: 5vh auto !important;
    }
  }

  .image-config-content {
    .style-grid {
      grid-template-columns: 1fr;
    }
  }
}

// 图片查看器样式
.image-viewer-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.9);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  outline: none;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.image-viewer-close {
  position: absolute;
  top: 20px;
  right: 20px;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  cursor: pointer;
  transition: all 0.2s ease;
  z-index: 10000;

  &:hover {
    background: rgba(255, 255, 255, 0.2);
    transform: rotate(90deg);
  }
}

.image-viewer-nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  cursor: pointer;
  transition: all 0.2s ease;
  z-index: 10000;

  &:hover:not(.is-disabled) {
    background: rgba(255, 255, 255, 0.2);
    transform: translateY(-50%) scale(1.1);
  }

  &.is-disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }
}

.image-viewer-nav-left {
  left: 20px;
}

.image-viewer-nav-right {
  right: 20px;
}

.image-viewer-content {
  max-width: 90vw;
  max-height: 85vh;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: zoomIn 0.3s ease;
}

@keyframes zoomIn {
  from {
    transform: scale(0.9);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

.image-viewer-image {
  max-width: 100%;
  max-height: 85vh;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

.image-viewer-toolbar {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(10px);
  padding: 12px 24px;
  border-radius: 50px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  z-index: 10000;
}

.image-viewer-info {
  .image-viewer-count {
    color: #fff;
    font-size: 14px;
    font-weight: 500;
  }
}

.image-viewer-actions {
  .toolbar-btn {
    background: rgba(255, 255, 255, 0.15);
    border: none;
    color: #fff;
    border-radius: 20px;
    padding: 8px 20px;
    transition: all 0.2s ease;

    &:hover {
      background: rgba(255, 255, 255, 0.25);
    }
  }
}

// 图片查看器响应式
@media (max-width: 768px) {
  .image-viewer-nav {
    width: 44px;
    height: 44px;

    .el-icon {
      font-size: 28px !important;
    }
  }

  .image-viewer-nav-left {
    left: 10px;
  }

  .image-viewer-nav-right {
    right: 10px;
  }

  .image-viewer-close {
    width: 40px;
    height: 40px;
    top: 10px;
    right: 10px;

    .el-icon {
      font-size: 22px !important;
    }
  }

  .image-viewer-toolbar {
    padding: 10px 16px;
    gap: 16px;

    .toolbar-btn {
      padding: 6px 16px;
      font-size: 13px;
    }
  }
}
</style>
