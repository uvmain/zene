import type { AlbumMetadata, ArtistMetadata, Queue, TrackMetadata, TrackMetadataWithImageUrl } from '~/types'
import { useBackendFetch } from './useBackendFetch'
import { useLogic } from './useLogic'
import { useRandomSeed } from './useRandomSeed'

const { backendFetchRequest, getAlbumTracks, getArtistTracks } = useBackendFetch()
const { randomSeed, refreshRandomSeed, getRandomSeed } = useRandomSeed()
const { trackWithImageUrl, getRandomInteger } = useLogic()

const currentlyPlayingTrack = ref<TrackMetadataWithImageUrl | undefined>()
const currentQueue = ref<Queue | undefined>()

export function usePlaybackQueue() {
  const resetCurrentlyPlayingTrack = () => {
    currentlyPlayingTrack.value = undefined
  }

  const setCurrentlyPlayingTrack = (track: TrackMetadata | TrackMetadataWithImageUrl) => {
    currentlyPlayingTrack.value = trackWithImageUrl(track)
  }

  const setCurrentlyPlayingTrackInQueue = (track: TrackMetadataWithImageUrl) => {
    if (!currentQueue.value) {
      return
    }
    const index = currentQueue.value.tracks.indexOf(track)
    currentQueue.value.position = index
    currentlyPlayingTrack.value = trackWithImageUrl(track)
  }

  const setCurrentQueue = (tracks: TrackMetadata[] | TrackMetadataWithImageUrl[], playFirstTrack: boolean = true) => {
    currentQueue.value = {
      tracks: tracks.map(track => trackWithImageUrl(track)),
      position: 0,
    }
    if (playFirstTrack && tracks.length > 0) {
      setCurrentlyPlayingTrack(tracks[0])
    }
  }

  const clearQueue = () => {
    currentQueue.value = undefined
  }

  const getRandomTrack = async (): Promise<TrackMetadataWithImageUrl> => {
    const seed = getRandomInteger()
    const formData = new FormData()
    formData.append('random', seed.toString())
    formData.append('limit', '1')
    const response = await backendFetchRequest(`tracks`, {
      method: 'POST',
      body: formData,
    })
    const json = await response.json() as TrackMetadata[]
    const randomTrack = trackWithImageUrl(json[0])
    setCurrentlyPlayingTrack(randomTrack)
    currentQueue.value = undefined
    return randomTrack
  }

  const play = async (artist?: ArtistMetadata, album?: AlbumMetadata, track?: TrackMetadata | TrackMetadataWithImageUrl) => {
    if (track) {
      setCurrentlyPlayingTrack(trackWithImageUrl(track))
      currentQueue.value = undefined
    }
    else if (album) {
      const tracks = await getAlbumTracks(album.musicbrainz_album_id)
      setCurrentQueue(tracks)
    }
    else if (artist) {
      const tracks = await getArtistTracks(artist.musicbrainz_artist_id)
      setCurrentQueue(tracks)
    }
  }

  const getRandomTracks = async (): Promise<TrackMetadataWithImageUrl[]> => {
    if (randomSeed.value === 0) {
      randomSeed.value = getRandomInteger()
    }
    const formData = new FormData()
    formData.append('random', randomSeed.value.toString())
    formData.append('limit', '100')
    const response = await backendFetchRequest(`tracks`, {
      method: 'POST',
      body: formData,
    })
    const json = await response.json() as TrackMetadata[]
    const randomTracks = json.map((randomTrack) => {
      return trackWithImageUrl(randomTrack)
    })
    setCurrentQueue(randomTracks)
    return randomTracks
  }

  const getNextTrack = async (): Promise<TrackMetadataWithImageUrl | undefined> => {
    if (currentQueue.value && currentQueue.value.tracks.length) {
      const currentIndex = currentQueue.value.position
      let nextTrack: TrackMetadataWithImageUrl
      if (currentIndex < currentQueue.value.tracks.length - 1) {
        nextTrack = currentQueue.value.tracks[currentIndex + 1]
        currentQueue.value.position = currentIndex + 1
      }
      else {
        nextTrack = currentQueue.value.tracks[0]
        currentQueue.value.position = 0
      }
      setCurrentlyPlayingTrack(nextTrack)
      return nextTrack
    }
    else {
      const randomTrack = await getRandomTrack()
      return randomTrack
    }
  }

  const getPreviousTrack = async (): Promise<TrackMetadataWithImageUrl | undefined> => {
    if (currentQueue.value && currentQueue.value.tracks.length) {
      const currentIndex = currentQueue.value.position
      let prevTrack: TrackMetadataWithImageUrl
      if (currentIndex > 0) {
        prevTrack = currentQueue.value.tracks[currentIndex - 1]
        currentQueue.value.position = currentIndex - 1
      }
      else {
        prevTrack = currentQueue.value.tracks[currentQueue.value.tracks.length - 1]
        currentQueue.value.position = currentQueue.value.tracks.length - 1
      }
      setCurrentlyPlayingTrack(prevTrack)
      return prevTrack
    }
    else {
      const randomTrack = await getRandomTrack()
      return randomTrack
    }
  }

  return {
    currentlyPlayingTrack,
    currentQueue,
    clearQueue,
    setCurrentlyPlayingTrackInQueue,
    resetCurrentlyPlayingTrack,
    setCurrentlyPlayingTrack,
    setCurrentQueue,
    play,
    getNextTrack,
    getPreviousTrack,
    getRandomTrack,
    getRandomTracks,
    randomSeed,
    refreshRandomSeed,
    getRandomSeed,
  }
}
