import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

import IncidentList from '@/components/incidents/IncidentList.vue'
import IncidentDetail from '@/components/incidents/IncidentDetail.vue'
import IncidentCreate from '@/components/incidents/IncidentCreate.vue'
import LoginForm from '@/components/auth/LoginForm.vue'
import StatusesView from '@/components/admin/Statuses.vue'

const routes = [
  {
    path: '/login',
    component: LoginForm,
  },
  {
    path: '/',
    component: IncidentList,
    meta: { requiresAuth: true },
  },
  {
    path: '/incident/create',
    component: IncidentCreate,
    props: true,
    meta: { requiresAuth: true },
  },
  {
    path: '/incident/:id',
    component: IncidentDetail,
    props: true,
    meta: { requiresAuth: true },
  },
  {
    path: '/admin/statuses',
    component: StatusesView,
    meta: { requiresAuth: true },
  },
]

export const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()

  if (!authStore.initialized) {
    await authStore.refresh()
  }

  if (to.meta.requiresAuth && !authStore.isAuth) {
    return '/login'
  }

  if (to.path === '/login' && authStore.isAuth) {
    return '/'
  }

  return true
})
