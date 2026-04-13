<template>
  <Back />
  <el-card style="max-width: 500px; margin: 40px auto">
    <template #header>
      <span>Профиль пользователя</span>
    </template>

    <UserData v-if="user" :user="user" />

    <el-empty v-else description="Пользователь не найден" />
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUsersStore } from '@/stores/usersStore'
import UserData from '@/components/users/UserData.vue'
import Back from '@/components/layout/Back.vue'

const route = useRoute()
const usersStore = useUsersStore()

const userId = computed(() => Number(route.params.id))

const user = computed(() => usersStore.items.find((u) => u.id === userId.value) || null)
</script>
