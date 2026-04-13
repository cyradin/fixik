import { defineStore } from 'pinia'
import { statusesApi } from '@/api/client'
import { notifyError } from '@/utils/notify'
import { extractUserMessage } from '@/utils/errors'

export interface Status {
  id: number
  code: string
  name: string
  sort: number
  description: string
  isFinal: boolean
}

interface StatusesState {
  items: Status[]
  pollingId: ReturnType<typeof setInterval> | null
}

export const useStatusesStore = defineStore('statuses', {
  state: (): StatusesState => ({
    items: [],
    pollingId: null,
  }),

  actions: {
    async fetchAll(): Promise<void> {
      try {
        const resp = await statusesApi.statusesGet({})
        this.items = resp.items.sort((a, b) => a.sort - b.sort)
      } catch (e) {
        notifyError('Не удалось обновить список статусов инцидентов')
        console.error('statuses fetch error:', e)
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
        const res = await statusesApi.statusesPost({
          request: data,
        })

        this.items.push(res)

        return res
      } catch (e) {
        console.error('status create error:', e)
        throw new Error(await extractUserMessage(e))
      }
    },

    async update(id: number, data: any) {
      try {
        await statusesApi.statusesIdPut({
          id,
          request: data,
        })

        const item = this.items.find((i) => i.id === id)
        if (item) Object.assign(item, data)
      } catch (e) {
        console.error('status update error:', e)
        throw new Error(await extractUserMessage(e))
      }
    },

    async remove(id: number) {
      try {
        await statusesApi.statusesIdDelete({ id })

        this.items = this.items.filter((i) => i.id !== id)
      } catch (e) {
        console.error('status delete error:', e)
        throw new Error(await extractUserMessage(e))
      }
    },
  },
})
