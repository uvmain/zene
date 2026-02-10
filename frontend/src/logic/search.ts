import type { SearchResult } from '~/types'
import { useSessionStorage } from '@vueuse/core'
import { fetchSearchResults } from '~/logic/backendFetch'

export const searchInput = useSessionStorage<string>('searchInput', '')

export function closeSearch() {
  searchInput.value = ''
}

export async function getSearchResults(): Promise<SearchResult> {
  if (!searchInput.value || searchInput.value.length < 3) {
    return Promise.resolve({} as SearchResult)
  }
  return fetchSearchResults(searchInput.value)
}
