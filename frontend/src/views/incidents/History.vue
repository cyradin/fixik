<template>
  <div>
    <IncidentFilters :with-status="true" />

    <el-row justify="space-between" align="middle" style="margin-bottom: 16px">
      <h2 style="margin: 0">Инциденты</h2>

      <el-button type="primary" size="large" @click="goToCreate"> + Создать инцидент </el-button>
    </el-row>

    <el-table
      :data="incidentsStore.filteredItems"
      style="width: 100%"
      row-key="id"
      @row-click="goToIncident"
      size="small"
      class="incidents-table"
    >
      <el-table-column prop="id" label="#" width="70" />

      <el-table-column prop="title" label="Инцидент" min-width="240">
        <template #default="{ row }">
          <strong>{{ row.title }}</strong>
        </template>
      </el-table-column>

      <el-table-column label="Статус" width="140">
        <template #default="{ row }">
          {{ row.status?.name }}
        </template>
      </el-table-column>

      <el-table-column label="Приоритет" width="140">
        <template #default="{ row }">
          <IncidentPriorityColor :priority-id="row.priority.id" :label="row.priority.name" />
        </template>
      </el-table-column>

      <el-table-column label="Команда" width="140">
        <template #default="{ row }">
          <span v-if="row.team">{{ row.team.name }}</span>
          <span v-else style="opacity: 0.5">—</span>
        </template>
      </el-table-column>

      <el-table-column label="Исп." width="120">
        <template #default="{ row }">
          <UserLink v-if="row.user" :user="row.user" :is-link="false" />
          <span v-else style="opacity: 0.5">—</span>
        </template>
      </el-table-column>

      <el-table-column label="💬" width="80" align="center">
        <template #default="{ row }">
          <span>{{ row.commentsCount }}</span>
        </template>
      </el-table-column>

      <el-table-column label="Создан" width="160">
        <template #default="{ row }">
          {{ formatDateTime(row.createdAt) }}
        </template>
      </el-table-column>

      <el-table-column label="Обновлён" width="160">
        <template #default="{ row }">
          {{ formatDateTime(row.updatedAt) }}
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { useIncidentsStore } from '@/stores/incidentsStore'
import { useRouter } from 'vue-router'
import IncidentPriorityColor from '@/components/incidents/IncidentPriorityColor.vue'
import UserLink from '@/components/users/UserLink.vue'
import { formatDateTime } from '@/utils/date'
import IncidentFilters from '@/components/incidents/IncidentFilters.vue'

const incidentsStore = useIncidentsStore()
const router = useRouter()

const goToIncident = (row: any) => {
  router.push(`/incident/${row.id}`)
}

const goToCreate = () => {
  router.push('/incident/create')
}
</script>

<style scoped>
.incidents-table :deep(.el-table__row) {
  cursor: pointer;
}

.incidents-table :deep(.el-table__row:hover) {
  background-color: var(--el-fill-color-light);
}
</style>
