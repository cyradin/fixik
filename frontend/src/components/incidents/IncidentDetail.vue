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

    <el-input
      type="textarea"
      v-model="editable.description"
      :disabled="loading.description"
      @blur="updateField('description', editable.description)"
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
          placeholder="Выберите статус"
          style="width: 100%"
          @change="updateField('status', editable.statusId)"
        >
          <el-option v-for="s in statuses" :key="s.id" :label="s.name" :value="s.id" />
        </el-select>
        <el-icon v-if="loading.status" class="is-loading" style="margin-left: 8px">
          <Loading />
        </el-icon>

        <p style="margin-top: 16px"><b>Команда:</b></p>
        <el-select
          v-model="editable.teamId"
          clearable
          :disabled="loading.team"
          placeholder="Выберите команду"
          style="width: 100%"
          @change="updateField('team', editable.teamId)"
        >
          <el-option v-for="t in teams" :key="t.id" :label="t.name" :value="t.id" />
        </el-select>
        <el-icon v-if="loading.team" class="is-loading" style="margin-left: 8px">
          <Loading />
        </el-icon>
      </el-col>

      <el-col :span="12">
        <p><b>Приоритет:</b></p>
        <el-select
          v-model="editable.priorityId"
          :disabled="loading.priority"
          placeholder="Выберите приоритет"
          style="width: 100%"
          @change="updateField('priority', editable.priorityId)"
        >
          <el-option v-for="p in priorities" :key="p.id" :label="p.name" :value="p.id" />
        </el-select>
        <el-icon v-if="loading.priority" class="is-loading" style="margin-left: 8px">
          <Loading />
        </el-icon>

        <p style="margin-top: 16px"><b>Исполнитель:</b></p>
        <el-select
          v-model="editable.userId"
          clearable
          :disabled="loading.user"
          placeholder="Выберите исполнителя"
          style="width: 100%"
          @change="updateField('user', editable.userId)"
        >
          <el-option v-for="u in users" :key="u.id" :label="u.name" :value="u.id" />
        </el-select>
        <el-icon v-if="loading.user" class="is-loading" style="margin-left: 8px">
          <Loading />
        </el-icon>
      </el-col>
    </el-row>
  </el-card>

  <el-empty v-else description="Инцидент не найден" />
</template>

<script setup lang="ts">
import { reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { useIncidentsStore } from '@/stores/incidentsStore'
import { useStatusesStore } from '@/stores/statusesStore'
import { useUsersStore } from '@/stores/usersStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { usePrioritiesStore } from '@/stores/prioritiesStore'
import { formatDateTime } from '@/utils/date'
import { watch } from 'vue'
import UserLink from '@/components/users/UserLink.vue'

const router = useRouter()
const incidentsStore = useIncidentsStore()
const statusesStore = useStatusesStore()
const prioritiesStore = usePrioritiesStore()
const teamsStore = useTeamsStore()
const usersStore = useUsersStore()
const props = defineProps<{
  id: string
}>()

const incident = computed(() => {
  const id = Number(props.id)
  return incidentsStore.items.find((i) => i.id === id) || null
})

const editable = reactive({
  description: incident.value?.description || '',
  statusId: incident.value?.status.id || null,
  priorityId: incident.value?.priority.id || null,
  teamId: incident.value?.team?.id || null,
  userId: incident.value?.user?.id || null,
})

watch(
  incident,
  (inc) => {
    if (!inc) return
    editable.description = inc.description
    editable.statusId = inc.status?.id ?? null
    editable.priorityId = inc.priority?.id ?? null
    editable.teamId = inc.team?.id ?? null
    editable.userId = inc.user?.id ?? null
  },
  { immediate: true },
)

const loading = reactive({
  description: false,
  status: false,
  priority: false,
  team: false,
  user: false,
  delete: false,
})

const statuses = computed(() => statusesStore.items)
const priorities = computed(() => prioritiesStore.items)
const teams = computed(() => teamsStore.items)
const users = computed(() => {
  const userStore = usersStore
  const teamId = editable.teamId
  return userStore.byTeam(teamId ?? undefined)
})

type Field = keyof typeof loading

const updateField = async (field: Field, value: any) => {
  if (!incident.value) return

  loading[field] = true
  try {
    switch (field) {
      case 'description':
        await incidentsStore.updateDescription(incident.value.id, value)
        incident.value.description = value
        break
      case 'status':
        await incidentsStore.updateStatus(incident.value.id, value)
        const status = statusesStore.items.find((s) => s.id === value)
        if (status) incident.value.status = { ...status }
        break
      case 'priority':
        await incidentsStore.updatePriority(incident.value.id, value)
        const priority = prioritiesStore.items.find((p) => p.id === value)
        if (priority) incident.value.priority = { ...priority }
        break
      case 'team':
        await incidentsStore.updateTeam(incident.value.id, value ?? undefined)
        const team = teamsStore.items.find((t) => t.id === value)
        incident.value.team = team ? { ...team } : null
        break

      case 'user':
        await incidentsStore.updateUser(incident.value.id, value ?? undefined)
        const user = usersStore.items.find((u) => u.id === value)
        incident.value.user = user ? { ...user } : null
        break
    }
  } catch (e) {
    ElMessage({
      message: `Не удалось обновить ${field} инцидента #${incident.value.id}`,
      type: 'error',
      duration: 3000,
    })
    console.error(`Ошибка при обновлении ${field}`, e)
  } finally {
    loading[field] = false
  }
}

const deleteIncident = async () => {
  if (!incident.value) return
  if (!confirm(`Вы уверены, что хотите удалить инцидент #${incident.value.id}?`)) return

  loading.delete = true
  const id = incident.value.id

  try {
    await incidentsStore.delete(id) // нужно добавить delete в стор
    ElMessage({
      message: `Инцидент #${id} удалён`,
      type: 'success',
      duration: 3000,
    })
    router.push('/')
  } catch (e) {
    ElMessage({
      message: `Не удалось удалить инцидент #${id}`,
      type: 'error',
      duration: 3000,
    })
    console.error('Ошибка удаления инцидента', e)
  } finally {
    loading.delete = false
  }
}

const goBack = () => {
  router.push('/')
}
</script>
