import { IncidentsApi } from './apis/IncidentsApi'
import { PrioritiesApi } from './apis/PrioritiesApi'
import { RolesApi } from './apis/RolesApi'
import { StatusesApi } from './apis/StatusesApi'
import { TeamsApi } from './apis/TeamsApi'
import { UsersApi } from './apis/UsersApi'
import { Configuration } from './runtime'

const basePath = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

const configuration = new Configuration({ basePath })

export const incidentsApi = new IncidentsApi(configuration)
export const prioritiesApi = new PrioritiesApi(configuration)
export const rolesApi = new RolesApi(configuration)
export const statusesApi = new StatusesApi(configuration)
export const teamsApi = new TeamsApi(configuration)
export const usersApi = new UsersApi(configuration)