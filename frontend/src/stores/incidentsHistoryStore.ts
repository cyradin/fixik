import { defineStore } from 'pinia'
import { incidentsApi } from '@/api/client'
import { Incident, mapApiIncidentToIncident } from './incidentsStore'
import { FilterState } from '@/types/filters'

interface IncidentsState {
  filters: FilterState
  items: Incident[]
  historyLoading: boolean
}

export const pageSize = 50

export const useIncidentsHistoryStore = defineStore('incidentsHistory', {
  state: (): IncidentsState => ({
    items: [],
    historyLoading: false,
    filters: {
      priorityIds: [] as number[],
      authorIds: [] as number[],
      userIds: [] as (number | null)[],
      teamIds: [] as (number | null)[],
      statusIds: [] as (number | null)[],
    },
  }),

  actions: {
    async fetch(): Promise<void> {
      this.historyLoading = true

      try {
        const mapNullToZero = (arr: (number | null)[]): number[] => arr.map((id) => id ?? 0)

        const resp = await incidentsApi.incidentsGet({
          limit: 100,
          offset: 0,
          statusIds: mapNullToZero(this.filters.statusIds),
          priorityIds: this.filters.priorityIds,
          userIds: mapNullToZero(this.filters.userIds),
          authorIds: this.filters.authorIds,
          teamIds: mapNullToZero(this.filters.teamIds),
        })

        this.items = (resp.items ?? []).map(mapApiIncidentToIncident)
      } catch (e) {
        console.error('history fetch error:', e)
        throw e
      } finally {
        this.historyLoading = false
      }
    },

    resetFilters() {
      this.filters.priorityIds = []
      this.filters.authorIds = []
      this.filters.userIds = []
      this.filters.teamIds = []
      this.filters.statusIds = []
    },

    togglePriority(priorityId: number) {
      const idx = this.filters.priorityIds.indexOf(priorityId)
      if (idx === -1) {
        this.filters.priorityIds.push(priorityId)
      } else {
        this.filters.priorityIds.splice(idx, 1)
      }
    },
  },
})
