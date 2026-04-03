<script setup lang="ts">
import { watch } from 'vue'

import { useIncidentsStore } from '@/stores/incidentsStore'
import { useStatusesStore } from '@/stores/statusesStore'
import { usePrioritiesStore } from '@/stores/prioritiesStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { useUsersStore } from '@/stores/usersStore'
import { useAuthStore } from '@/stores/authStore'

import Header from '@/components/layout/Header.vue'

const incidentsStore = useIncidentsStore()
const statusesStore = useStatusesStore()
const prioritiesStore = usePrioritiesStore()
const teamsStore = useTeamsStore()
const usersStore = useUsersStore()
const authStore = useAuthStore()

watch(
  () => authStore.isAuth,
  (isAuth) => {
    if (isAuth) {
      statusesStore.fetchAll()
      prioritiesStore.fetchAll()
      teamsStore.fetchAll()
      usersStore.fetchAll()

      incidentsStore.startPolling(5000)
      usersStore.startPolling(5000)
    } else {
      incidentsStore.stopPolling()
      usersStore.stopPolling()
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
