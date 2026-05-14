<template>
  <div>
    <IncidentFilters :store="incidentsStore" :with-status="false" />

    <el-row justify="space-between" align="middle" style="margin: 16px">
      <h2 style="margin: 0">Инциденты</h2>

      <el-button type="primary" size="large" @click="goToCreate"> + Создать инцидент </el-button>
    </el-row>
    <div class="statuses-scroll">
      <div v-for="status in statusesStore.items" :key="status.code" class="status-column">
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
                  {{ incident.description?.slice(0, 50)
                  }}{{ incident.description && incident.description.length > 50 ? '…' : '' }}
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
                    <el-space size="small" class="incident-user">
                      <el-icon><User /></el-icon>

                      <span v-if="incident.user" class="incident-user-name">
                        <UserLink :user="incident.user" :is-link="false" />
                      </span>

                      <span v-else class="incident-user-empty"> Не назначено </span>
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
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useIncidentsStore } from '@/stores/incidentsStore'
import { useStatusesStore } from '@/stores/statusesStore'
import { useRouter } from 'vue-router'
import { Container, Draggable } from 'vue3-smooth-dnd'
import { User, ChatDotRound } from '@element-plus/icons-vue'
import IncidentPriorityColor from '@/components/incidents/IncidentPriorityColor.vue'
import UserLink from '@/components/users/UserLink.vue'
import { formatDateTime } from '@/utils/date'
import IncidentFilters from '@/components/incidents/IncidentFilters.vue'
import { notifyError } from '@/utils/notify'

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
  } catch (e: any) {
    notifyError(e.message)
  }
}
</script>

<style scoped>
.statuses-scroll {
  display: flex;
  gap: 20px;
  overflow-x: auto;
  padding-bottom: 12px;
}

.status-column {
  flex: 0 0 320px;
}

.status-column :deep(.el-card) {
  border-radius: 12px;
}

.status-column > .el-card {
  background: #f5f7fa;
  border: 1px solid #dcdfe6;
}

.status-column h3 {
  margin: 0 0 16px;
  padding-bottom: 8px;

  font-size: 18px;
  font-weight: 700;
  color: #303133;

  border-bottom: 2px solid #dcdfe6;
}

.status-column .el-card .el-card {
  border: 1px solid #dcdfe6;
  border-left: 4px solid #409eff;

  background: #ffffff;

  transition:
    transform 0.15s ease,
    box-shadow 0.15s ease,
    border-color 0.15s ease;
}

.status-column .el-card .el-card:hover {
  transform: translateY(-2px);

  border-color: #409eff;

  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

.status-column strong {
  display: block;

  margin-bottom: 8px;

  font-size: 15px;
  font-weight: 700;
  line-height: 1.4;

  color: #1f2329;
}

.status-column p {
  margin: 6px 0;

  font-size: 13px;
  line-height: 1.5;

  color: #4e5969;
}

.status-column p b {
  color: #303133;
}

.status-column .el-card__footer {
  background: #fafafa;
}

.incident-user {
  font-size: 14px;
  font-weight: 600;
}

.incident-user-name {
  color: #409eff;
}

.incident-user-empty {
  color: #909399;
  font-style: italic;
}

.incident-user :deep(a) {
  color: #409eff;
  text-decoration: none;
  font-weight: 600;
}

.incident-user :deep(a:hover) {
  text-decoration: underline;
}
</style>
