<template>
  <el-button type="danger" :loading="loading" @click="deleteIncident"> Удалить </el-button>
</template>

<script setup lang="ts">
import { ref, h, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { useIncidentsStore } from '@/stores/incidentsStore'
import { notifyError, notifySuccess } from '@/utils/notify'

const props = defineProps<{ id: number }>()

const router = useRouter()
const incidentsStore = useIncidentsStore()

const loading = ref(false)

const incident = computed(() => {
  return incidentsStore.items.find((i) => i.id === props.id) || null
})

const deleteIncident = () => {
  if (!incident.value) return

  const id = incident.value.id

  incidentsStore.delete(id)

  let msg: any = null

  const undo = () => {
    incidentsStore.undoDelete(id)
    msg?.close()
    notifySuccess(`Инцидент #${id} не удалён`, 'Удаление отменено')
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

  router.push('/')
}
</script>
