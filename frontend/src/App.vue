<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => {
  console.log('小红书内容生成与管理系统 v2.0.0 已启动')
})
</script>

<style>
#app {
  width: 100%;
  height: 100%;
}
</style>
