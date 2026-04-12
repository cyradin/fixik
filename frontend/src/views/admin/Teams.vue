<template>
  <EditableAdminTable
    title="Команды"
    :items="store.items"
    :columns="columns"
    :create="create"
    :update="update"
    :remove="remove"
    :default-row="getDefaultRow"
    @refresh="refresh"
  />
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import EditableAdminTable from '@/components/admin/EditTable.vue'
import { useTeamsStore } from '@/stores/teamsStore'
import { notifyError, notifySuccess } from '@/utils/notify'

const store = useTeamsStore()

onMounted(() => {
  store.fetchAll()
  store.stopPolling()
})

onUnmounted(() => {
  store.startPolling()
})

const columns = [
  { key: 'id', label: 'ID', width: '80', editor: 'span' },

  { key: 'code', label: 'Code', editor: 'el-input', required: true },

  { key: 'name', label: 'Название', editor: 'el-input', required: true },

  { key: 'description', label: 'Описание', editor: 'el-input', required: true },

  {
    key: 'sort',
    label: 'Sort',
    width: '160',
    editor: 'el-input-number',
    required: true,
    editorProps: { min: 0 },
  },
]

const getDefaultSort = () => {
  if (!store.items.length) return 10
  const max = Math.max(...store.items.map((i) => i.sort || 0))
  return max + 10
}

const getDefaultRow = () => ({
  code: '',
  name: '',
  description: '',
  sort: getDefaultSort(),
})

const create = async (data: any) => {
  try {
    const res = await store.create(data)
    notifySuccess('Команда создана')
    return res
  } catch (e) {
    notifyError('Ошибка создания команды')
    throw e
  }
}

const update = async (id: number, data: any) => {
  try {
    await store.update(id, data)
    notifySuccess('Команда обновлена')
  } catch (e) {
    notifyError('Ошибка обновления команды')
    throw e
  }
}

const remove = async (id: number) => {
  try {
    await store.remove(id)
    notifySuccess('Команда удалена')
  } catch (e) {
    notifyError('Ошибка удаления команды')
    throw e
  }
}

const refresh = () => {
  store.fetchAll()
}
</script>
