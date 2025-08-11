export interface SubsonicLyricsResponse {
  'subsonic-response': {
    status: string
    version: string
    type: string
    serverVersion: string
    openSubsonic: boolean
    lyricsList: LyricsList
  }
}

export interface LyricsList {
  structuredLyrics: StructuredLyric[]
}

export interface StructuredLyric {
  displayArtist: string
  displayTitle: string
  lang: string
  offset: number
  synced: boolean
  line: StructuredLyricLine[]
}

export interface StructuredLyricLine {
  start?: number
  value: string
}
