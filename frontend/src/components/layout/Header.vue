<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { useRouter } from 'vue-router'
import { ElDropdown, ElDropdownMenu, ElDropdownItem, ElButton, ElNotification } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import { notifyError } from '@/utils/notify'

const authStore = useAuthStore()
const router = useRouter()

const userName = computed(() => authStore.user?.name || authStore.user?.username || 'Пользователь')

const logout = async () => {
  try {
    await authStore.logout()
    router.push('/login')
  } catch {
    notifyError('Не удалось выйти из системы')
  }
}
</script>

<template>
  <el-header
    style="display: flex; justify-content: space-between; align-items: center; position: relative"
  >
    <h1 style="margin: 0">Fixik: Система инцидентов</h1>

    <el-dropdown v-if="authStore.isAuth">
      <el-button dashed :icon="User" style="display: flex; align-items: center">
        {{ userName }}
      </el-button>

      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item @click="logout">Выйти</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </el-header>
</template>
