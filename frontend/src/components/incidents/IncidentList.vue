<template>
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
            v-for="incident in incidentsStore.items.filter(i => i.status.code === status.code)"
            :key="incident.id"
          >
            <el-card shadow="never">
              <strong>#{{incident.id}} {{ incident.title }}</strong>
              <p>{{ incident.description }}</p>
              <el-tag type="info">{{ incident.priority.name }}</el-tag>
            </el-card>
          </Draggable>
        </Container>
      </el-card>
    </el-col>
  </el-row>
</template>

<script setup>
import { computed } from 'vue'
import { useIncidentsStore } from '@/stores/incidentsStore'
import { useStatusesStore } from '@/stores/statusesStore'
import { ElMessage } from 'element-plus'
import { Container, Draggable } from 'vue3-smooth-dnd'

const incidentsStore = useIncidentsStore()
const statusesStore = useStatusesStore()

const columns = computed(() =>
  statusesStore.items.map(status => ({
    ...status,
    incidents: incidentsStore.items.filter(i => i.status.code === status.code)
  }))
)

const getChildPayload = (statusCode) => (index) => {
  const items = incidentsStore.items.filter(i => i.status.code === statusCode)
  return items[index]
}

const onDrop = async (statusCode, dropResult) => {
  const { removedIndex, addedIndex, payload } = dropResult

  if (addedIndex === null || !payload) return

  const movedItem = payload

  if (movedItem.status?.code === statusCode) return

  try {
    await incidentsStore.updateStatus(movedItem.id, statusCode)
  } catch (err) {
    ElMessage({
      message: `Не удалось обновить статус инцидента #${movedItem.id}`,
      type: 'error',
      duration: 3000,
    })
  }
}
</script>