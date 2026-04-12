import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

import IncidentList from '@/views/incidents/List.vue'
import IncidentDetail from '@/views/incidents/Detail.vue'
import IncidentCreate from '@/views/incidents/Create.vue'

import LoginForm from '@/views/auth/Login.vue'
import AdminStatuses from '@/views/admin/Statuses.vue'
import AdminPriorities from '@/views/admin/Priorities.vue'
import AdminTeams from '@/views/admin/Teams.vue'
import AdminUsers from '@/views/admin/Users.vue'

import UserProfile from '@/views/user/Profile.vue'

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
    component: AdminStatuses,
    meta: { requiresAuth: true },
  },
  {
    path: '/admin/priorities',
    component: AdminPriorities,
    meta: { requiresAuth: true },
  },
  {
    path: '/admin/teams',
    component: AdminTeams,
    meta: { requiresAuth: true },
  },
  {
    path: '/admin/users',
    component: AdminUsers,
    meta: { requiresAuth: true },
  },
  {
    path: '/profile',
    component: UserProfile,
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
