import { createRouter, createWebHashHistory } from 'vue-router'
import IncidentList from '@/components/incidents/IncidentList.vue'
import IncidentDetail from '@/components/incidents/IncidentDetail.vue'

const routes = [
  {
    path: '/',
    component: IncidentList,
  },
  {
    path: '/incident/:id',
    component: IncidentDetail,
    props: true,
  },
]

export const router = createRouter({
  history: createWebHashHistory(),
  routes,
})
