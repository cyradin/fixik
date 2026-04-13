import { defineStore } from 'pinia'
import { prioritiesApi } from '@/api/client'
import { notifyError } from '@/utils/notify'
import { extractUserMessage } from '@/utils/errors'

export interface Priority {
  id: number
  name: string
  code: string
  description: string
  sort: number
}

interface PrioritiesState {
  items: Priority[]
  pollingId: ReturnType<typeof setInterval> | null
}

export const usePrioritiesStore = defineStore('priorities', {
  state: (): PrioritiesState => ({
    items: [],
    pollingId: null,
  }),

  actions: {
    async fetchAll(): Promise<void> {
      try {
        const resp = await prioritiesApi.prioritiesGet({})
        this.items = resp.items.sort((a, b) => a.sort - b.sort)
      } catch (e) {
        notifyError('Не удалось обновить список приоритетов')
        console.error('priorities fetch error:', e)
      }
    },

    startPolling(interval = 30000): void {
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

    async create(data: any) {
      try {
        const res = await prioritiesApi.prioritiesPost({
          request: data,
        })

        this.items.push(res)
        return res
      } catch (e) {
        console.error('priority create error:', e)
        throw new Error(await extractUserMessage(e))
      }
    },

    async update(id: number, data: any) {
      try {
        await prioritiesApi.prioritiesIdPut({
          id,
          request: data,
        })

        const item = this.items.find((i) => i.id === id)
        if (item) Object.assign(item, data)
      } catch (e) {
        console.error('priority update error:', e)
        throw new Error(await extractUserMessage(e))
      }
    },

    async remove(id: number) {
      try {
        await prioritiesApi.prioritiesIdDelete({ id })
        this.items = this.items.filter((i) => i.id !== id)
      } catch (e) {
        console.error('priority delete error:', e)
        throw new Error(await extractUserMessage(e))
      }
    },
  },
})
