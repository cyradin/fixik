<template>
  <el-card shadow="hover" style="width: 100%">
    <template #header>
      <span>Комментарии</span>
    </template>

    <el-row :gutter="16">
      <el-col :span="12">
        <el-input
          v-model="text"
          type="textarea"
          :rows="3"
          placeholder="Написать комментарий..."
          @keydown.ctrl.enter.prevent="submit"
        />

        <div style="margin-top: 10px; display: flex; justify-content: flex-end">
          <el-button type="primary" :loading="loadingCreate" @click="submit"> Отправить </el-button>
        </div>
      </el-col>
    </el-row>

    <el-row>
      <el-col :span="24">
        <el-divider />

        <el-skeleton v-if="loading" :rows="3" animated />

        <el-empty v-else-if="comments.length === 0" description="Комментариев нет" />

        <el-timeline v-else>
          <el-timeline-item v-for="c in comments" :key="c.id" :timestamp="formatDate(c.createdAt)">
            <el-card shadow="never">
              <div>
                <b>{{ c.author.name }}</b>
              </div>
              <div style="margin-top: 4px">
                {{ c.text }}
              </div>
            </el-card>
          </el-timeline-item>
        </el-timeline>

        <el-pagination
          v-if="total > limit"
          style="margin-top: 16px; display: flex; justify-content: center"
          background
          layout="prev, pager, next"
          :page-size="limit"
          :total="total"
          v-model:current-page="page"
          @current-change="onPageChange"
        />
      </el-col>
    </el-row>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useCommentsStore } from '@/stores/commentsStore'
import { formatDateTime } from '@/utils/date'
import { notifyError } from '@/utils/notify'

const props = defineProps<{
  incidentId: number
}>()

const store = useCommentsStore()

const text = ref('')
const limit = 20
const loadingCreate = ref(false)

const comments = computed(() => store.getByIncident(props.incidentId))
const loading = computed(() => store.isLoading(props.incidentId))
const page = computed({
  get: () => store.getPage(props.incidentId),
  set: (val) => store.setPage(props.incidentId, val),
})
const total = computed(() => store.getTotal(props.incidentId))

const load = () => {
  store.fetch(props.incidentId, page.value, limit)
}

const onPageChange = (p: number) => {
  store.setPage(props.incidentId, p)
  store.fetch(props.incidentId, p, limit)
}

const submit = async () => {
  if (!text.value.trim()) return

  loadingCreate.value = true

  try {
    await store.create(props.incidentId, text.value.trim())
    text.value = ''

    store.setPage(props.incidentId, 1)
    store.fetch(props.incidentId, 1, limit)
  } catch (e) {
    notifyError('Не удалось добавить комментарий')
  } finally {
    loadingCreate.value = false
  }
}

const formatDate = (d: string) => formatDateTime(d)

onMounted(load)
</script>
