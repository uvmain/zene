import { debugEnabled } from '~/logic/store'

export function toggleDebug(): boolean {
  debugEnabled.value = !debugEnabled.value
  return debugEnabled.value
}

export function debugLog(logMessage: string) {
  if (debugEnabled.value) {
    console.log(`[DEBUG] ${logMessage}`)
  }
}
