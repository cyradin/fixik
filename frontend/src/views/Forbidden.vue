<template>
  <div class="forbidden-page">
    <el-result icon="error" title="403" sub-title="У вас нет доступа к этой странице">
      <template #extra>
        <el-button type="primary" @click="goHome"> На главную </el-button>

        <el-button @click="goBack"> Назад </el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}

const goBack = () => {
  router.back()
}
</script>

<style scoped>
.forbidden-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-bg-color-page);
}
</style>
