import { backendUrl } from '~/stores/main'
import type { Share } from '~/types/share'

export async function fetchShare(shareToken: string): Promise<Share> {
  const url = `${backendUrl.value}/share/${shareToken}`
  const response = await fetch(url)
  const data = await response.json() as Share

  const visitCountUrl = `${backendUrl.value}/share/${shareToken}/incrementvisitcount`
  void fetch(visitCountUrl, { method: 'POST' })
  return data
}
