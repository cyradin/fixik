import { defineStore } from 'pinia'
import { teamsApi } from '@/api/client'

export interface Team {
    id: number
    code: string
    name: string
}

interface TeamsState {
    items: Team[]
}

export const useTeamsStore = defineStore('teams', {
    state: (): TeamsState => ({
        items: [],
    }),

    actions: {
        async fetchAll(): Promise<void> {
            try {
                const resp = await teamsApi.teamsGet({})
                this.items = resp.items
            } catch (e) {
                console.error('teams fetch error:', e)
            }
        },
    },
})