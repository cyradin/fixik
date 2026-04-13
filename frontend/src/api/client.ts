import { notifyError } from '@/utils/notify'
import { AuthApi } from './apis/AuthApi'
import { IncidentsApi } from './apis/IncidentsApi'
import { PrioritiesApi } from './apis/PrioritiesApi'
import { RolesApi } from './apis/RolesApi'
import { StatusesApi } from './apis/StatusesApi'
import { TeamsApi } from './apis/TeamsApi'
import { UsersApi } from './apis/UsersApi'
import { Configuration } from './runtime'

const basePath = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

const configuration = new Configuration({
  basePath,
  credentials: 'include',
})

export const authApi = new AuthApi(configuration)

function withAuthRefresh<T extends (...args: any[]) => Promise<any>>(fn: T): T {
  return (async (...args: Parameters<T>) => {
    try {
      return await fn(...args)
    } catch (e: any) {
      if (e.response?.status === 401) {
        try {
          await authApi.authRefreshPost()
        } catch {
          throw e
        }

        return fn(...args)
      }
      throw e
    }
  }) as T
}

function wrapApiWithAuthRefresh<T extends object>(api: T): T {
  return new Proxy(api, {
    get(target, prop: string) {
      const value = (target as any)[prop]
      if (typeof value === 'function') {
        return withAuthRefresh(value.bind(target))
      }
      return value
    },
  })
}

export const incidentsApi = wrapApiWithAuthRefresh(new IncidentsApi(configuration))
export const prioritiesApi = wrapApiWithAuthRefresh(new PrioritiesApi(configuration))
export const rolesApi = wrapApiWithAuthRefresh(new RolesApi(configuration))
export const statusesApi = wrapApiWithAuthRefresh(new StatusesApi(configuration))
export const teamsApi = wrapApiWithAuthRefresh(new TeamsApi(configuration))
export const usersApi = wrapApiWithAuthRefresh(new UsersApi(configuration))
