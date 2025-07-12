export {}

declare global {
  // chrome.cast types
  namespace chrome.cast {
    const VERSION: string
    const isAvailable: boolean

    enum AutoJoinPolicy {
      TAB_AND_ORIGIN_SCOPED,
      ORIGIN_SCOPED,
      PAGE_SCOPED,
    }

    enum SessionStatus {
      STOPPED,
      CONNECTED,
      DISCONNECTED,
    }

    enum ReceiverAvailability {
      AVAILABLE,
      UNAVAILABLE,
    }

    enum Capability {
      VIDEO_OUT,
      AUDIO_OUT,
      VIDEO_IN,
      AUDIO_IN,
    }

    class SessionRequest {
      constructor(appId: string)
    }

    class ApiConfig {
      constructor(
        sessionRequest: SessionRequest,
        sessionListener: (session: Session) => void,
        receiverListener: (availability: ReceiverAvailability) => void
      )
    }

    class MediaInfo {
      contentId: string
      contentType: string
      metadata?: any
      constructor(contentId: string, contentType: string)
    }

    class LoadRequest {
      media: MediaInfo
      constructor(mediaInfo: MediaInfo)
    }

    class Session {
      loadMedia(
        request: LoadRequest,
        successCallback: (media: any) => void,
        errorCallback: (error: any) => void
      ): void
    }

    const media: {
      DEFAULT_MEDIA_RECEIVER_APP_ID: string
      MediaInfo: typeof MediaInfo
      LoadRequest: typeof LoadRequest
    }
  }

  // cast.framework types
  namespace cast.framework {
    function getCastContext(): CastContext
    class CastContext {
      static getInstance(): CastContext
      setOptions(options: {
        receiverApplicationId: string
        autoJoinPolicy: chrome.cast.AutoJoinPolicy
      }): void
      getCurrentSession(): CastSession | null
    }

    interface CastSession {
      loadMedia: (request: chrome.cast.media.LoadRequest) => Promise<any>
    }
  }

  interface Window {
    __onGCastApiAvailable?: (isAvailable: boolean) => void
    cast?: {
      isAvailable?: boolean
      framework?: typeof cast.framework
    }
    chrome?: typeof chrome
  }
}
