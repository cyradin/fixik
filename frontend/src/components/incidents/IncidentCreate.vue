<template>
  <el-button link @click="goBack"> ← Назад </el-button>

  <el-card shadow="hover">
    <el-form label-position="top">
      <el-form-item label="Название">
        <el-input v-model="model.title" placeholder="Введите название" />
      </el-form-item>

      <el-form-item label="Описание">
        <el-input type="textarea" v-model="model.description" placeholder="Введите описание" />
      </el-form-item>

      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="Статус">
            <el-select v-model="model.statusId" placeholder="Выберите статус">
              <el-option v-for="s in statuses" :key="s.id" :label="s.name" :value="s.id" />
            </el-select>
          </el-form-item>

          <el-form-item label="Команда">
            <el-select v-model="model.teamId" clearable placeholder="Выберите команду">
              <el-option v-for="t in teams" :key="t.id" :label="t.name" :value="t.id" />
            </el-select>
          </el-form-item>
        </el-col>

        <el-col :span="12">
          <el-form-item label="Приоритет">
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
import { ElMessage } from 'element-plus'
import { useIncidentsStore, Incident } from '@/stores/incidentsStore'
import { useStatusesStore } from '@/stores/statusesStore'
import { useUsersStore } from '@/stores/usersStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { usePrioritiesStore } from '@/stores/prioritiesStore'

const router = useRouter()
const incidentsStore = useIncidentsStore()
const statusesStore = useStatusesStore()
const prioritiesStore = usePrioritiesStore()
const teamsStore = useTeamsStore()
const usersStore = useUsersStore()

const model = reactive({
  title: '',
  description: '',
  statusId: null as number | null,
  priorityId: null as number | null,
  teamId: null as number | null,
  userId: null as number | null,
})

const statuses = computed(() => statusesStore.items)
const priorities = computed(() => prioritiesStore.items)
const teams = computed(() => teamsStore.items)
const users = computed(() => usersStore.items)

const loading = ref(false)

const submit = async () => {
  if (!model.title || !model.statusId || !model.priorityId || !model.description) {
    ElMessage({ message: 'Заполните название, описание статус и приоритет', type: 'warning' })
    return
  }

  loading.value = true
  try {
    const status = statusesStore.items.find((s) => s.id === model.statusId)
    const priority = prioritiesStore.items.find((p) => p.id === model.priorityId)

    if (!status || !priority) {
      ElMessage({ message: 'Выберите статус и приоритет', type: 'warning' })
      return
    }

    let userId
    if (model.userId) {
      userId = model.userId
    }

    let teamId
    if (model.teamId) {
      teamId = model.teamId
    }

    const newIncident = await incidentsStore.create({
      title: model.title,
      description: model.description,
      statusId: status.id,
      priorityId: priority.id,
      userId: userId,
      teamId: teamId,
    })
    ElMessage({ message: `Инцидент #${newIncident.id} создан`, type: 'success' })
    router.push(`/incident/${newIncident.id}`)
  } catch (e) {
    ElMessage({ message: 'Не удалось создать инцидент', type: 'error' })
    console.error(e)
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push('/')
}
</script>
