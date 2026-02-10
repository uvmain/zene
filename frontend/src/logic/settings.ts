import { useLocalStorage } from '@vueuse/core'

export const StreamQualities = [96, 128, 160, 192, 256, 'native'] as const

export type StreamQuality = typeof StreamQualities[number]

const defaultQuality: StreamQuality = 160

export const streamQuality = useLocalStorage<StreamQuality>('streamQuality', defaultQuality)
