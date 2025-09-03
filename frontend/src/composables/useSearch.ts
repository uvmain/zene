import type { SearchResult } from '~/types'
import { useSessionStorage } from '@vueuse/core'
import { fetchSearchResults } from './backendFetch'

const searchInput = useSessionStorage<string>('searchInput', '')

export function useSearch() {
  const closeSearch = () => {
    searchInput.value = ''
  }

  const getSearchResults = async (): Promise<SearchResult> => {
    const searchResult = await fetchSearchResults(searchInput.value)
    return searchResult
  }

  return {
    getSearchResults,
    searchInput,
    closeSearch,
  }
}
