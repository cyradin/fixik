<template>
  <span v-if="user">
    <el-link type="primary" @click="goToUser"> {{ user.name }} (@{{ user.username }}) </el-link>
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
}>()

const router = useRouter()

const goToUser = () => {
  if (!props.user) return
  router.push(`/user/${props.user.id}`)
}
</script>
