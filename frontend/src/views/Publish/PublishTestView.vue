<template>
  <div>
    <!-- 页面头部 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-xiaohongshu-dark">🧪 发布功能测试</h1>
      <p class="text-gray-500 mt-1">测试内容在小红书平台的发布流程</p>
    </div>

    <!-- 主内容卡片 -->
    <div class="bg-white rounded-xl shadow-xiaohongshu p-5">
      <!-- 步骤指示器 -->
      <div class="mb-8">
        <el-steps :active="currentStep" finish-status="success" align-center simple>
          <el-step title="发布前预览" description="模拟小红书展示效果" />
          <el-step title="格式校验" description="检测平台格式要求" />
          <el-step title="合规性检查" description="检测违规内容" />
          <el-step title="发布执行" description="执行发布操作" />
        </el-steps>
      </div>

      <!-- 内容选择区域 -->
      <div class="mb-8">
        <h3 class="text-lg font-semibold text-gray-800 mb-4">选择要发布的内容</h3>
        <el-select
          v-model="selectedContent"
          placeholder="请选择要发布的内容"
          style="width: 100%"
          size="large"
          @change="handleContentSelect"
        >
          <el-option
            v-for="content in contentList"
            :key="content.id"
            :label="content.title"
            :value="content"
          />
        </el-select>
      </div>

      <!-- 测试结果区域 -->
      <div v-if="selectedContent">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- 左侧：小红书预览效果 -->
          <div>
            <div class="bg-gray-50 rounded-xl p-6">
              <h4 class="font-semibold text-gray-800 mb-4 flex items-center gap-2">
                <span class="w-2 h-2 bg-xiaohongshu-red rounded-full"></span>
                小红书预览效果
              </h4>
              <div class="bg-white rounded-xl border border-gray-100 p-4">
                <!-- 用户信息 -->
                <div class="flex items-center gap-3 mb-4">
                  <el-avatar :size="40" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
                  <div>
                    <div class="font-medium text-gray-800">测试用户</div>
                    <div class="text-xs text-gray-400">刚刚</div>
                  </div>
                </div>
                <!-- 内容预览 -->
                <div class="mb-4">
                  <h3 class="text-lg font-semibold text-gray-800 mb-3">{{ selectedContent.title }}</h3>
                  <div class="text-gray-600 leading-relaxed" v-html="selectedContent.content" />
                </div>
                <!-- 话题标签 -->
                <div class="flex flex-wrap gap-2">
                  <el-tag
                    v-for="tag in selectedContent.tags"
                    :key="tag"
                    size="small"
                    class="bg-xiaohongshu-red/10 text-xiaohongshu-red border-xiaohongshu-red/20"
                  >
                    #{{ tag }}
                  </el-tag>
                </div>
              </div>
            </div>
          </div>

          <!-- 右侧：检查结果 -->
          <div>
            <div class="bg-gray-50 rounded-xl p-6">
              <div class="flex items-center justify-between mb-4">
                <h4 class="font-semibold text-gray-800 flex items-center gap-2">
                  <span class="w-2 h-2 bg-blue-500 rounded-full"></span>
                  检查结果
                </h4>
                <el-button
                  type="primary"
                  @click="runChecks"
                  :loading="checking"
                  class="bg-xiaohongshu-red border-none hover:opacity-90"
                >
                  开始检查
                </el-button>
              </div>

              <div class="space-y-3">
                <div
                  v-for="check in checkList"
                  :key="check.id"
                  class="bg-white rounded-lg p-4 border border-gray-100"
                >
                  <div class="flex items-center justify-between">
                    <span class="font-medium text-gray-700">{{ check.name }}</span>
                    <div class="flex items-center gap-2">
                      <el-icon
                        v-if="check.status === 'pass'"
                        class="text-green-500"
                        :size="18"
                      >
                        <CircleCheck />
                      </el-icon>
                      <el-icon
                        v-else-if="check.status === 'fail'"
                        class="text-red-500"
                        :size="18"
                      >
                        <CircleClose />
                      </el-icon>
                      <el-icon
                        v-else
                        class="text-gray-400"
                        :size="18"
                      >
                        <Clock />
                      </el-icon>
                      <span
                        :class="[
                          'text-sm font-medium',
                          check.status === 'pass' ? 'text-green-600' :
                          check.status === 'fail' ? 'text-red-600' : 'text-gray-500'
                        ]"
                      >
                        {{ check.status === 'pass' ? '通过' : check.status === 'fail' ? '未通过' : '待检查' }}
                      </span>
                    </div>
                  </div>
                  <div v-if="check.message" class="mt-2 text-sm text-red-500">
                    {{ check.message }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="mt-8 flex justify-center gap-4">
          <el-button
            type="primary"
            size="large"
            :disabled="!canPublish"
            :loading="publishing"
            class="px-8 bg-gradient-to-r from-xiaohongshu-red to-xiaohongshu-pink border-none hover:opacity-90"
            @click="handlePublish"
          >
            {{ publishing ? '发布中...' : '执行发布' }}
          </el-button>
          <el-button size="large" @click="handleSaveDraft">
            保存草稿
          </el-button>
        </div>
      </div>
    </div>

    <!-- 小贴士 -->
    <div class="mt-6 bg-xiaohongshu-red/5 border border-xiaohongshu-red/20 rounded-xl p-4">
      <div class="flex items-start gap-3">
        <span class="text-xl">💡</span>
        <div>
          <h4 class="font-medium text-gray-800 mb-1">发布前检查清单</h4>
          <ul class="text-sm text-gray-600 space-y-1">
            <li>• 确保标题吸引人，长度在10-30字之间</li>
            <li>• 正文内容清晰，段落分明，使用表情符号增加可读性</li>
            <li>• 添加3-5个相关话题标签，提高内容曝光</li>
            <li>• 检查是否有违规内容，避免被平台限流</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { CircleCheck, CircleClose, Clock } from '@element-plus/icons-vue'

// 当前步骤
const currentStep = ref(0)

// 选中的内容
const selectedContent = ref<any>(null)

// 加载状态
const checking = ref(false)
const publishing = ref(false)

// 内容列表
const contentList = ref([
  {
    id: 1,
    title: '夏日穿搭分享',
    content: '<p>今天给大家分享一套超好看的夏日穿搭✨</p><p>上衣选了一件淡紫色的吊带，面料超级舒服，版型也很显瘦！</p><p>下装搭配了一条白色牛仔短裤，简单又清爽～</p>',
    tags: ['夏日穿搭', '穿搭分享', 'OOTD']
  },
  {
    id: 2,
    title: '探店｜这家咖啡店太好拍了',
    content: '<p>周末和闺蜜一起去了一家藏在巷子里的咖啡店☕️</p><p>装修风格是ins风，超级适合拍照！</p><p>咖啡也很好喝，推荐大家试试他们的招牌拿铁～</p>',
    tags: ['探店', '咖啡店', '拍照打卡']
  }
])

// 检查列表
const checkList = ref([
  { id: 1, name: '标题长度检查', status: 'pending', message: '' },
  { id: 2, name: '正文内容检查', status: 'pending', message: '' },
  { id: 3, name: '话题标签检查', status: 'pending', message: '' },
  { id: 4, name: '敏感词检测', status: 'pending', message: '' },
  { id: 5, name: '图片合规性检查', status: 'pending', message: '' }
])

// 是否可以发布
const canPublish = computed(() => {
  return selectedContent.value && checkList.value.every(check => check.status === 'pass')
})

// 选择内容
const handleContentSelect = () => {
  currentStep.value = 0
  checkList.value.forEach(check => {
    check.status = 'pending'
    check.message = ''
  })
}

// 运行检查
const runChecks = async () => {
  checking.value = true
  currentStep.value = 1

  // 模拟检查过程
  for (let i = 0; i < checkList.value.length; i++) {
    await new Promise(resolve => setTimeout(resolve, 500))
    const check = checkList.value[i]
    check.status = Math.random() > 0.15 ? 'pass' : 'fail'
    check.message = check.status === 'fail' ? '发现一些小问题需要修改' : ''
    
    if (i === 1) currentStep.value = 2
    if (i === 3) currentStep.value = 3
  }

  checking.value = false
  if (canPublish.value) {
    ElMessage.success('所有检查通过，可以发布！')
  } else {
    ElMessage.warning('部分检查未通过，请修改后重试')
  }
}

// 执行发布
const handlePublish = async () => {
  if (!canPublish.value) return

  publishing.value = true
  currentStep.value = 3

  // 模拟发布过程
  await new Promise(resolve => setTimeout(resolve, 2000))

  publishing.value = false
  ElMessage.success('🎉 发布成功！内容已成功发布到小红书平台')
}

// 保存草稿
const handleSaveDraft = () => {
  ElMessage.success('草稿已保存')
}
</script>
