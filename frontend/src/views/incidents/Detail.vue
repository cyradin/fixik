<template>
  <el-button link @click="goBack"> ← Назад </el-button>

  <el-card v-if="incident" shadow="hover" style="width: 100%">
    <el-space style="width: 100%" alignment="center" justify="space-between">
      <h2 style="margin: 0">#{{ incident.id }} {{ incident.title }}</h2>

      <IncidentDeleteButton :id="Number(id)" />
    </el-space>

    <p v-if="incident.author">
      <b>Автор:</b>
      <UserLink :user="incident.author" :is-link="true" />
    </p>

    <p><b>Создан:</b> {{ formatDateTime(incident.createdAt) }}</p>
    <p><b>Обновлён:</b> {{ formatDateTime(incident.updatedAt) }}</p>

    <IncidentEditForm :id="id" />

    <el-row>
      <IncidentComments :incident-id="Number(id)" />
    </el-row>
  </el-card>

  <el-empty v-else description="Инцидент не найден" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import IncidentComments from '@/components/incidents/IncidentComments.vue'
import IncidentDeleteButton from '@/components/incidents/IncidentDeleteButton.vue'
import IncidentEditForm from '@/components/incidents/IncidentEditForm.vue'

import { useIncidentsStore } from '@/stores/incidentsStore'
import { formatDateTime } from '@/utils/date'
import UserLink from '@/components/users/UserLink.vue'

const router = useRouter()

const incidentsStore = useIncidentsStore()

const props = defineProps<{
  id: string
}>()

const incident = computed(() => {
  const id = Number(props.id)
  return incidentsStore.items.find((i) => i.id === id) || null
})

const goBack = () => {
  router.push('/')
}
</script>
