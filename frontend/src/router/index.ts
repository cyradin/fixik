import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

import IncidentList from '@/views/incidents/List.vue'
import IncidentHistory from '@/views/incidents/History.vue'
import IncidentDetail from '@/views/incidents/Detail.vue'
import IncidentCreate from '@/views/incidents/Create.vue'

import LoginForm from '@/views/auth/Login.vue'
import AdminStatuses from '@/views/admin/Statuses.vue'
import AdminPriorities from '@/views/admin/Priorities.vue'
import AdminTeams from '@/views/admin/Teams.vue'
import AdminUsers from '@/views/admin/Users.vue'

import UserProfile from '@/views/user/Profile.vue'
import UserInfo from '@/views/user/Info.vue'
import { PERMISSION_GROUPS } from '@/constants/permissions'
import Forbidden from '@/views/Forbidden.vue'
import { useRolesStore } from '@/stores/rolesStore'

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
    path: '/history',
    component: IncidentHistory,
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
    meta: {
      requiresAuth: true,
      permissions: PERMISSION_GROUPS.STATUS_ADMIN,
    },
  },
  {
    path: '/admin/priorities',
    component: AdminPriorities,
    meta: {
      requiresAuth: true,
      permissions: PERMISSION_GROUPS.PRIORITY_ADMIN,
    },
  },
  {
    path: '/admin/teams',
    component: AdminTeams,
    meta: {
      requiresAuth: true,
      permissions: PERMISSION_GROUPS.TEAM_ADMIN,
    },
  },
  {
    path: '/admin/users',
    component: AdminUsers,
    meta: {
      requiresAuth: true,
      permissions: PERMISSION_GROUPS.USER_ADMIN,
    },
  },
  {
    path: '/profile',
    component: UserProfile,
    meta: { requiresAuth: true },
  },
  {
    path: '/user/:id',
    component: UserInfo,
  },
  {
    path: '/403',
    component: Forbidden,
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

  const permissions = to.meta.permissions as string[] | undefined

  if (permissions && permissions.length > 0) {
    const allowed = authStore.can(permissions)

    if (!allowed) {
      return '/403'
    }
  }

  return true
})
