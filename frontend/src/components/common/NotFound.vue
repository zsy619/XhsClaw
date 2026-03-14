<template>
  <div class="not-found-container">
    <div class="not-found-content">
      <h1 class="error-code">404</h1>
      <p class="error-message">页面未找到</p>
      <p class="error-description">抱歉，您访问的页面不存在</p>
      <el-button type="primary" @click="goHome">返回首页</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/dashboard')
}
</script>

<style scoped lang="scss">
.not-found-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;

  .not-found-content {
    text-align: center;

    .error-code {
      font-size: 120px;
      font-weight: 700;
      color: #ff2442;
      margin-bottom: 16px;
    }

    .error-message {
      font-size: 24px;
      color: #333;
      margin-bottom: 8px;
    }

    .error-description {
      font-size: 14px;
      color: #666;
      margin-bottom: 32px;
    }
  }
}
</style>
