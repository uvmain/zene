import { useLocalStorage } from '@vueuse/core'

export const StreamQualities = [128, 160, 192, 320, 'native'] as const
export type StreamQuality = typeof StreamQualities[number]

const defaultQuality: StreamQuality = 160

const streamQuality = useLocalStorage<StreamQuality>('streamQuality', defaultQuality)

export function useSettings() {
  return {
    streamQuality,
    StreamQualities,
  }
}
