import { useLocalStorage } from '@vueuse/core'

const useDebugBool = useLocalStorage('useDebugBool', false)

export function useDebug() {
  const toggleDebug = (): boolean => {
    useDebugBool.value = !useDebugBool.value
    return useDebugBool.value
  }

  const debugLog = (logMessage: string) => {
    if (useDebugBool.value) {
      console.log(`[DEBUG] ${logMessage}`)
    }
  }

  return {
    toggleDebug,
    debugLog,
    useDebugBool,
  }
}
