import { defineStore } from 'pinia'
import { usersApi } from '@/api/client'
import { notifyError } from '@/utils/notify'
import { WebUpdateUserRequestRoleEnum, WebUserResponse } from '@/api'

export interface User {
  id: number
  name: string
  username: string
  email: string
  teamId: number | null
  role: string
}

type CreateUserRequest = {
  name: string
  username: string
  email: string
  role: WebUpdateUserRequestRoleEnum
  teamId?: number
  password: string
}

type UpdateUserRequest = {
  name?: string
  username?: string
  email?: string
  role?: WebUpdateUserRequestRoleEnum
  teamId?: number
  password?: string
}

interface UsersState {
  items: User[]
  pollingId: ReturnType<typeof setInterval> | null
}

export const useUsersStore = defineStore('users', {
  state: (): UsersState => ({
    items: [],
    pollingId: null,
  }),

  getters: {
    byTeam: (state) => {
      return (teamId?: number) => {
        if (!teamId || teamId === 0) return state.items
        return state.items.filter((u) => u.teamId === teamId)
      }
    },
  },

  actions: {
    async fetchAll(): Promise<void> {
      try {
        let allItems: User[] = []
        let offset = 0
        const limit = 100

        while (true) {
          const resp = await usersApi.usersGet({ limit, offset })
          const items = resp.items ?? []

          allItems = allItems.concat(items.map(mapApiUser))

          if (items.length < limit) break
          offset += limit
        }

        this.items = allItems
      } catch (e) {
        notifyError('Не удалось обновить список пользователей')
        console.error('users fetch error:', e)
      }
    },

    async create(data: CreateUserRequest) {
      try {
        const res = await usersApi.usersPost({
          request: data,
        })

        const user = mapApiUser(res)
        this.items.push(user)

        return user
      } catch (e) {
        console.error('user create error:', e)
        throw e
      }
    },

    async update(id: number, data: UpdateUserRequest) {
      try {
        await usersApi.usersIdPatch({
          id,
          request: data,
        })

        const item = this.items.find((u) => u.id === id)
        if (item) {
          Object.assign(item, data)
        }
      } catch (e) {
        console.error('user update error:', e)
        throw e
      }
    },

    async remove(id: number) {
      try {
        await usersApi.usersIdDelete({ id })
        this.items = this.items.filter((u) => u.id !== id)
      } catch (e) {
        console.error('user delete error:', e)
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
  },
})

function mapApiUser(item: WebUserResponse): User {
  return {
    id: item.id,
    name: item.name,
    username: item.username,
    email: item.email,
    teamId: item.teamId ?? null,
    role: item.role,
  }
}
