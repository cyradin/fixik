<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({ username: '', password: '' })

const onSubmit = async () => {
  try {
    await authStore.login(form.username, form.password)
    router.push('/')
  } catch {}
}
</script>

<template>
  <el-row justify="center">
    <el-col :span="8">
      <el-card>
        <h2 style="text-align: center; margin-bottom: 20px">Вход</h2>

        <el-form @submit.prevent="onSubmit">
          <el-form-item label="Логин">
            <el-input v-model="form.username" size="small" />
          </el-form-item>

          <el-form-item label="Пароль">
            <el-input v-model="form.password" type="password" show-password size="small" />
          </el-form-item>

          <el-alert
            v-if="authStore.error"
            :title="authStore.error"
            type="error"
            show-icon
            style="margin-bottom: 12px"
            @close="authStore.clearError()"
          />

          <el-button
            type="primary"
            :loading="authStore.loading"
            @click="onSubmit"
            style="width: 100%"
          >
            Войти
          </el-button>
        </el-form>
      </el-card>
    </el-col>
  </el-row>
</template>
