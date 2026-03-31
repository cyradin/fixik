<template>
  <el-tag :style="{ backgroundColor: color, color: '#fff' }">
    {{ label }}
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePrioritiesStore } from '@/stores/prioritiesStore'
import { useIncidentsStore } from '@/stores/incidentsStore'

const props = defineProps<{
  priorityId: number
  label?: string
}>()

// цвета по убыванию критичности
const colors = [
  '#D32F2F', // красный
  '#FBC02D', // желтый
  '#388E3C', // зеленый
  '#455A64', // серый
]

const prioritiesStore = usePrioritiesStore()
const incidentsStore = useIncidentsStore()

// индекс цвета по возрастанию sort
const priorityIndex = computed(() => {
  // сортируем приоритеты по sort
  const sortedPriorities = [...prioritiesStore.items].sort((a, b) => (a.sort ?? 0) - (b.sort ?? 0))
  const idx = sortedPriorities.findIndex(p => p.id === props.priorityId)
  return idx >= 0 ? Math.min(idx, colors.length - 1) : colors.length - 1
})

const color = computed(() => colors[priorityIndex.value])
const label = computed(() => props.label || incidentsStore.items.find(i => i.priority.id === props.priorityId)?.priority.name || '—')
</script>