<template>
  <el-button link @click="goBack"> ← Назад </el-button>

  <el-card v-if="incident" shadow="hover">
    <h2>#{{ incident.id }} {{ incident.title }}</h2>

    <el-button type="danger" @click="deleteIncident" style="float: right" :loading="loading.delete">
      Удалить
    </el-button>

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
import { reactive, computed, watch, h } from 'vue'
import { useRouter } from 'vue-router'
import IncidentComments from '@/components/incidents/IncidentComments.vue'
import IncidentEditForm from '@/components/incidents/IncidentEditForm.vue'
import { ElMessage } from 'element-plus'

import { useIncidentsStore } from '@/stores/incidentsStore'
import { formatDateTime } from '@/utils/date'
import UserLink from '@/components/users/UserLink.vue'
import { notifyError, notifySuccess } from '@/utils/notify'

const router = useRouter()

const incidentsStore = useIncidentsStore()

const props = defineProps<{
  id: string
}>()

const incident = computed(() => {
  const id = Number(props.id)
  return incidentsStore.items.find((i) => i.id === id) || null
})

const loading = reactive({
  delete: false,
})

let deleteTimer: ReturnType<typeof setTimeout> | null = null
let pendingDeleteId: number | null = null
let pendingBackup: any = null

const cleanup = () => {
  if (deleteTimer) {
    clearTimeout(deleteTimer)
    deleteTimer = null
  }
  pendingDeleteId = null
  pendingBackup = null
}

const deleteIncident = () => {
  if (!incident.value) return

  const id = incident.value.id

  pendingDeleteId = id
  pendingBackup = { ...incident.value }

  incidentsStore.removeLocal(id)

  let msg: any = null

  const undo = () => {
    if (!pendingDeleteId) return

    incidentsStore.addLocal(pendingBackup)
    msg?.close()
    notifySuccess(`Инцидент #${pendingDeleteId} не удалён`, 'Удаление отменено')
    cleanup()
  }

  msg = ElMessage({
    type: 'warning',
    duration: 5000,
    showClose: true,
    message: h('span', [
      `Инцидент #${id} будет удален через 5 секунд`,
      h(
        'span',
        {
          style: 'color: var(--el-color-primary); cursor: pointer; margin-left: 8px;',
          onClick: undo,
        },
        'Отменить',
      ),
    ]),
  })

  deleteTimer = setTimeout(async () => {
    if (!pendingDeleteId) return

    try {
      await incidentsStore.delete(pendingDeleteId)
      notifySuccess(`Инцидент #${pendingDeleteId} удалён окончательно`, 'Инцидент удален')
    } catch (e) {
      incidentsStore.addLocal(pendingBackup)
      notifyError('Ошибка удаления, восстановлено')
    }

    cleanup()
    msg?.close()
  }, 5000)

  router.push('/')
}

const goBack = () => {
  router.push('/')
}
</script>
