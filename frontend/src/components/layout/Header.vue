<template>
  <el-header
    style="display: flex; justify-content: space-between; align-items: center; position: relative"
  >
    <h1 style="margin: 0">Fixik: Система инцидентов</h1>

    <el-menu v-if="authStore.isAuth" mode="horizontal" :ellipsis="false" router>
      <el-menu-item index="/">Инциденты</el-menu-item>
      <el-menu-item index="/admin/statuses" v-if="can(PERMISSION_GROUPS.STATUS_ADMIN)"
        >Статусы</el-menu-item
      >
      <el-menu-item index="/admin/priorities" v-if="can(PERMISSION_GROUPS.PRIORITY_ADMIN)"
        >Приоритеты</el-menu-item
      >
      <el-menu-item index="/admin/teams" v-if="can(PERMISSION_GROUPS.TEAM_ADMIN)"
        >Команды</el-menu-item
      >
      <el-menu-item index="/admin/users" v-if="can(PERMISSION_GROUPS.USER_ADMIN)"
        >Пользователи</el-menu-item
      >
    </el-menu>

    <el-dropdown v-if="authStore.isAuth">
      <el-button dashed :icon="User" style="display: flex; align-items: center">
        {{ userName }}
      </el-button>

      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item @click="goProfile">Профиль</el-dropdown-item>
          <el-dropdown-item @click="logout">Выйти</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </el-header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { useRouter } from 'vue-router'
import { ElDropdown, ElDropdownMenu, ElDropdownItem, ElButton, ElNotification } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import { notifyError } from '@/utils/notify'
import { PERMISSION_GROUPS } from '@/constants/permissions'

const authStore = useAuthStore()
const router = useRouter()

const can = authStore.can

const userName = computed(() => authStore.user?.name || authStore.user?.username || 'Пользователь')

const logout = async () => {
  try {
    await authStore.logout()
    router.push('/login')
  } catch {
    notifyError('Не удалось выйти из системы')
  }
}

const goProfile = () => {
  router.push('/profile')
}
</script>
