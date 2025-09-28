import type { SearchResult } from '~/types'
import { useDebounce, useSessionStorage } from '@vueuse/core'
import { fetchSearchResults } from './backendFetch'

const searchInput = useSessionStorage<string>('searchInput', '')

export function useSearch() {
  const closeSearch = () => {
    searchInput.value = ''
  }

  const debouncedInput = useDebounce(searchInput, 1000)

  const getSearchResults = async (): Promise<SearchResult> => {
    return fetchSearchResults(debouncedInput.value)
  }

  return {
    getSearchResults,
    searchInput,
    closeSearch,
  }
}
