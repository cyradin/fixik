<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'

import { useIncidentsStore } from '@/stores/incidentsStore'
import { useStatusesStore } from '@/stores/statusesStore'
import { usePrioritiesStore } from '@/stores/prioritiesStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { useUsersStore } from '@/stores/usersStore'

const incidentsStore = useIncidentsStore()
const statusesStore = useStatusesStore()
const prioritiesStore = usePrioritiesStore()
const teamsStore = useTeamsStore()
const usersStore = useUsersStore()

onMounted(() => {
  statusesStore.fetchAll()
  prioritiesStore.fetchAll()
  teamsStore.fetchAll()
  usersStore.fetchAll()

  incidentsStore.startPolling(5000)
  usersStore.startPolling(5000)
})

onUnmounted(() => {
  incidentsStore.stopPolling()
  usersStore.stopPolling()
})
</script>

<template>
  <el-container style="padding: 20px">
    <el-header>
      <h1>Fixik: Система инцидентов</h1>
    </el-header>

    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<style>
body {
  margin: 0;
  font-family: 'Helvetica Neue', Arial, sans-serif;
}
</style>
