import type { SearchResult } from '~/types'
import { fetchSearchResults } from '~/logic/backendFetch'

export const searchInput = ref<string>('')

export function closeSearch() {
  searchInput.value = ''
}

export async function getSearchResults(): Promise<SearchResult> {
  if (!searchInput.value || searchInput.value.length < 3) {
    return Promise.resolve({} as SearchResult)
  }
  return fetchSearchResults(searchInput.value)
}
