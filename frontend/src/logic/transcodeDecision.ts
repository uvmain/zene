import type { DirectPlayProfile, TranscodeClientInfo } from '~/types/subsonicTranscode'
import { directPlayMaxQuality, transcodeStreamQuality } from '~/stores/main'
import { debugLog } from './logger'

async function getDirectPlayProfiles(): Promise<DirectPlayProfile[]> {
  const directPlayBitrate = (directPlayMaxQuality.value === 'Unlimited' ? 10000000 : directPlayMaxQuality.value as number) * 1000
  const audioConfigs: {
    mp3Config: MediaDecodingConfiguration,
    flacConfig: MediaDecodingConfiguration,
    opusConfig: MediaDecodingConfiguration,
    opusWebmConfig: MediaDecodingConfiguration,
    aacConfig: MediaDecodingConfiguration
  } = {
    mp3Config: { type: 'file', audio: { contentType: 'audio/mp3', channels: '2', bitrate: directPlayBitrate, samplerate: 44100 } },
    flacConfig: { type: 'file', audio: { contentType: 'audio/flac', channels: '2', bitrate: directPlayBitrate, samplerate: 44100 } },
    opusConfig: {   type: "file", audio: { contentType: 'audio/ogg; codecs="opus"', channels: '2', bitrate: directPlayBitrate, samplerate: 48000 } },
    opusWebmConfig: { type: "file", audio: { contentType: 'audio/webm; codecs="opus"', channels: '2', bitrate: directPlayBitrate, samplerate: 48000 } },
    aacConfig: { type: "file", audio: { contentType: 'audio/mp4; codecs="mp4a.40.2"', channels: '2', bitrate: directPlayBitrate, samplerate: 44100 } },
  }

  const directPlayProfiles: DirectPlayProfile[] = []

  if (typeof navigator !== 'undefined' && navigator.mediaCapabilities) {

    navigator.mediaCapabilities.decodingInfo(audioConfigs.mp3Config).then((result) => {
      if (result.supported) {
        directPlayProfiles.push({ containers: ['mp3'], audioCodecs: ['mp3'], protocols: ['http'] })
      }
    })

    navigator.mediaCapabilities.decodingInfo(audioConfigs.flacConfig).then((result) => {
      if (result.supported) {
        directPlayProfiles.push({ containers: ['flac'], audioCodecs: ['flac'], protocols: ['http'] })
      }
    })

    navigator.mediaCapabilities.decodingInfo(audioConfigs.opusConfig).then((result) => {
      if (result.supported) {
        directPlayProfiles.push({ containers: ['ogg'], audioCodecs: ['opus'], protocols: ['http'] })
      }
    })

    navigator.mediaCapabilities.decodingInfo(audioConfigs.opusWebmConfig).then((result) => {
      if (result.supported) {
        directPlayProfiles.push({ containers: ['webm'], audioCodecs: ['opus'], protocols: ['http'] })
      }
    })

    navigator.mediaCapabilities.decodingInfo(audioConfigs.aacConfig).then((result) => {
      if (result.supported) {
        directPlayProfiles.push({ containers: ['mp4'], audioCodecs: ['aac'], protocols: ['http'] })
      }
    })
  }
  else {
    debugLog("MediaCapabilities API not supported in this browser. Direct play profiles will be limited.")
  }

  return directPlayProfiles
}

export async function getTranscodeClientInfo(): Promise<TranscodeClientInfo> {
  const transcodeBitrate = (transcodeStreamQuality.value === 'Unlimited' ? 10000000 : transcodeStreamQuality.value as number) * 1000
  const directPlayBitrate = (directPlayMaxQuality.value === 'Unlimited' ? 10000000 : directPlayMaxQuality.value as number) * 1000
  const directPlayProfiles = await getDirectPlayProfiles()
  const transcodeClientInfo: TranscodeClientInfo = {
    "name": "zeneclient",
    "platform": navigator.platform || navigator.userAgent || "unknown",
    "maxAudioBitrate": directPlayBitrate,
    "maxTranscodingAudioBitrate": transcodeBitrate,
    "directPlayProfiles": directPlayProfiles,
    "transcodingProfiles": [
      {
        "container": "mp4",
        "audioCodec": "aac",
        "protocol": "http"
      },
      {
        "container": "mp3",
        "audioCodec": "mp3",
        "protocol": "http"
      }
    ],
    "codecProfiles": []
  }
  return transcodeClientInfo
}