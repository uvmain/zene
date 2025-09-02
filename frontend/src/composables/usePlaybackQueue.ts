import type { Queue } from '~/types'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicArtist } from '~/types/subsonicArtist'
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchAlbum, fetchArtistTopSongs, fetchRandomTracks } from './backendFetch'

const currentlyPlayingTrack = ref<SubsonicSong | undefined>()
const currentQueue = ref<Queue | undefined>()

export function usePlaybackQueue() {
  const resetCurrentlyPlayingTrack = () => {
    currentlyPlayingTrack.value = undefined
  }

  const setCurrentlyPlayingTrack = (track: SubsonicSong) => {
    currentlyPlayingTrack.value = track
  }

  const setCurrentlyPlayingTrackInQueue = (track: SubsonicSong) => {
    if (!currentQueue.value) {
      return
    }
    const index = currentQueue.value.tracks.indexOf(track)
    currentQueue.value.position = index
    currentlyPlayingTrack.value = track
  }

  const setCurrentQueue = (tracks: SubsonicSong[], playFirstTrack: boolean = true) => {
    currentQueue.value = {
      tracks,
      position: 0,
    }
    if (playFirstTrack && tracks.length > 0) {
      setCurrentlyPlayingTrack(tracks[0])
    }
  }

  const clearQueue = () => {
    currentQueue.value = undefined
  }

  const getRandomTrack = async (): Promise<SubsonicSong> => {
    const randomTracks = await fetchRandomTracks(1)
    const randomTrack = randomTracks[0]
    setCurrentlyPlayingTrack(randomTrack)
    currentQueue.value = undefined
    return randomTrack
  }

  const play = async (artist?: SubsonicArtist, album?: SubsonicAlbum, track?: SubsonicSong) => {
    if (track) {
      setCurrentlyPlayingTrack(track)
      currentQueue.value = undefined
    }
    else if (album) {
      const tracks = await fetchAlbum(album.id).then(fetchedAlbum => fetchedAlbum.song)
      setCurrentQueue(tracks)
    }
    else if (artist) {
      const tracks = await fetchArtistTopSongs(artist.id, 100)
      setCurrentQueue(tracks)
    }
  }

  const getRandomTracks = async (): Promise<SubsonicSong[]> => {
    const randomTracks = await fetchRandomTracks()
    setCurrentQueue(randomTracks)
    return randomTracks
  }

  const getNextTrack = async (): Promise<SubsonicSong | undefined> => {
    if (currentQueue.value && currentQueue.value.tracks.length) {
      const currentIndex = currentQueue.value.position
      let nextTrack: SubsonicSong | undefined
      if (currentIndex < currentQueue.value.tracks.length - 1) {
        nextTrack = currentQueue.value.tracks[currentIndex + 1]
        currentQueue.value.position = currentIndex + 1
      }
      else {
        nextTrack = currentQueue.value.tracks[0]
        currentQueue.value.position = 0
      }
      if (nextTrack !== undefined) {
        setCurrentlyPlayingTrack(nextTrack)
      }
      return nextTrack
    }
    else {
      const randomTrack = await getRandomTrack()
      return randomTrack
    }
  }

  const getPreviousTrack = async (): Promise<SubsonicSong | undefined> => {
    if (currentQueue.value && currentQueue.value.tracks.length) {
      const currentIndex = currentQueue.value.position
      let prevTrack: SubsonicSong | undefined
      if (currentIndex > 0) {
        prevTrack = currentQueue.value.tracks[currentIndex - 1]
        currentQueue.value.position = currentIndex - 1
      }
      else {
        prevTrack = currentQueue.value.tracks[currentQueue.value.tracks.length - 1]
        currentQueue.value.position = currentQueue.value.tracks.length - 1
      }
      if (prevTrack !== undefined) {
        setCurrentlyPlayingTrack(prevTrack)
      }
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
  }
}
