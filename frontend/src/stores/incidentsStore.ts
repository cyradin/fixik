import { defineStore } from 'pinia'
import { incidentsApi } from '@/api/client'
import { useStatusesStore } from '@/stores/statusesStore'
import { usePrioritiesStore } from '@/stores/prioritiesStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { useUsersStore } from '@/stores/usersStore'
import { notifyError } from '@/utils/notify'

export interface IncidentStatus {
  id: number
  code: string
  name: string
}

export interface Priority {
  id: number
  name: string
}

export interface IncidentUser {
  id: number
  name: string
  username: string
}

export interface IncidentTeam {
  id: number
  name: string
  code: string
}

export interface Incident {
  id: number
  title: string
  description: string
  status: IncidentStatus
  priority: Priority
  user: IncidentUser | null
  author: IncidentUser | null
  team: IncidentTeam | null
  createdAt: string
  updatedAt: string
  commentsCount: number
}

interface CreateIncident {
  title: string
  description: string
  statusId: number
  priorityId: number
  teamId?: number
  userId?: number
}

interface IncidentsState {
  items: Incident[]
  pollingId: ReturnType<typeof setInterval> | null
  filters: {
    priorityIds: number[]
    unassignedOnly: boolean
  }
}

export const useIncidentsStore = defineStore('incidents', {
  state: (): IncidentsState => ({
    items: [],
    pollingId: null,
    filters: {
      priorityIds: [],
      unassignedOnly: false,
    },
  }),

  getters: {
    filteredItems: (state) => {
      return state.items.filter((incident) => {
        if (
          state.filters.priorityIds.length > 0 &&
          !state.filters.priorityIds.includes(incident.priority.id)
        ) {
          return false
        }

        if (state.filters.unassignedOnly && incident.user) {
          return false
        }

        return true
      })
    },
  },

  actions: {
    set(items: Incident[]) {
      this.items = items
    },

    addLocal(incident: Incident) {
      const exists = this.items.find((i) => i.id === incident.id)
      if (!exists) {
        this.items.push(incident)
      }
    },

    removeLocal(id: number) {
      this.items = this.items.filter((i) => i.id !== id)
    },

    async fetchAll(): Promise<void> {
      try {
        let allItems: Incident[] = []
        let offset = 0
        const limit = 100

        while (true) {
          const resp = await incidentsApi.incidentsGet({ limit, offset })
          const items = resp.items ?? []

          allItems = allItems.concat(items.map(mapApiIncidentToIncident))

          if (items.length < limit) break
          offset += limit
        }

        this.items = allItems
      } catch (e) {
        notifyError('Не удалось обновить список инцидентов')
        console.error('incidents fetch error:', e)
      }
    },

    async updateStatus(id: number, statusId: number) {
      try {
        await incidentsApi.incidentsIdPatch({ id, request: { statusId } })
        const item = this.items.find((i) => i.id === id)
        if (item) {
          const statusesStore = useStatusesStore()
          const status = statusesStore.items.find((s) => s.id === statusId)
          if (status) {
            item.status = { ...status }
          }
        }
      } catch (e) {
        console.error('update status error:', e)
        throw e
      }
    },

    async updatePriority(id: number, priorityId: number) {
      try {
        await incidentsApi.incidentsIdPatch({ id, request: { priorityId } })
        const item = this.items.find((i) => i.id === id)
        if (item) {
          const prioritiesStore = usePrioritiesStore()
          const priority = prioritiesStore.items.find((p) => p.id === priorityId)
          if (priority) {
            item.priority = { ...priority }
          }
        }
      } catch (e) {
        console.error('update priority error:', e)
        throw e
      }
    },

    async updateDescription(id: number, description: string) {
      try {
        await incidentsApi.incidentsIdPatch({ id, request: { description } })
        const item = this.items.find((i) => i.id === id)
        if (item) item.description = description
      } catch (e) {
        console.error('update description error:', e)
        throw e
      }
    },

    async updateTeam(id: number, teamId?: number) {
      try {
        await incidentsApi.incidentsIdPatch({
          id,
          request: { teamId: teamId || 0 },
        })

        const item = this.items.find((i) => i.id === id)
        if (!item) return

        const teamsStore = useTeamsStore()
        item.team =
          teamId && teamId > 0 ? (teamsStore.items.find((t) => t.id === teamId) ?? null) : null
      } catch (e) {
        console.error('update team error:', e)
        throw e
      }
    },

    async updateUser(id: number, userId?: number) {
      try {
        await incidentsApi.incidentsIdPatch({
          id,
          request: { userId: userId || 0 },
        })

        const item = this.items.find((i) => i.id === id)
        if (!item) return

        const usersStore = useUsersStore()
        item.user =
          userId && userId > 0 ? (usersStore.items.find((u) => u.id === userId) ?? null) : null
      } catch (e) {
        console.error('update user error:', e)
        throw e
      }
    },

    async create(dto: CreateIncident) {
      try {
        const res = await incidentsApi.incidentsPost({ request: dto })
        const incident = mapApiIncidentToIncident(res)

        this.addLocal(incident)

        return incident
      } catch (e) {
        notifyError('Не удалось создать инцидент')
        console.error(e)
        throw e
      }
    },

    async delete(id: number) {
      try {
        await incidentsApi.incidentsIdDelete({ id })
        this.removeLocal(id)
      } catch (e) {
        console.error('Ошибка удаления инцидента', e)
        throw e
      }
    },

    startPolling(interval = 5000): void {
      this.stopPolling()
      this.fetchAll()
      this.pollingId = setInterval(() => this.fetchAll(), interval)
    },

    stopPolling(): void {
      if (this.pollingId) {
        clearInterval(this.pollingId)
        this.pollingId = null
      }
    },

    togglePriority(priorityId: number) {
      const idx = this.filters.priorityIds.indexOf(priorityId)
      if (idx === -1) {
        this.filters.priorityIds.push(priorityId)
      } else {
        this.filters.priorityIds.splice(idx, 1)
      }
    },

    toggleUnassigned() {
      this.filters.unassignedOnly = !this.filters.unassignedOnly
    },

    resetFilters() {
      this.filters.priorityIds = []
      this.filters.unassignedOnly = false
    },
  },
})

function mapApiIncidentToIncident(item: any): Incident {
  const statusesStore = useStatusesStore()
  const prioritiesStore = usePrioritiesStore()
  const teamsStore = useTeamsStore()
  const usersStore = useUsersStore()

  return {
    ...item,
    status: {
      id: item.status.id ?? 0,
      code: item.status.code,
      name: item.status.name,
    },
    priority: item.priority
      ? { id: item.priority.id, name: item.priority.name }
      : { id: 0, name: 'Неизвестный' },
    user: item.user
      ? { id: item.user.id, name: item.user.name, username: item.user.username }
      : null,
    author: item.author
      ? { id: item.author.id, name: item.author.name, username: item.author.username }
      : null,
    team: item.team ? { id: item.team.id, name: item.team.name, code: item.team.code } : null,
    createdAt: item.createdAt,
    updatedAt: item.updatedAt,
    commentsCount: item.commentsCount ?? 0,
  }
}
