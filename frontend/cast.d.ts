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
      currentTime?: number
      autoplay?: boolean
      playbackRate?: number
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
      PlayerState: {
        IDLE: string
        PLAYING: string
        PAUSED: string
        BUFFERING: string
      }
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
      addEventListener(eventType: string, listener: (event: any) => void): void
      removeEventListener(eventType: string, listener: (event: any) => void): void
    }

    interface CastSession {
      loadMedia: (request: chrome.cast.media.LoadRequest) => Promise<any>
      getMediaSession: () => MediaSession | null
      addUpdateListener: (listener: (isAlive: boolean) => void) => void
      removeUpdateListener: (listener: (isAlive: boolean) => void) => void
      addMediaListener: (listener: (event: any) => void) => void
      removeMediaListener: (listener: (event: any) => void) => void
    }

    interface MediaSession {
      media: MediaInfo | null
      addEventListener: (eventType: string, listener: (event: any) => void) => void
      removeEventListener: (eventType: string, listener: (event: any) => void) => void
    }

    interface MediaInfo {
      contentId: string
      contentType: string
      metadata?: any
    }

    enum CastContextEventType {
      CAST_STATE_CHANGED = 'caststatechanged',
      SESSION_STATE_CHANGED = 'sessionstatechanged',
    }

    enum CastState {
      NO_DEVICES_AVAILABLE = 'NO_DEVICES_AVAILABLE',
      NOT_CONNECTED = 'NOT_CONNECTED',
      CONNECTING = 'CONNECTING',
      CONNECTED = 'CONNECTED',
    }

    enum SessionState {
      NO_SESSION = 'NO_SESSION',
      SESSION_STARTING = 'SESSION_STARTING',
      SESSION_STARTED = 'SESSION_STARTED',
      SESSION_START_FAILED = 'SESSION_START_FAILED',
      SESSION_ENDING = 'SESSION_ENDING',
      SESSION_ENDED = 'SESSION_ENDED',
      SESSION_RESUMED = 'SESSION_RESUMED',
    }

    class RemotePlayer {
      canControlVolume: boolean
      canPause: boolean
      canSeek: boolean
      currentTime: number
      displayName: string
      duration: number
      imageUrl: string
      isConnected: boolean
      isMuted: boolean
      isPaused: boolean
      playerState: string
      title: string
      volumeLevel: number
      mediaInfo: MediaInfo | null
      constructor()
    }

    class RemotePlayerController {
      constructor(player: RemotePlayer)
      addEventListener(eventType: string, listener: (event: any) => void): void
      removeEventListener(eventType: string, listener: (event: any) => void): void
      playOrPause(): void
      stop(): void
      seek(): void
      previousTrack(): void
      nextTrack(): void
      setVolumeLevel(): void
      muteOrUnmute(): void
      getSeekPosition(currentTime: number, duration: number): number
      getSeekTime(currentTime: number, duration: number): number
    }

    enum RemotePlayerEventType {
      CURRENT_TIME_CHANGED = 'currentTimeChanged',
      DURATION_CHANGED = 'durationChanged',
      IS_CONNECTED_CHANGED = 'isConnectedChanged',
      IS_MUTED_CHANGED = 'isMutedChanged',
      IS_PAUSED_CHANGED = 'isPausedChanged',
      PLAYER_STATE_CHANGED = 'playerStateChanged',
      VOLUME_LEVEL_CHANGED = 'volumeLevelChanged',
      CAN_CONTROL_VOLUME_CHANGED = 'canControlVolumeChanged',
      CAN_PAUSE_CHANGED = 'canPauseChanged',
      CAN_SEEK_CHANGED = 'canSeekChanged',
      MEDIA_INFO_CHANGED = 'mediaInfoChanged',
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
