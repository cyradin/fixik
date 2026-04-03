<script setup lang="ts">
import { onUnmounted, watchEffect } from 'vue'

import { useIncidentsStore } from '@/stores/incidentsStore'
import { useStatusesStore } from '@/stores/statusesStore'
import { usePrioritiesStore } from '@/stores/prioritiesStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { useUsersStore } from '@/stores/usersStore'
import { useAuthStore } from '@/stores/authStore'

const incidentsStore = useIncidentsStore()
const statusesStore = useStatusesStore()
const prioritiesStore = usePrioritiesStore()
const teamsStore = useTeamsStore()
const usersStore = useUsersStore()
const authStore = useAuthStore()

let stopPolling = false

watchEffect(() => {
  if (authStore.isAuth && !stopPolling) {
    statusesStore.fetchAll()
    prioritiesStore.fetchAll()
    teamsStore.fetchAll()
    usersStore.fetchAll()

    incidentsStore.startPolling(5000)
    usersStore.startPolling(5000)
    stopPolling = true
  }
})

onUnmounted(() => {
  incidentsStore.stopPolling()
  usersStore.stopPolling()
})
</script>

<template>
  <el-config-provider>
    <el-container v-if="authStore.initialized" style="padding: 20px">
      <el-header>
        <h1>Fixik: Система инцидентов</h1>
      </el-header>

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
}
</style>
