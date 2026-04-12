<template>
  <el-descriptions v-if="user" :column="1" border size="default" style="margin-bottom: 20px">
    <el-descriptions-item label="Имя">
      {{ user.name }}
    </el-descriptions-item>

    <el-descriptions-item label="Username">
      {{ user.username }}
    </el-descriptions-item>

    <el-descriptions-item label="Email">
      {{ user.email }}
    </el-descriptions-item>

    <el-descriptions-item label="Роль">
      {{ roleStore.nameByCode(user.role) }}
    </el-descriptions-item>
  </el-descriptions>
</template>

<script setup lang="ts">
import type { CurrentUser } from '@/stores/authStore'
import { useRolesStore } from '@/stores/rolesStore'

const roleStore = useRolesStore()

defineProps<{
  user: CurrentUser | null
}>()
</script>
