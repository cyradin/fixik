<template>
  <div>
    <el-input
      type="textarea"
      v-model="editable.description"
      :disabled="loading.description"
      @blur="updateDescription(editable.description)"
      placeholder="Введите описание"
    />
    <el-icon v-if="loading.description" class="is-loading" style="margin-left: 8px">
      <Loading />
    </el-icon>

    <el-row :gutter="16">
      <el-col :span="12">
        <p><b>Статус:</b></p>
        <el-select
          v-model="editable.statusId"
          :disabled="loading.status"
          style="width: 100%"
          @change="updateStatus(editable.statusId)"
        >
          <el-option v-for="s in statuses" :key="s.id" :label="s.name" :value="s.id" />
        </el-select>

        <p style="margin-top: 16px"><b>Команда:</b></p>
        <el-select
          v-model="editable.teamId"
          clearable
          :disabled="loading.team"
          style="width: 100%"
          @change="updateTeam(editable.teamId)"
        >
          <el-option v-for="t in teams" :key="t.id" :label="t.name" :value="t.id" />
        </el-select>
      </el-col>

      <el-col :span="12">
        <p><b>Приоритет:</b></p>
        <el-select
          v-model="editable.priorityId"
          :disabled="loading.priority"
          style="width: 100%"
          @change="updatePriority(editable.priorityId)"
        >
          <el-option v-for="p in priorities" :key="p.id" :label="p.name" :value="p.id" />
        </el-select>

        <p style="margin-top: 16px"><b>Исполнитель:</b></p>
        <el-select
          v-model="editable.userId"
          clearable
          :disabled="loading.user"
          style="width: 100%"
          @change="updateUser(editable.userId)"
        >
          <el-option v-for="u in users" :key="u.id" :label="u.name" :value="u.id" />
        </el-select>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { reactive, watch, computed } from 'vue'
import { Loading } from '@element-plus/icons-vue'
import { notifyError, notifySuccess } from '@/utils/notify'

import { Incident, useIncidentsStore } from '@/stores/incidentsStore'
import { useStatusesStore } from '@/stores/statusesStore'
import { useUsersStore } from '@/stores/usersStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { usePrioritiesStore } from '@/stores/prioritiesStore'

const props = defineProps<{
  incident: Incident
}>()

const incidentsStore = useIncidentsStore()
const statusesStore = useStatusesStore()
const prioritiesStore = usePrioritiesStore()
const teamsStore = useTeamsStore()
const usersStore = useUsersStore()

const statuses = computed(() => statusesStore.items)
const priorities = computed(() => prioritiesStore.items)
const teams = computed(() => teamsStore.items)

const users = computed(() => {
  const teamId = editable.teamId
  return usersStore.byTeam(teamId ?? undefined)
})

const loading = reactive({
  description: false,
  status: false,
  priority: false,
  team: false,
  user: false,
})

const editable = reactive({
  description: '',
  statusId: 0,
  priorityId: 0,
  teamId: null as number | null,
  userId: null as number | null,
})

watch(
  () => props.incident,
  (inc) => {
    if (!inc) return

    editable.description = inc.description
    editable.statusId = inc.status?.id ?? 0
    editable.priorityId = inc.priority?.id ?? 0
    editable.teamId = inc.team?.id ?? null
    editable.userId = inc.user?.id ?? null
  },
  { immediate: true },
)

const updateDescription = async (value: string) => {
  const inc = props.incident
  if (!inc) return

  loading.description = true
  try {
    await incidentsStore.updateDescription(inc.id, value)

    inc.description = value

    notifySuccess('Описание обновлено')
  } catch (e: any) {
    notifyError(e.message)
  } finally {
    loading.description = false
  }
}

const updateStatus = async (value: number) => {
  const inc = props.incident
  if (!inc) return

  loading.status = true
  try {
    await incidentsStore.updateStatus(inc.id, value)

    const status = statusesStore.items.find((s) => s.id === value)
    if (status) inc.status = { ...status }

    notifySuccess('Статус обновлен')
  } catch (e: any) {
    notifyError(e.message)
  } finally {
    loading.status = false
  }
}

const updatePriority = async (value: number) => {
  const inc = props.incident
  if (!inc) return

  loading.priority = true
  try {
    await incidentsStore.updatePriority(inc.id, value)

    const priority = prioritiesStore.items.find((p) => p.id === value)
    if (priority) inc.priority = { ...priority }

    notifySuccess('Приоритет обновлен')
  } catch (e: any) {
    notifyError(e.message)
  } finally {
    loading.priority = false
  }
}

const updateTeam = async (value: number | null) => {
  const inc = props.incident
  if (!inc) return

  loading.team = true
  try {
    await incidentsStore.updateTeam(inc.id, value ?? undefined)

    const team = teamsStore.items.find((t) => t.id === value)
    inc.team = team ? { ...team } : null

    notifySuccess('Команда обновлена')
  } catch (e: any) {
    notifyError(e.message)
  } finally {
    loading.team = false
  }
}

const updateUser = async (value: number | null) => {
  const inc = props.incident
  if (!inc) return

  loading.user = true
  try {
    await incidentsStore.updateUser(inc.id, value ?? undefined)

    const user = usersStore.items.find((u) => u.id === value)
    inc.user = user ? { ...user } : null

    notifySuccess('Исполнитель обновлен')
  } catch (e: any) {
    notifyError(e.message)
  } finally {
    loading.user = false
  }
}
</script>
