<template>
  <el-button link @click="goBack"> ← Назад </el-button>

  <el-card v-if="loading" shadow="never"> Загрузка... </el-card>

  <el-card v-else-if="incident" shadow="hover" style="width: 100%">
    <el-space style="width: 100%" alignment="center" justify="space-between">
      <h2 style="margin: 0">#{{ incident.id }} {{ incident.title }}</h2>

      <IncidentDeleteButton :id="Number(id)" v-if="can(PERMISSIONS.INCIDENT_DELETE)" />
    </el-space>

    <p v-if="incident.author">
      <b>Автор:</b>
      <UserLink :user="incident.author" :is-link="true" />
    </p>

    <p v-if="incident.user">
      <b>Исполнитель:</b>
      <UserLink :user="incident.user" :is-link="true" />
    </p>

    <p><b>Создан:</b> {{ formatDateTime(incident.createdAt) }}</p>
    <p><b>Обновлён:</b> {{ formatDateTime(incident.updatedAt) }}</p>

    <IncidentEditForm :incident="incident" />

    <el-row>
      <IncidentComments :incident-id="Number(id)" />
    </el-row>
  </el-card>

  <el-empty v-else description="Инцидент не найден" />
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import IncidentComments from '@/components/incidents/IncidentComments.vue'
import IncidentDeleteButton from '@/components/incidents/IncidentDeleteButton.vue'
import IncidentEditForm from '@/components/incidents/IncidentEditForm.vue'

import { Incident, useIncidentsStore } from '@/stores/incidentsStore'
import { formatDateTime } from '@/utils/date'
import UserLink from '@/components/users/UserLink.vue'
import { useAuthStore } from '@/stores/authStore'
import { PERMISSIONS } from '@/constants/permissions'
import { notifyError } from '@/utils/notify'

const router = useRouter()
const incident = ref<Incident | null>(null)
const loading = ref(false)
const incidentsStore = useIncidentsStore()
const authStore = useAuthStore()

const can = authStore.can

const props = defineProps<{
  id: string
}>()

const load = async (id: number) => {
  loading.value = true
  try {
    incident.value = await incidentsStore.getById(id)
  } catch (e: any) {
    notifyError(e.message)
    incident.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  load(Number(props.id))
})

watch(
  () => props.id,
  (newId) => {
    load(Number(newId))
  },
)

const goBack = () => {
  router.back()
}
</script>
