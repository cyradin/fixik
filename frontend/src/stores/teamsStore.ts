import { defineStore } from 'pinia'
import { teamsApi } from '@/api/client'
import { notifyError } from '@/utils/notify'
import { extractUserMessage } from '@/utils/errors'

export interface Team {
  id: number
  name: string
  code: string
  description: string
  sort: number
}

interface TeamsState {
  items: Team[]
  pollingId: ReturnType<typeof setInterval> | null
}

export const useTeamsStore = defineStore('teams', {
  state: (): TeamsState => ({
    items: [],
    pollingId: null,
  }),

  actions: {
    async fetchAll(): Promise<void> {
      try {
        const resp = await teamsApi.teamsGet({})
        this.items = resp.items.sort((a, b) => a.sort - b.sort)
      } catch (e) {
        notifyError('Не удалось обновить список команд')
        console.error('teams fetch error:', e)
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
        const res = await teamsApi.teamsPost({
          request: data,
        })

        this.items.push(res)
        return res
      } catch (e) {
        console.error('team create error:', e)
        throw new Error(await extractUserMessage(e))
      }
    },

    async update(id: number, data: any) {
      try {
        await teamsApi.teamsIdPut({
          id,
          request: data,
        })

        const item = this.items.find((i) => i.id === id)
        if (item) Object.assign(item, data)
      } catch (e) {
        console.error('team update error:', e)
        throw new Error(await extractUserMessage(e))
      }
    },

    async remove(id: number) {
      try {
        await teamsApi.teamsIdDelete({ id })

        this.items = this.items.filter((i) => i.id !== id)
      } catch (e) {
        console.error('team delete error:', e)
        throw new Error(await extractUserMessage(e))
      }
    },
  },
})
