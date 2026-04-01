import { defineStore } from 'pinia'
import { incidentsApi } from '@/api/client'
import { useStatusesStore } from '@/stores/statusesStore'
import { usePrioritiesStore } from '@/stores/prioritiesStore'
import { useTeamsStore } from '@/stores/teamsStore'
import { useUsersStore } from '@/stores/usersStore'

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
}

interface IncidentsState {
  items: Incident[]
  loading: boolean
  pollingId: ReturnType<typeof setInterval> | null
  filters: {
    priorityIds: number[]
    unassignedOnly: boolean
  }
}

export const useIncidentsStore = defineStore('incidents', {
  state: (): IncidentsState => ({
    items: [],
    loading: false,
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
    async fetchAll(): Promise<void> {
      this.loading = true
      try {
        let allItems: Incident[] = []
        let offset = 0
        const limit = 100

        while (true) {
          const resp = await incidentsApi.incidentsGet({ limit, offset })
          const items = resp.items ?? []

          allItems = allItems.concat(
            items.map((item) => ({
              ...item,
              status: {
                id: item.status.id ?? 0,
                code: item.status.code,
                name: item.status.name,
              },
              user: item.user
                ? {
                    id: item.user.id,
                    name: item.user.name,
                    username: item.user.username,
                  }
                : null,
              author: item.author
                ? {
                    id: item.author.id,
                    name: item.author.name,
                    username: item.author.username,
                  }
                : null,
              team: item.team
                ? {
                    id: item.team.id,
                    name: item.team.name,
                    code: item.team.code,
                  }
                : null,
            })),
          )

          if (items.length < limit) break
          offset += limit
        }

        this.items = allItems
      } catch (e) {
        console.error('incidents fetch error:', e)
      } finally {
        this.loading = false
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

    async updateTeam(id: number, teamId: number | undefined) {
      try {
        await incidentsApi.incidentsIdPatch({ id, request: { teamId } })
        const item = this.items.find((i) => i.id === id)
        if (item) {
          if (teamId) {
            const teamsStore = useTeamsStore()
            const team = teamsStore.items.find((t) => t.id === teamId)
            item.team = team ? { ...team } : null
          } else {
            item.team = null
          }
        }
      } catch (e) {
        console.error('update team error:', e)
        throw e
      }
    },

    async updateUser(id: number, userId: number | undefined) {
      try {
        await incidentsApi.incidentsIdPatch({ id, request: { userId } })
        const item = this.items.find((i) => i.id === id)
        if (item) {
          if (userId) {
            const usersStore = useUsersStore()
            const user = usersStore.items.find((u) => u.id === userId)
            item.user = user ? { ...user } : null
          } else {
            item.user = null
          }
        }
      } catch (e) {
        console.error('update user error:', e)
        throw e
      }
    },

    startPolling(interval = 10000): void {
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
