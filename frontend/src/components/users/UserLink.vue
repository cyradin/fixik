<template>
  <span v-if="user">
    <el-link v-if="isLink" type="primary" @click="goToUser">
      <el-tooltip v-if="user.username" placement="top">
        <template #content>
          <span>@{{ user.username }}</span>
        </template>
        {{ user.name }}
      </el-tooltip>
    </el-link>
    <span v-else>
      <el-tooltip v-if="user.username" placement="top">
        <template #content>
          <span>@{{ user.username }}</span>
        </template>
        {{ user.name }}
      </el-tooltip>
    </span>
  </span>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ElIcon } from 'element-plus'
import { ChatLineRound } from '@element-plus/icons-vue'

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
