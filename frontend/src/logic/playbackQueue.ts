import type { Queue } from '~/types'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicIndexArtist } from '~/types/subsonicArtist'
import type { SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchAlbum, fetchArtistTopSongs, fetchRandomTracks } from '~/logic/backendFetch'
import { debugLog } from '~/logic/logger'
import { postPlaycount } from '~/logic/playCounts'
import { repeatStatus, routeTracks, shuffleEnabled } from '~/logic/store'

type AudioElement = HTMLAudioElement | null | undefined

export const currentlyPlayingTrack = ref<SubsonicSong | undefined>()
export const currentlyPlayingPodcastEpisode = ref<SubsonicPodcastEpisode | undefined>()
export const currentQueue = ref<Queue | undefined>()
export const isPlaying = ref(false)
export const playcountPosted = ref(false)
export const currentTime = ref(0)
export const previousVolume = ref(1)
export const currentVolume = ref(1)
export const trackUrl = ref('')
export const audioElement = ref<AudioElement>(null)
const previousIndexes = ref<number[]>([])

export function resetCurrentlyPlayingTrack() {
  currentlyPlayingTrack.value = undefined
  currentlyPlayingPodcastEpisode.value = undefined
  if (audioElement.value) {
    audioElement.value.currentTime = 0
  }
}

export function setCurrentlyPlayingTrack(track: SubsonicSong) {
  if (currentQueue.value) {
    const index = currentQueue.value.tracks.indexOf(track)
    currentQueue.value.position = index
  }
  currentlyPlayingPodcastEpisode.value = undefined
  currentlyPlayingTrack.value = track
}

export function setCurrentlyPlayingPodcastEpisode(episode: SubsonicPodcastEpisode) {
  currentlyPlayingTrack.value = undefined
  currentlyPlayingPodcastEpisode.value = episode
  clearQueue()
}

export function setCurrentQueue(tracks: SubsonicSong[]) {
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
  previousIndexes.value = []
  currentQueue.value = undefined
}

async function getRandomTrack(): Promise<SubsonicSong> {
  const randomTracks = await fetchRandomTracks(1)
  const randomTrack = randomTracks[0]
  setCurrentlyPlayingTrack(randomTrack)
  clearQueue()
  return randomTrack
}

export async function play(artist?: SubsonicIndexArtist, album?: SubsonicAlbum, track?: SubsonicSong, podcastEpisode?: SubsonicPodcastEpisode) {
  if (track) {
    setCurrentlyPlayingTrack(track)
    clearQueue()
  }
  else if (album) {
    const tracks = await fetchAlbum(album.id).then(fetchedAlbum => fetchedAlbum.song)
    setCurrentQueue(tracks)
  }
  else if (artist) {
    const tracks = await fetchArtistTopSongs(artist.id, 100)
    setCurrentQueue(tracks)
  }
  else if (podcastEpisode) {
    setCurrentlyPlayingPodcastEpisode(podcastEpisode)
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
      while (previousIndexes.value.includes(randomIndex) && whileCounter < 50) {
        randomIndex = Math.floor(Math.random() * currentQueue.value.tracks.length)
        whileCounter++
      }
      previousIndexes.value.push(randomIndex)
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
      if (previousIndexes.value.length) {
        previousIndexes.value.pop()
      }
      const previousIndex = previousIndexes.value[previousIndexes.value.length - 1]
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

export function toggleShuffle() {
  shuffleEnabled.value = !shuffleEnabled.value
}

export function toggleMute() {
  if (audioElement.value) {
    debugLog('Changing volume')
    if (audioElement.value.volume !== 0) {
      previousVolume.value = audioElement.value.volume
      audioElement.value.volume = 0
      currentVolume.value = 0
    }
    else {
      audioElement.value.volume = previousVolume.value
      currentVolume.value = previousVolume.value
    }
  }
}

export async function togglePlayback() {
  if (!audioElement.value) {
    console.error('Audio element not found')
    return
  }
  else if (!(currentQueue.value && currentQueue.value.tracks.length) && routeTracks.value.length) {
    setCurrentQueue(routeTracks.value)
  }
  else if ((currentQueue.value && currentQueue.value.tracks.length) && !currentlyPlayingTrack.value) {
    setCurrentQueue(currentQueue.value?.tracks)
  }

  if (isPlaying.value) {
    audioElement.value.pause()
    isPlaying.value = false
  }
  else {
    if (currentlyPlayingTrack.value || currentlyPlayingPodcastEpisode.value) {
      await audioElement.value.play()
      isPlaying.value = true
    }
  }
}

export async function stopPlayback() {
  if (audioElement.value) {
    if (audioElement.value.currentTime < 1) {
      resetCurrentlyPlayingTrack()
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

export async function updateProgress() {
  if (!audioElement.value) {
    return
  }
  currentTime.value = audioElement.value.currentTime

  if (currentlyPlayingTrack.value && !playcountPosted.value) {
    const halfwayPoint = currentlyPlayingTrack.value.duration / 2
    if (currentTime.value >= halfwayPoint) {
      await postPlaycount(currentlyPlayingTrack.value.musicBrainzId)
      playcountPosted.value = true
    }
  }
  else if (currentlyPlayingPodcastEpisode.value && !playcountPosted.value) {
    const halfwayPoint = Number(currentlyPlayingPodcastEpisode.value.duration) / 2
    if (currentTime.value >= halfwayPoint) {
      await postPlaycount(currentlyPlayingPodcastEpisode.value.streamId)
      playcountPosted.value = true
    }
  }
}

export function volumeInput(volumeString: string) {
  if (!audioElement.value) {
    return
  }
  const volume = Number.parseFloat(volumeString)
  audioElement.value.volume = volume
  currentVolume.value = volume
}

export function toggleRepeat() {
  switch (repeatStatus.value) {
    case 'off':
      repeatStatus.value = '1'
      break
    case '1':
      repeatStatus.value = 'all'
      break
    case 'all':
      repeatStatus.value = 'off'
      break
  }
}
