export const chromecastAvailable = { value: false }
export const chromecastConnected = ref(false)
export const castPlayer = { value: null as cast.framework.RemotePlayer | null }
export const castPlayerController = { value: null as cast.framework.RemotePlayerController | null }
export const castContext = { value: null as cast.framework.CastContext | null }
export const castSession = { value: null as cast.framework.CastSession | null }
export const isCasting = { value: false }
export const savedLocalPosition = { value: 0 }
