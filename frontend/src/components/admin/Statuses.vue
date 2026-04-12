<template>
  <EditableAdminTable
    title="Статусы"
    :items="store.items"
    :columns="columns"
    :create="create"
    :update="update"
    :delete="remove"
    @refresh="refresh"
  />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import EditableAdminTable from '@/components/admin/EditTable.vue'
import { useStatusesStore } from '@/stores/statusesStore'
import { statusesApi } from '@/api/client'
import { notifyError, notifySuccess } from '@/utils/notify'

const store = useStatusesStore()

onMounted(() => {
  store.fetchAll()
})

const columns = [
  {
    key: 'id',
    label: 'ID',
    width: '80',
    editor: 'span',
  },
  {
    key: 'code',
    label: 'Code',
    editor: 'el-input',
  },
  {
    key: 'name',
    label: 'Название',
    editor: 'el-input',
  },
  {
    key: 'description',
    label: 'Описание',
    editor: 'el-input',
  },
  {
    key: 'sort',
    label: 'Sort',
    width: '100',
    editor: 'el-input-number',
    editorProps: {
      min: 0,
    },
  },
  {
    key: 'isFinal',
    label: 'Финальный',
    width: '120',
    editor: 'el-switch',
    formatter: (v: boolean) => (v ? 'Да' : 'Нет'),
  },
]

const create = async (data: any) => {
  if (!data.code || !data.name) {
    notifyError('Заполните code и name')
    throw new Error('validation')
  }

  const res = await statusesApi.statusesPost({
    request: data,
  })

  notifySuccess('Статус создан')
  return res
}

const update = async (id: number, data: any) => {
  await statusesApi.statusesIdPut({
    id,
    request: data,
  })

  notifySuccess('Сохранено')
}

const remove = async (id: number) => {
  await statusesApi.statusesIdDelete({ id })
  notifySuccess('Удалено')
}

const refresh = () => {
  store.fetchAll()
}
</script>
