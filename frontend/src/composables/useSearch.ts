import type { SearchResult } from '~/types'
import { useSessionStorage } from '@vueuse/core'
import { fetchSearchResults } from './backendFetch'

const searchInput = useSessionStorage<string>('searchInput', '')

export function useSearch() {
  const closeSearch = () => {
    searchInput.value = ''
  }

  const getSearchResults = async (): Promise<SearchResult> => {
    if (!searchInput.value || searchInput.value.length < 3) {
      return Promise.resolve({} as SearchResult)
    }
    return fetchSearchResults(searchInput.value)
  }

  return {
    getSearchResults,
    searchInput,
    closeSearch,
  }
}
