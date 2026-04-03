<template>
  <el-space wrap>
    <el-button
      v-for="priority in sortedPriorities"
      :key="priority.id"
      size="default"
      :type="incidentsStore.filters.priorityIds.includes(priority.id) ? 'primary' : 'default'"
      @click="incidentsStore.togglePriority(priority.id)"
    >
      {{ priority.name }}
    </el-button>

    <el-button
      size="default"
      :type="incidentsStore.filters.unassignedOnly ? 'danger' : 'default'"
      @click="incidentsStore.toggleUnassigned"
    >
      Не назначенные
    </el-button>

    <el-button size="default" @click="incidentsStore.resetFilters"> Сбросить </el-button>
  </el-space>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useIncidentsStore } from '@/stores/incidentsStore'
import { usePrioritiesStore } from '@/stores/prioritiesStore'

const incidentsStore = useIncidentsStore()
const prioritiesStore = usePrioritiesStore()

onMounted(() => {
  if (prioritiesStore.items.length === 0) {
    prioritiesStore.fetchAll()
  }
})

const sortedPriorities = computed(() => {
  return [...prioritiesStore.items].sort((a, b) => a.sort - b.sort)
})
</script>
