<template>
  <div v-if="user" style="margin-bottom: 20px">
    <p><b>Имя:</b> {{ user.name }}</p>
    <p><b>Username:</b> {{ user.username }}</p>
    <p><b>Email:</b> {{ user.email }}</p>
    <p><b>Роль:</b> {{ user.role }}</p>
  </div>
</template>

<script setup lang="ts">
import type { CurrentUser } from '@/stores/authStore'

defineProps<{
  user: CurrentUser | null
}>()
</script>
