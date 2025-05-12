import { useSessionStorage } from '@vueuse/core'
import { backendFetchRequest } from './fetchFromBackend'

interface AllSlugs {
  timeLastChecked: Date
  slugs: string[]
}

const allSlugs = useSessionStorage<AllSlugs>('all-slugs', {} as AllSlugs)

const lastThirtyRandomSlugs: string[] = []

export async function getAllSlugs() {
  if (allSlugs.value.slugs?.length > 0) {
    const now = new Date()
    const timeLastChecked = new Date(allSlugs.value.timeLastChecked)
    const timeAgo = (now.getTime() - timeLastChecked.getTime()) / 1000
    if (timeAgo < 300) { // 5 minutes
      return allSlugs.value.slugs
    }
  }

  try {
    const response = await backendFetchRequest('slugs')
    const jsonData = await response.json() as string[]
    const newSlugs: AllSlugs = {
      timeLastChecked: new Date(),
      slugs: jsonData,
    }
    allSlugs.value = newSlugs
    return jsonData
  }
  catch (error) {
    console.error('Failed to fetch thumbnails:', error)
    return []
  }
}

export async function getRandomSlug(): Promise<string> {
  const allSlugs: string[] = await getAllSlugs()
  let randomSlug = allSlugs[Math.floor(Math.random() * allSlugs.length)]
  if (allSlugs.length > 30 && lastThirtyRandomSlugs.includes(randomSlug)) {
    randomSlug = await getRandomSlug()
  }
  if (lastThirtyRandomSlugs.length >= 30) {
    lastThirtyRandomSlugs.shift()
  }
  lastThirtyRandomSlugs.push(randomSlug)
  return randomSlug
}

export async function getSlugPosition(slug: string): Promise<number> {
  const allSlugs: string[] = await getAllSlugs()
  const slugPosition = allSlugs.indexOf(slug)
  return slugPosition
}

export async function getPreviousSlug(slug: string): Promise<string> {
  const allSlugs: string[] = await getAllSlugs()
  const slugPosition = await getSlugPosition(slug)
  const previousPosition = (slugPosition > 0 ? slugPosition - 1 : allSlugs.length)
  return allSlugs[previousPosition]
}

export async function getNextSlug(slug: string): Promise<string> {
  const allSlugs: string[] = await getAllSlugs()
  const slugPosition = await getSlugPosition(slug)
  const nextPosition = (slugPosition < allSlugs.length ? slugPosition + 1 : 0)
  return allSlugs[nextPosition]
}
