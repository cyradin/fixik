<template>
  <el-button link @click="goBack"> ← Назад </el-button>
  <el-card v-if="incident" shadow="hover">
    <h2>#{{ incident.id }} {{ incident.title }}</h2>

    <p v-if="incident.author"><b>Автор:</b> <UserLink :user="incident.author" /></p>

    <p><b>Описание:</b> {{ incident.description }}</p>

    <p><b>Статус:</b> {{ incident.status.name }}</p>
    <p>
      <b>Приоритет:</b>
      <IncidentPriorityColor :priority-id="incident.priority.id" :label="incident.priority.name" />
    </p>

    <p v-if="incident.team"><b>Команда:</b> {{ incident.team.name }}</p>
    <p v-if="incident.user"><b>Исполнитель:</b> <UserLink :user="incident.user" /></p>
  </el-card>

  <el-empty v-else description="Инцидент не найден" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useIncidentsStore } from '@/stores/incidentsStore'
import IncidentPriorityColor from '@/components/incidents/IncidentPriorityColor.vue'
import UserLink from '@/components/users/UserLink.vue'

const route = useRoute()
const router = useRouter()
const incidentsStore = useIncidentsStore()

const incident = computed(() => {
  const id = Number(route.params.id)
  return incidentsStore.items.find((i) => i.id === id) || null
})

const goBack = () => {
  router.back()
}
</script>
