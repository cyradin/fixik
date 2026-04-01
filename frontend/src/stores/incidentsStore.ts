import { defineStore } from 'pinia'
import { incidentsApi } from '@/api/client'
import { useStatusesStore } from '@/stores/statusesStore'

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
}

export const useIncidentsStore = defineStore('incidents', {
  state: (): IncidentsState => ({
    items: [],
    loading: false,
    pollingId: null,
  }),

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

    async updateStatus(id: number, statusCode: string): Promise<void> {
      try {
        const statusesStore = useStatusesStore()
        const status = statusesStore.items.find((s) => s.code === statusCode)
        if (!status) throw new Error(`Status not found for code ${statusCode}`)

        await incidentsApi.incidentsIdPatch({ id, request: { statusId: status.id } })

        const item = this.items.find((i) => i.id === id)
        if (item) {
          item.status = {
            id: status.id,
            code: status.code,
            name: status.name,
          }
        }
      } catch (e) {
        console.error('update status error:', e)
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
  },
})
