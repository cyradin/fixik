<template>
  <el-space wrap>
    <el-button
      v-for="priority in sortedPriorities"
      :key="priority.id"
      size="default"
      :type="store.filters.priorityIds.includes(priority.id) ? 'primary' : 'default'"
      @click="store.togglePriority(priority.id)"
    >
      {{ priority.name }}
    </el-button>

    <el-select
      v-if="withStatus"
      v-model="store.filters.statusIds"
      multiple
      clearable
      collapse-tags
      placeholder="Статус"
      style="min-width: 220px"
    >
      <el-option v-for="s in statuses" :key="s.id" :label="s.name" :value="s.id" />
    </el-select>

    <el-select
      v-model="store.filters.authorIds"
      multiple
      clearable
      collapse-tags
      placeholder="Автор"
      style="min-width: 220px"
    >
      <el-option v-for="u in users" :key="u.id" :label="u.name" :value="u.id" />
    </el-select>

    <el-select
      v-model="store.filters.userIds"
      multiple
      clearable
      collapse-tags
      placeholder="Исполнитель"
      style="min-width: 220px"
    >
      <el-option label="Не назначено" :value="null" />
      <el-option v-for="u in users" :key="u.id" :label="u.name" :value="u.id" />
    </el-select>

    <el-select
      v-model="store.filters.teamIds"
      multiple
      clearable
      collapse-tags
      placeholder="Команда"
      style="min-width: 220px"
    >
      <el-option label="Не назначено" :value="null" />
      <el-option v-for="t in teams" :key="t.id" :label="t.name" :value="t.id" />
    </el-select>

    <el-button size="default" @click="store.resetFilters"> Сбросить </el-button>
  </el-space>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import type { FiltersStoreInterface } from '@/types/filters'
import { usePrioritiesStore } from '@/stores/prioritiesStore'
import { useUsersStore } from '@/stores/usersStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { useStatusesStore } from '@/stores/statusesStore'

const props = defineProps<{
  store: FiltersStoreInterface
  withStatus?: boolean
}>()

const prioritiesStore = usePrioritiesStore()
const usersStore = useUsersStore()
const teamsStore = useTeamsStore()
const statusesStore = useStatusesStore()

onMounted(() => {
  if (!prioritiesStore.items.length) prioritiesStore.fetchAll()
  if (!usersStore.items.length) usersStore.fetchAll()
  if (!teamsStore.items.length) teamsStore.fetchAll()
  if (!statusesStore.items.length) statusesStore.fetchAll()
})

const sortedPriorities = computed(() => [...prioritiesStore.items].sort((a, b) => a.sort - b.sort))
const users = computed(() => usersStore.items)
const teams = computed(() => teamsStore.items)
const statuses = computed(() => statusesStore.items)
</script>
