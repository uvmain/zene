import type { AlbumMetadata, ArtistMetadata, Playlist, TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { backendFetchRequest, getAlbumTracks, getArtistTracks } from './fetchFromBackend'
import { trackWithImageUrl } from './logic'

const currentlyPlayingTrack = ref<TrackMetadataWithImageUrl | undefined>()
const currentPlaylist = ref<Playlist | undefined>()

export function usePlaybackQueue() {
  const resetCurrentlyPlayingTrack = () => {
    currentlyPlayingTrack.value = undefined
  }

  const setCurrentlyPlayingTrack = (track: TrackMetadata | TrackMetadataWithImageUrl) => {
    console.log(`setting current track to ${track.filename}`)
    currentlyPlayingTrack.value = trackWithImageUrl(track)
  }

  const setCurrentPlaylist = (tracks: TrackMetadata[] | TrackMetadataWithImageUrl[], playFirstTrack: boolean = true) => {
    currentPlaylist.value = {
      tracks: tracks.map(track => trackWithImageUrl(track)),
      position: 0,
    }
    if (playFirstTrack && tracks.length > 0) {
      setCurrentlyPlayingTrack(tracks[0])
    }
  }

  const getRandomTrack = async (): Promise<TrackMetadataWithImageUrl> => {
    const response = await backendFetchRequest('tracks?random=true&limit=1')
    const json = await response.json() as TrackMetadata[]
    const randomTrack = trackWithImageUrl(json[0])
    setCurrentlyPlayingTrack(randomTrack)
    currentPlaylist.value = undefined
    return randomTrack
  }

  const play = async (artist?: ArtistMetadata, album?: AlbumMetadata, track?: TrackMetadata | TrackMetadataWithImageUrl) => {
    if (track) {
      setCurrentlyPlayingTrack(trackWithImageUrl(track))
      currentPlaylist.value = undefined
    }
    else if (album) {
      const tracks = await getAlbumTracks(album.musicbrainz_album_id)
      setCurrentPlaylist(tracks)
    }
    else if (artist) {
      const tracks = await getArtistTracks(artist.musicbrainz_artist_id)
      setCurrentPlaylist(tracks)
    }
  }

  const getRandomTracks = async (): Promise<TrackMetadataWithImageUrl[]> => {
    const response = await backendFetchRequest('tracks?random=true&limit=100')
    const json = await response.json() as TrackMetadata[]
    const randomTracks = json.map((randomTrack) => {
      return trackWithImageUrl(randomTrack)
    })
    setCurrentPlaylist(randomTracks)
    return randomTracks
  }

  const getNextTrack = async (): Promise<TrackMetadataWithImageUrl | undefined> => {
    if (currentPlaylist.value && currentPlaylist.value.tracks.length) {
      const currentIndex = currentPlaylist.value.position
      let nextTrack: TrackMetadataWithImageUrl
      if (currentIndex < currentPlaylist.value.tracks.length - 1) {
        nextTrack = currentPlaylist.value.tracks[currentIndex + 1]
        currentPlaylist.value.position = currentIndex + 1
      }
      else {
        nextTrack = currentPlaylist.value.tracks[0]
        currentPlaylist.value.position = 0
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
    if (currentPlaylist.value && currentPlaylist.value.tracks.length) {
      const currentIndex = currentPlaylist.value.position
      let prevTrack: TrackMetadataWithImageUrl
      if (currentIndex > 0) {
        prevTrack = currentPlaylist.value.tracks[currentIndex - 1]
        currentPlaylist.value.position = currentIndex - 1
      }
      else {
        prevTrack = currentPlaylist.value.tracks[currentPlaylist.value.tracks.length - 1]
        currentPlaylist.value.position = currentPlaylist.value.tracks.length - 1
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
    currentPlaylist,
    resetCurrentlyPlayingTrack,
    setCurrentlyPlayingTrack,
    setCurrentPlaylist,
    play,
    getNextTrack,
    getPreviousTrack,
    getRandomTrack,
    getRandomTracks,
  }
}
