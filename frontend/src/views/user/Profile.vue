<template>
  <Back />
  <el-card style="max-width: 500px; margin: 40px auto">
    <template #header>
      <span>Профиль</span>
    </template>

    <UserData :user="authStore.user" />

    <el-divider>Смена пароля</el-divider>

    <ChangePasswordForm />
  </el-card>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/authStore'
import ChangePasswordForm from '@/components/users/ChangePasswordForm.vue'
import UserData from '@/components/users/UserData.vue'
import Back from '@/components/layout/Back.vue'

const authStore = useAuthStore()
</script>
