<template>
  <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
    <el-form-item label="Текущий пароль" prop="currentPassword">
      <el-input v-model="form.currentPassword" type="password" show-password />
    </el-form-item>

    <el-form-item label="Новый пароль" prop="newPassword">
      <el-input v-model="form.newPassword" type="password" show-password />
    </el-form-item>

    <el-button type="primary" :loading="loading" @click="submit"> Сменить пароль </el-button>
  </el-form>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { notifyError, notifySuccess } from '@/utils/notify'

const authStore = useAuthStore()

const formRef = ref()
const loading = ref(false)

const form = reactive({
  currentPassword: '',
  newPassword: '',
})

const rules = {
  currentPassword: [{ required: true, message: 'Введите текущий пароль', trigger: 'blur' }],
  newPassword: [
    { required: true, message: 'Введите новый пароль', trigger: 'blur' },
    { min: 6, message: 'Минимум 6 символов', trigger: 'blur' },
  ],
}

const submit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true

    await authStore.changePassword(form.currentPassword, form.newPassword)

    notifySuccess('Пароль изменён')

    form.currentPassword = ''
    form.newPassword = ''
  } catch (e: any) {
    notifyError(e.message || 'Не удалось изменить пароль')
    console.error(e)
  } finally {
    loading.value = false
  }
}
</script>
