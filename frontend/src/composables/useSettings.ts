import { useLocalStorage } from '@vueuse/core'

export const StreamQualities = [96, 128, 160, 192, 256, 'native'] as const

type StreamQuality = typeof StreamQualities[number]

const defaultQuality: StreamQuality = 160

const streamQuality = useLocalStorage<StreamQuality>('streamQuality', defaultQuality)

export function useSettings() {
  return {
    streamQuality,
    StreamQualities,
  }
}
