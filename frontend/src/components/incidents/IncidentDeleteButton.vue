<template>
  <el-button type="danger" :loading="loading" @click="deleteIncident"> Удалить </el-button>
</template>

<script setup lang="ts">
import { ref, h, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { useIncidentsStore } from '@/stores/incidentsStore'
import { notifyError, notifySuccess } from '@/utils/notify'

const props = defineProps<{
  id: number
}>()

const incident = computed(() => {
  const id = Number(props.id)
  return incidentsStore.items.find((i) => i.id === id) || null
})

const router = useRouter()
const incidentsStore = useIncidentsStore()

const loading = ref(false)

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
</script>
