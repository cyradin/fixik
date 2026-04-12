<script setup lang="ts">
import { watch } from 'vue'

import { useIncidentsStore } from '@/stores/incidentsStore'
import { useStatusesStore } from '@/stores/statusesStore'
import { usePrioritiesStore } from '@/stores/prioritiesStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { useUsersStore } from '@/stores/usersStore'
import { useAuthStore } from '@/stores/authStore'

import Header from '@/components/layout/Header.vue'
import { useRolesStore } from './stores/rolesStore'

const incidentsStore = useIncidentsStore()
const statusesStore = useStatusesStore()
const prioritiesStore = usePrioritiesStore()
const teamsStore = useTeamsStore()
const usersStore = useUsersStore()
const authStore = useAuthStore()
const rolesStore = useRolesStore()

watch(
  () => authStore.isAuth,
  (isAuth) => {
    if (isAuth) {
      statusesStore.startPolling()
      prioritiesStore.startPolling()
      teamsStore.startPolling()
      usersStore.startPolling()
      incidentsStore.startPolling()
      rolesStore.startPolling()
    } else {
      statusesStore.stopPolling()
      prioritiesStore.stopPolling()
      teamsStore.stopPolling()
      usersStore.stopPolling()
      incidentsStore.stopPolling()
      rolesStore.stopPolling()
    }
  },
  { immediate: true },
)
</script>

<template>
  <el-config-provider>
    <el-container v-if="authStore.initialized" style="height: 100vh; flex-direction: column">
      <Header />
      <el-main>
        <router-view />
      </el-main>
    </el-container>

    <div v-else v-loading.fullscreen.lock="true" element-loading-text="Проверка авторизации..." />
  </el-config-provider>
</template>

<style>
body {
  margin: 0;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  height: 100vh;
}
</style>
