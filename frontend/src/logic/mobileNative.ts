import { Capacitor } from '@capacitor/core'
import { useLocalStorage } from '@vueuse/core'

export const isMobileNative = Capacitor.isNativePlatform()
export const serverBaseUrl = useLocalStorage('serverBaseUrl', '')
