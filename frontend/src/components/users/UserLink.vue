<template>
  <span v-if="user">
    <span v-if="isLink" class="user-link" @click="goToUser">
      <el-tooltip v-if="user.username" placement="top">
        <template #content>
          <span>@{{ user.username }}</span>
        </template>
        {{ user.name }}
      </el-tooltip>

      <template v-else>
        {{ user.name }}
      </template>
    </span>

    <span v-else>
      <el-tooltip v-if="user.username" placement="top">
        <template #content>
          <span>@{{ user.username }}</span>
        </template>
        {{ user.name }}
      </el-tooltip>

      <template v-else>
        {{ user.name }}
      </template>
    </span>
  </span>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

interface User {
  id: number
  name: string
  username: string
}

const props = defineProps<{
  user: User | null
  isLink: boolean
}>()

const router = useRouter()

const goToUser = () => {
  if (!props.user) return
  router.push(`/user/${props.user.id}`)
}
</script>

<style scoped>
.user-link {
  color: var(--el-color-primary);
  cursor: pointer;
  text-decoration: underline;
  text-decoration-color: transparent;
  transition:
    text-decoration-color 0.2s,
    opacity 0.2s;
}

.user-link:hover {
  text-decoration-color: currentColor;
  opacity: 0.85;
}
</style>
