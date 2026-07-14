export interface TranscodeClientInfo {
  name: string
  platform: string
  maxAudioBitrate: number
  maxTranscodingAudioBitrate: number
  directPlayProfiles: DirectPlayProfile[]
  transcodingProfiles: TranscodingProfile[]
  codecProfiles: CodecProfile[]
}

export interface DirectPlayProfile {
  containers: string[]
  audioCodecs: string[]
  protocols: string[]
  maxAudioChannels?: number
}

interface TranscodingProfile {
  container: string
  audioCodec: string
  protocol: string
  maxAudioChannels?: number
}

interface CodecProfile {
  type: string
  name: string
  limitations: Limitation[]
}

interface Limitation {
  name: string
  comparison: string
  values: string[]
  required: boolean
}

export interface SubsonicTranscodeDecisionResponse {
  "subsonic-response": {
    status: string
    version: string
    type: string
    serverVersion: string
    openSubsonic: boolean
    transcodeDecision: TranscodeDecision
  }
}

export interface TranscodeDecision {
  canDirectPlay: boolean
  canTranscode: boolean
  transcodeReason?: string[]
  errorReason?: string
  transcodeParams?: string
  sourceStream?: StreamDetails
  transcodeStream?: StreamDetails
}

interface StreamDetails {
  protocol: string
  container: string
  codec: string
  audioChannels: number
  audioBitrate: number
  audioProfile: string
  audioSamplerate: number
  audioBitdepth: number
}