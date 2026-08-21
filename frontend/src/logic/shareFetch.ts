import { backendUrl } from '~/stores/main'
import type { Share } from '~/types/share'

export async function fetchShare(shareToken: string): Promise<Share> {
  const url = `${backendUrl.value}/share/${shareToken}`
  const response = await fetch(url)
  const data = await response.json() as Share
  return data
}
