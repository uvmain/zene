import { useLocalStorage } from '@vueuse/core'

export const debugEnabled = useLocalStorage('debugEnabled', false)

export function toggleDebug(): boolean {
  debugEnabled.value = !debugEnabled.value
  return debugEnabled.value
}

export function debugLog(logMessage: string) {
  if (debugEnabled.value) {
    console.log(`[DEBUG] ${logMessage}`)
  }
}
