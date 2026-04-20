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

    <el-select
      v-if="props.withStatus"
      v-model="incidentsStore.filters.statusIds"
      multiple
      clearable
      collapse-tags
      placeholder="Статус"
      style="min-width: 220px"
    >
      <el-option v-for="s in statusesStore.items" :key="s.id" :label="s.name" :value="s.id" />
    </el-select>

    <el-select
      v-model="incidentsStore.filters.authorIds"
      multiple
      clearable
      collapse-tags
      placeholder="Автор"
      style="min-width: 220px"
    >
      <el-option v-for="u in usersStore.items" :key="u.id" :label="u.name" :value="u.id" />
    </el-select>

    <el-select
      v-model="incidentsStore.filters.userIds"
      multiple
      clearable
      collapse-tags
      placeholder="Исполнитель"
      style="min-width: 220px"
    >
      <el-option label="Не назначено" :value="null" />

      <el-option v-for="u in usersStore.items" :key="u.id" :label="u.name" :value="u.id" />
    </el-select>

    <el-select
      v-model="incidentsStore.filters.teamIds"
      multiple
      clearable
      collapse-tags
      placeholder="Команда"
      style="min-width: 220px"
    >
      <el-option label="Не назначено" :value="null" />

      <el-option v-for="t in teamsStore.items" :key="t.id" :label="t.name" :value="t.id" />
    </el-select>

    <el-button size="default" @click="incidentsStore.resetFilters"> Сбросить </el-button>
  </el-space>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useIncidentsStore } from '@/stores/incidentsStore'
import { usePrioritiesStore } from '@/stores/prioritiesStore'
import { useUsersStore } from '@/stores/usersStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { useStatusesStore } from '@/stores/statusesStore'

const props = defineProps<{
  withStatus?: boolean
}>()

const usersStore = useUsersStore()
const teamsStore = useTeamsStore()
const incidentsStore = useIncidentsStore()
const prioritiesStore = usePrioritiesStore()
const statusesStore = useStatusesStore()

onMounted(() => {
  if (prioritiesStore.items.length === 0) {
    prioritiesStore.fetchAll()
  }
})

const sortedPriorities = computed(() => {
  return [...prioritiesStore.items].sort((a, b) => a.sort - b.sort)
})
</script>
