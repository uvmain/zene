import dayjs from 'dayjs'

export function niceDate(dateString: string): string {
  const date = dayjs(dateString)
  return date.isValid() ? date.format('DD/MM/YYYY') : 'Invalid Date'
}

export function getThumbnailPath(slug: string) {
  return `/api/thumbnail/${slug}`
}
