<template>
  <el-row justify="center">
    <el-col :span="8">
      <el-card>
        <h2 style="text-align: center; margin-bottom: 20px">Вход</h2>

        <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="onSubmit">
          <el-form-item label="Логин" prop="username">
            <el-input v-model="form.username" size="small" />
          </el-form-item>

          <el-form-item label="Пароль" prop="password">
            <el-input v-model="form.password" type="password" show-password size="small" />
          </el-form-item>

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

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { useAuthStore } from '@/stores/authStore'
import { notifyError } from '@/utils/notify'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref<FormInstance>()

const form = reactive({
  username: '',
  password: '',
})

const rules: FormRules = {
  username: [{ required: true, message: 'Введите логин', trigger: 'blur' }],
  password: [{ required: true, message: 'Введите пароль', trigger: 'blur' }],
}

const onSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      await authStore.login(form.username, form.password)
      router.push('/')
    } catch (e: any) {
      notifyError(e.message)
    }
  })
}
</script>
