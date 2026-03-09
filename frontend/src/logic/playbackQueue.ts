import type { PlayItem, Queue } from '~/types'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicIndexArtist } from '~/types/subsonicArtist'
import type { SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import type { SubsonicSong } from '~/types/subsonicSong'
import { computedAsync } from '@vueuse/core'
import { audioElement, playWhenReady } from '~/logic/audioElement'
import { fetchAlbum, fetchArtistTopSongs, fetchRandomTracks } from '~/logic/backendFetch'
import { getAuthenticatedTrackUrl } from '~/logic/common'
import { postPlaycount } from '~/logic/playerUtils'
import { routeTracks } from '~/logic/routeTracks'
import { repeatStatus, shuffleEnabled } from '~/logic/store'
import { episodeIsStored, getStoredEpisode } from '~/stores/usePodcastStore'

export const currentlyPlayingItem = ref<PlayItem>({})
export const currentQueue = ref<Queue | undefined>()
export const isPlaying = ref(false)
export const playcountPosted = ref(false)
export const currentTime = ref(0)

let previousIndexes: number[] = []

export const trackUrl = computedAsync(async () => {
  if (currentlyPlayingItem.value.track !== undefined) {
    return getAuthenticatedTrackUrl(currentlyPlayingItem.value.track.musicBrainzId)
  }
  else if (currentlyPlayingItem.value.podcastEpisode) {
    return episodeIsStored(currentlyPlayingItem.value.podcastEpisode.streamId).then(async (stored) => {
      if (!stored && currentlyPlayingItem.value.podcastEpisode) {
        return getAuthenticatedTrackUrl(currentlyPlayingItem.value.podcastEpisode.streamId, true)
      }
      else if (currentlyPlayingItem.value.podcastEpisode) {
        return getStoredEpisode(currentlyPlayingItem.value.podcastEpisode.streamId).then((blob) => {
          return URL.createObjectURL(blob)
        })
      }
    })
  }
  return ''
})

export function handlePlay(track: SubsonicSong) {
  if (currentQueue.value?.tracks.some(queueTrack => queueTrack.id === track.id)) {
    setCurrentlyPlayingTrack(track)
  }
  else if (routeTracks.value?.some(queueTrack => queueTrack.musicBrainzId === track.musicBrainzId)) {
    setCurrentQueue(routeTracks.value)
    setCurrentlyPlayingTrack(track)
  }
  else {
    void play({ track })
  }
}

export function setCurrentlyPlayingTrack(track: SubsonicSong) {
  if (currentQueue.value) {
    const index = currentQueue.value.tracks.indexOf(track)
    currentQueue.value.position = index
  }
  currentlyPlayingItem.value = { track }
  playcountPosted.value = false
  playWhenReady({ track })
}

function setCurrentQueue(tracks: SubsonicSong[]) {
  clearQueue()
  let index = 0
  if (shuffleEnabled.value && tracks.length > 0) {
    index = Math.floor(Math.random() * tracks.length)
  }
  currentQueue.value = {
    tracks,
    position: index,
  }
  setCurrentlyPlayingTrack(tracks[index])
}

export function clearQueue() {
  previousIndexes = []
  currentQueue.value = undefined
}

async function getRandomTrack(): Promise<SubsonicSong> {
  const randomTracks = await fetchRandomTracks(1)
  const randomTrack = randomTracks[0]
  setCurrentlyPlayingTrack(randomTrack)
  clearQueue()
  return randomTrack
}

interface PlayOptions {
  artist?: SubsonicIndexArtist
  album?: SubsonicAlbum
  track?: SubsonicSong
  podcastEpisode?: SubsonicPodcastEpisode
}

export async function play(playOptions: PlayOptions) {
  if (playOptions.track) {
    setCurrentlyPlayingTrack(playOptions.track)
    clearQueue()
  }
  else if (playOptions.album) {
    const tracks = await fetchAlbum(playOptions.album.id).then(fetchedAlbum => fetchedAlbum.song)
    setCurrentQueue(tracks)
  }
  else if (playOptions.artist) {
    const tracks = await fetchArtistTopSongs(playOptions.artist.id, 100)
    setCurrentQueue(tracks)
  }
  else if (playOptions.podcastEpisode) {
    currentlyPlayingItem.value = { podcastEpisode: playOptions.podcastEpisode }
    clearQueue()
  }
}

export async function getRandomTracks(size: number = 10): Promise<SubsonicSong[]> {
  const randomTracks = await fetchRandomTracks(size)
  setCurrentQueue(randomTracks)
  return randomTracks
}

export async function handleNextTrack(): Promise<SubsonicSong | undefined> {
  if (currentQueue.value && currentQueue.value.tracks.length) {
    const currentIndex = currentQueue.value.position
    let nextTrack: SubsonicSong | undefined
    if (shuffleEnabled.value) {
      let randomIndex: number
      randomIndex = Math.floor(Math.random() * currentQueue.value.tracks.length)
      let whileCounter = 0
      while (previousIndexes.includes(randomIndex) && whileCounter < 50) {
        randomIndex = Math.floor(Math.random() * currentQueue.value.tracks.length)
        whileCounter++
      }
      previousIndexes.push(randomIndex)
      nextTrack = currentQueue.value.tracks[randomIndex]
      currentQueue.value.position = randomIndex
      if (nextTrack !== undefined) {
        setCurrentlyPlayingTrack(nextTrack)
        return nextTrack
      }
    }
    else if (repeatStatus.value === 'all' && currentIndex === currentQueue.value.tracks.length - 1) {
      nextTrack = currentQueue.value.tracks[0]
      currentQueue.value.position = 0
      if (nextTrack !== undefined) {
        setCurrentlyPlayingTrack(nextTrack)
        return nextTrack
      }
    }
    else if (repeatStatus.value === '1') {
      nextTrack = currentQueue.value.tracks[currentIndex]
      currentQueue.value.position = currentIndex
      if (nextTrack !== undefined) {
        setCurrentlyPlayingTrack(nextTrack)
        if (audioElement.value) {
          audioElement.value.currentTime = 0
          await audioElement.value.play()
        }
        return nextTrack
      }
    }
    else {
      if (currentIndex < currentQueue.value.tracks.length - 1) {
        nextTrack = currentQueue.value.tracks[currentIndex + 1]
        currentQueue.value.position = currentIndex + 1
      }
      if (nextTrack !== undefined) {
        setCurrentlyPlayingTrack(nextTrack)
        return nextTrack
      }
    }
  }
  else {
    const randomTrack = await getRandomTrack()
    return randomTrack
  }
}

export async function handlePreviousTrack(): Promise<SubsonicSong | undefined> {
  if (currentQueue.value && currentQueue.value.tracks.length) {
    const currentIndex = currentQueue.value.position
    let prevTrack: SubsonicSong | undefined
    if (shuffleEnabled.value) {
      if (previousIndexes.length) {
        previousIndexes.pop()
      }
      const previousIndex = previousIndexes[previousIndexes.length - 1]
      const prevTrack = currentQueue.value.tracks[previousIndex]
      currentQueue.value.position = previousIndex
      if (prevTrack !== undefined) {
        setCurrentlyPlayingTrack(prevTrack)
      }
      return prevTrack
    }
    if (repeatStatus.value === 'all') {
      prevTrack = currentQueue.value.tracks[currentQueue.value.tracks.length - 1]
      currentQueue.value.position = currentQueue.value.tracks.length - 1
      return prevTrack
    }
    else if (repeatStatus.value === '1') {
      prevTrack = currentQueue.value.tracks[currentIndex]
      return prevTrack
    }
    else {
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
  }
  else {
    const randomTrack = await getRandomTrack()
    return randomTrack
  }
}

export function togglePlayback() {
  if (!audioElement.value) {
    console.error('Audio element not found')
    return
  }
  else if (!(currentQueue.value && currentQueue.value.tracks.length) && routeTracks.value.length) {
    setCurrentQueue(routeTracks.value)
  }
  else if ((currentQueue.value && currentQueue.value.tracks.length) && !currentlyPlayingItem.value.track) {
    setCurrentQueue(currentQueue.value?.tracks)
  }

  if (isPlaying.value) {
    audioElement.value.pause()
    isPlaying.value = false
  }
  else {
    if (currentlyPlayingItem.value.track || currentlyPlayingItem.value.podcastEpisode) {
      void audioElement.value.play()
      isPlaying.value = true
    }
  }
}

export function stopPlayback() {
  if (audioElement.value) {
    if (!isPlaying.value || audioElement.value.currentTime === 0) {
      currentlyPlayingItem.value = {}
      audioElement.value.currentTime = 0
      audioElement.value.removeAttribute('src')
      clearQueue()
    }
    audioElement.value.pause()
    audioElement.value.load()
  }

  isPlaying.value = false
}

export function seek(seekSeconds: number) {
  if (audioElement.value) {
    audioElement.value.currentTime = seekSeconds
  }
}

export function updateProgress() {
  if (!audioElement.value) {
    return
  }
  currentTime.value = audioElement.value.currentTime

  if (currentlyPlayingItem.value.track && !playcountPosted.value) {
    const halfwayPoint = currentlyPlayingItem.value.track.duration / 2
    if (currentTime.value >= halfwayPoint) {
      void postPlaycount(currentlyPlayingItem.value.track.musicBrainzId)
      playcountPosted.value = true
    }
  }
  else if (currentlyPlayingItem.value.podcastEpisode && !playcountPosted.value) {
    const halfwayPoint = Number(currentlyPlayingItem.value.podcastEpisode.duration) / 2
    if (currentTime.value >= halfwayPoint) {
      void postPlaycount(currentlyPlayingItem.value.podcastEpisode.streamId)
      playcountPosted.value = true
    }
  }
}
