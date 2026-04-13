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
import UserInfo from '@/views/user/Info.vue'
import { PERMISSIONS } from '@/constants/permissions'
import Forbidden from '@/views/Forbidden.vue'

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
    meta: {
      requiresAuth: true,
      permissions: [
        PERMISSIONS.STATUS_CREATE,
        PERMISSIONS.STATUS_UPDATE,
        PERMISSIONS.STATUS_DELETE,
      ],
    },
  },
  {
    path: '/admin/priorities',
    component: AdminPriorities,
    meta: {
      requiresAuth: true,
      permissions: [
        PERMISSIONS.PRIORITY_CREATE,
        PERMISSIONS.PRIORITY_UPDATE,
        PERMISSIONS.PRIORITY_DELETE,
      ],
    },
  },
  {
    path: '/admin/teams',
    component: AdminTeams,
    meta: {
      requiresAuth: true,
      permissions: [PERMISSIONS.TEAM_CREATE, PERMISSIONS.TEAM_UPDATE, PERMISSIONS.TEAM_DELETE],
    },
  },
  {
    path: '/admin/users',
    component: AdminUsers,
    meta: {
      requiresAuth: true,
      permissions: [PERMISSIONS.USER_CREATE, PERMISSIONS.USER_UPDATE, PERMISSIONS.USER_DELETE],
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
