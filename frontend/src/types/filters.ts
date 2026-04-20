export interface FilterState {
  priorityIds: number[]
  authorIds: number[]
  userIds: (number | null)[]
  teamIds: (number | null)[]
  statusIds: (number | null)[]
}

export interface FiltersStoreInterface {
  filters: FilterState
  togglePriority(id: number): void
  resetFilters(): void
}
