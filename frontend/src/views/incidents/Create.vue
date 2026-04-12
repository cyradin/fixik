<template>
  <el-button link @click="goBack"> ← Назад </el-button>

  <el-card shadow="hover">
    <el-form ref="formRef" :model="model" :rules="rules" label-position="top">
      <el-form-item label="Название" prop="title">
        <el-input v-model="model.title" placeholder="Введите название" />
      </el-form-item>

      <el-form-item label="Описание" prop="description">
        <el-input type="textarea" v-model="model.description" placeholder="Введите описание" />
      </el-form-item>

      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="Статус">
            <el-input :model-value="defaultStatus?.name || ''" disabled />
          </el-form-item>

          <el-form-item label="Команда">
            <el-select v-model="model.teamId" clearable placeholder="Выберите команду">
              <el-option v-for="t in teams" :key="t.id" :label="t.name" :value="t.id" />
            </el-select>
          </el-form-item>
        </el-col>

        <el-col :span="12">
          <el-form-item label="Приоритет" prop="priorityId">
            <el-select v-model="model.priorityId" placeholder="Выберите приоритет">
              <el-option v-for="p in priorities" :key="p.id" :label="p.name" :value="p.id" />
            </el-select>
          </el-form-item>

          <el-form-item label="Исполнитель">
            <el-select v-model="model.userId" clearable placeholder="Выберите исполнителя">
              <el-option v-for="u in users" :key="u.id" :label="u.name" :value="u.id" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-button type="primary" :loading="loading" @click="submit"> Создать </el-button>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useIncidentsStore } from '@/stores/incidentsStore'
import { useStatusesStore } from '@/stores/statusesStore'
import { useUsersStore } from '@/stores/usersStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { usePrioritiesStore } from '@/stores/prioritiesStore'
import { notifyError, notifySuccess } from '@/utils/notify'
import { useAuthStore } from '@/stores/authStore'

const router = useRouter()
const authStore = useAuthStore()
const incidentsStore = useIncidentsStore()
const statusesStore = useStatusesStore()
const prioritiesStore = usePrioritiesStore()
const teamsStore = useTeamsStore()
const usersStore = useUsersStore()

const formRef = ref()

const model = reactive({
  title: '',
  description: '',
  statusId: null as number | null,
  priorityId: null as number | null,
  teamId: null as number | null,
  userId: null as number | null,
})

const rules = {
  title: [{ required: true, message: 'Введите название', trigger: 'blur' }],
  description: [{ required: true, message: 'Введите описание', trigger: 'blur' }],
  priorityId: [{ required: true, message: 'Выберите приоритет', trigger: 'change' }],
}

const statuses = computed(() => statusesStore.items)
const priorities = computed(() => prioritiesStore.items)
const teams = computed(() => teamsStore.items)
const users = computed(() => usersStore.items)

const defaultStatus = computed(() => {
  if (!statuses.value.length) return null
  return [...statuses.value].sort((a: any, b: any) => (a.sort ?? 0) - (b.sort ?? 0))[0] ?? null
})

if (!model.statusId && defaultStatus.value) {
  model.statusId = defaultStatus.value.id
}

const loading = ref(false)

const submit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const status = defaultStatus.value
    const priority = prioritiesStore.items.find((p) => p.id === model.priorityId)

    const newIncident = await incidentsStore.create({
      title: model.title,
      description: model.description,
      statusId: status!.id,
      priorityId: priority!.id,
      userId: model.userId || undefined,
      teamId: model.teamId || undefined,
      authorId: authStore.user?.id,
    })

    notifySuccess(`Инцидент #${newIncident.id} создан`)

    router.push(`/incident/${newIncident.id}`)
  } catch (e) {
    notifyError('Не удалось создать инцидент')
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push('/')
}
</script>
