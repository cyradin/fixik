<template>
  <span v-if="user">
    <el-link v-if="isLink" type="primary" @click="goToUser">
      {{ user.name }} (@{{ user.username }})
    </el-link>
    <span v-else>{{ user.name }} (@{{ user.username }})</span>
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
