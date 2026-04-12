<template>
  <EditableAdminTable
    title="Пользователи"
    :items="store.items"
    :columns="columns"
    :create="create"
    :update="update"
    :remove="remove"
    :default-row="getDefaultRow"
    :on-edit-row="onEditRow"
    @refresh="refresh"
  />
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import EditableAdminTable from '@/components/admin/EditTable.vue'

import { useUsersStore } from '@/stores/usersStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { useRolesStore } from '@/stores/rolesStore'
import { notifyError, notifySuccess } from '@/utils/notify'

const usersStore = useUsersStore()
const teamsStore = useTeamsStore()
const rolesStore = useRolesStore()

onMounted(() => {
  usersStore.fetchAll()
  teamsStore.fetchAll()
  rolesStore.fetchAll()

  usersStore.stopPolling()
  teamsStore.stopPolling()
  rolesStore.stopPolling()
})

onUnmounted(() => {
  usersStore.startPolling()
  teamsStore.startPolling()
  rolesStore.startPolling()
})

const onEditRow = (row: any, formRow: any) => {
  if ('password' in formRow) {
    formRow.password = ''
  }
}

const store = usersStore

const columns = computed(() => [
  { key: 'id', label: 'ID', width: '80', editor: 'span' },

  { key: 'name', label: 'Имя', editor: 'el-input', required: true },
  { key: 'username', label: 'Username', editor: 'el-input', required: true },
  { key: 'email', label: 'Email', editor: 'el-input', required: true },

  {
    key: 'password',
    label: 'Пароль',
    editor: 'el-input',
    editorProps: {
      type: 'password',
      showPassword: true,
    },
    validation: {
      min: 6,
    },
  },

  {
    key: 'role',
    label: 'Роль',
    editor: 'el-select',
    required: true,
    options: rolesStore.items.map((r) => ({
      label: r.name,
      value: r.code,
    })),
  },

  {
    key: 'teamId',
    label: 'Команда',
    editor: 'el-select',
    options: teamsStore.items.map((t) => ({
      label: t.name,
      value: t.id,
    })),
  },
])

const getDefaultRow = () => ({
  name: '',
  username: '',
  email: '',
  role: '',
  teamId: null,
  password: '',
})

const create = async (data: any) => {
  try {
    const res = await store.create(data)
    notifySuccess('Пользователь создан')
    return res
  } catch (e) {
    notifyError('Ошибка создания пользователя')
    throw e
  }
}

const update = async (id: number, data: any) => {
  try {
    await store.update(id, data)
    notifySuccess('Пользователь обновлён')
  } catch (e) {
    notifyError('Ошибка обновления пользователя')
    throw e
  }
}

const remove = async (id: number) => {
  try {
    await store.remove(id)
    notifySuccess('Пользователь удален')
  } catch (e) {
    notifyError('Ошибка удаления пользователя')
    throw e
  }
}

const refresh = () => {
  store.fetchAll()
}
</script>
