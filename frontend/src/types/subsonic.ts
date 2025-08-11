export interface SubsonicResponse {
  'subsonic-response': {
    status: string
    version: string
    type: string
    serverVersion: string
    openSubsonic: boolean
    error?: {
      code: number
      message: string
      helpUrl?: string
    }
    [key: string]: any
  }
}
