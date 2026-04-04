<template>
  <div>
    <IncidentFilters />

    <el-row justify="space-between" align="middle" style="margin-bottom: 16px">
      <h2 style="margin: 0">Инциденты</h2>

      <el-button type="primary" size="large" @click="goToCreate"> + Создать инцидент </el-button>
    </el-row>
    <el-row :gutter="16">
      <el-col v-for="status in statusesStore.items" :key="status.code" :span="6">
        <el-card shadow="hover">
          <h3>{{ status.name }}</h3>
          <Container
            :group-name="'incidents'"
            :get-child-payload="getChildPayload(status.code)"
            @drop="onDrop(status.code, $event)"
          >
            <Draggable
              v-for="incident in incidentsStore.filteredItems.filter(
                (i) => i.status.code === status.code,
              )"
              :key="incident.id"
            >
              <el-card
                shadow="never"
                style="cursor: pointer; margin-bottom: 8px"
                @click="goToIncident(incident.id)"
              >
                <strong>#{{ incident.id }} {{ incident.title }}</strong>

                <p>
                  {{ incident.description?.slice(0, 200)
                  }}{{ incident.description && incident.description.length > 200 ? '…' : '' }}
                </p>

                <p><b>Создан:</b> {{ formatDateTime(incident.createdAt) }}</p>
                <p><b>Обновлён:</b> {{ formatDateTime(incident.updatedAt) }}</p>

                <el-space justify="space-between">
                  <IncidentPriorityColor
                    :priority-id="incident.priority.id"
                    :label="incident.priority.name"
                  />
                </el-space>

                <template #footer>
                  <el-row justify="space-between" align="middle">
                    <el-space size="small">
                      <el-icon><User /></el-icon>
                      <p v-if="incident.user" style="margin: 0">
                        <UserLink :user="incident.user" :is-link="false" />
                      </p>
                      <p v-else style="margin: 0">Не назначено</p>
                    </el-space>

                    <el-space size="small">
                      <el-badge v-if="incident.commentsCount > 0" :value="incident.commentsCount">
                        <el-icon><ChatDotRound /></el-icon>
                      </el-badge>
                    </el-space>
                  </el-row>
                </template>
              </el-card>
            </Draggable>
          </Container>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { useIncidentsStore } from '@/stores/incidentsStore'
import { useStatusesStore } from '@/stores/statusesStore'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Container, Draggable } from 'vue3-smooth-dnd'
import { User, ChatDotRound } from '@element-plus/icons-vue'
import IncidentPriorityColor from '@/components/incidents/IncidentPriorityColor.vue'
import UserLink from '@/components/users/UserLink.vue'
import { formatDateTime } from '@/utils/date'
import IncidentFilters from '@/components/incidents/IncidentFilters.vue'

const incidentsStore = useIncidentsStore()
const statusesStore = useStatusesStore()
const router = useRouter()

const goToIncident = (id: number) => {
  router.push(`/incident/${id}`)
}

const goToCreate = () => {
  router.push('/incident/create')
}

const getChildPayload = (statusCode: string) => (index: number) => {
  const items = incidentsStore.filteredItems.filter((i) => i.status.code === statusCode)
  return items[index]
}

const onDrop = async (statusCode: string, dropResult: any) => {
  const { addedIndex, payload } = dropResult

  if (addedIndex === null || !payload) return
  const movedItem = payload
  if (movedItem.status?.code === statusCode) return

  const status = statusesStore.items.find((s) => s.code === statusCode)
  if (!status) return

  try {
    await incidentsStore.updateStatus(movedItem.id, status.id)
  } catch (err) {
    ElMessage({
      message: `Не удалось обновить статус инцидента #${movedItem.id}`,
      type: 'error',
      duration: 3000,
    })
  }
}
</script>
