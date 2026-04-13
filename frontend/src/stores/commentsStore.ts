import { defineStore } from 'pinia'
import { incidentsApi } from '@/api/client'
import type { WebIncidentComment } from '@/api/models'
import { notifyError } from '@/utils/notify'

interface CommentsState {
  itemsByIncident: Record<number, WebIncidentComment[]>
  totalByIncident: Record<number, number>
  pageByIncident: Record<number, number>
  loadingByIncident: Record<number, boolean>
}

export const useCommentsStore = defineStore('comments', {
  state: (): CommentsState => ({
    itemsByIncident: {},
    totalByIncident: {},
    pageByIncident: {},
    loadingByIncident: {},
  }),

  getters: {
    getByIncident: (state) => (incidentId: number) => state.itemsByIncident[incidentId] ?? [],

    getTotal: (state) => (incidentId: number) => state.totalByIncident[incidentId] ?? 0,

    getPage: (state) => (incidentId: number) => state.pageByIncident[incidentId] ?? 1,

    isLoading: (state) => (incidentId: number) => !!state.loadingByIncident[incidentId],
  },

  actions: {
    async fetch(incidentId: number, page = 1, limit = 20) {
      this.loadingByIncident[incidentId] = true

      try {
        const res = await incidentsApi.incidentsIdCommentsGet({
          id: incidentId,
          limit,
          offset: (page - 1) * limit,
        })

        this.itemsByIncident[incidentId] = res.items ?? []
        this.totalByIncident[incidentId] = res.pagination?.total ?? 0
        this.pageByIncident[incidentId] = page
      } catch (e) {
        notifyError('Не удалось загрузить комментарии')
        console.error(e)
      } finally {
        this.loadingByIncident[incidentId] = false
      }
    },

    async setPage(incidentId: number, page: number, limit = 20) {
      await this.fetch(incidentId, page, limit)
    },

    async create(incidentId: number, text: string) {
      try {
        const res = await incidentsApi.incidentsIdCommentsPost({
          id: incidentId,
          request: { text },
        })

        const comment: WebIncidentComment = res

        if (!this.itemsByIncident[incidentId]) {
          this.itemsByIncident[incidentId] = []
        }

        this.itemsByIncident[incidentId].unshift(comment)
        this.totalByIncident[incidentId] = (this.totalByIncident[incidentId] || 0) + 1

        return comment
      } catch (e) {
        throw e
      }
    },

    clear(incidentId: number) {
      delete this.itemsByIncident[incidentId]
      delete this.totalByIncident[incidentId]
      delete this.pageByIncident[incidentId]
      delete this.loadingByIncident[incidentId]
    },
  },
})
