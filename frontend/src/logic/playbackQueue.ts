import type { PlayItem } from '~/types'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicIndexArtist } from '~/types/subsonicArtist'
import type { SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import type { SubsonicSong } from '~/types/subsonicSong'
import { audioElement, seek as elementSeek, playWhenReady } from '~/logic/audioElement'
import { fetchAlbum, fetchArtistTopSongs, fetchRandomTracks } from '~/logic/backendFetch'
import { castAudio } from '~/logic/castAudio'
import { castPlayer, castPlayerController, chromecastConnected } from '~/logic/castRefs'
import { getAuthenticatedTrackUrl } from '~/logic/common'
import { postPlaycount } from '~/logic/playerUtils'
import { routeTracks } from '~/logic/routeTracks'
import { repeatStatus, shuffleEnabled } from '~/logic/store'
import { episodeIsStored, getStoredEpisode } from '~/stores/usePodcastStore'

export const currentlyPlayingItem = ref<PlayItem>({})
export const currentQueue = ref<SubsonicSong[] | undefined>()
export const currentQueuePosition = ref<number>(0)
export const isPlaying = ref(false)
export const playcountPosted = ref(false)
export const currentTime = ref(0)
export const trackUrl = ref('')

let previousIndexes: number[] = []
let currentHalfwayPoint: number = 0

export function handlePlay(track: SubsonicSong) {
  if (currentQueue.value?.some(queueTrack => queueTrack.id === track.id)) {
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

export function setCurrentlyPlayingTrack(track: SubsonicSong, autoPlay: boolean = true) {
  if (currentQueue.value) {
    const index = currentQueue.value.indexOf(track)
    currentQueuePosition.value = index
  }
  currentlyPlayingItem.value = { track }
  trackUrl.value = getAuthenticatedTrackUrl(track.musicBrainzId)
  playcountPosted.value = false
  currentHalfwayPoint = track.duration / 2
  if (chromecastConnected.value) {
    void castAudio()
  }
  else if (autoPlay) {
    playWhenReady({ track })
  }
}

export async function setCurrentlyPlayingPodcast(episode: SubsonicPodcastEpisode) {
  currentlyPlayingItem.value = { podcastEpisode: episode }
  playcountPosted.value = false
  currentHalfwayPoint = Number.parseInt(episode.duration) / 2

  if (currentlyPlayingItem.value.podcastEpisode) {
    await episodeIsStored(currentlyPlayingItem.value.podcastEpisode.streamId).then(async (stored) => {
      if (!stored && currentlyPlayingItem.value.podcastEpisode) {
        trackUrl.value = getAuthenticatedTrackUrl(currentlyPlayingItem.value.podcastEpisode.streamId, true)
      }
      else if (currentlyPlayingItem.value.podcastEpisode) {
        await getStoredEpisode(currentlyPlayingItem.value.podcastEpisode.streamId).then((blob) => {
          trackUrl.value = URL.createObjectURL(blob)
        })
      }
    })
  }
  clearQueue()
  if (chromecastConnected.value) {
    await castAudio()
  }
  else {
    playWhenReady({ podcastEpisode: episode })
  }
}

function setCurrentQueue(tracks: SubsonicSong[]) {
  clearQueue()
  let index = 0
  if (shuffleEnabled.value && tracks.length > 0) {
    index = Math.floor(Math.random() * tracks.length)
  }
  currentQueue.value = tracks
  currentQueuePosition.value = index
  setCurrentlyPlayingTrack(tracks[index])
}

export function clearQueue() {
  previousIndexes = []
  currentQueue.value = undefined
}

async function getRandomTrack(): Promise<SubsonicSong> {
  const randomTracks = await fetchRandomTracks({ limit: 1 })
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
    await setCurrentlyPlayingPodcast(playOptions.podcastEpisode)
  }
}

export async function getRandomTracks(size: number = 10): Promise<SubsonicSong[]> {
  const randomTracks = await fetchRandomTracks({ limit: size })
  setCurrentQueue(randomTracks)
  return randomTracks
}

export function handleNextTrack() {
  if (currentQueue.value && currentQueue.value.length) {
    const currentIndex = currentQueuePosition.value
    let nextTrack: SubsonicSong | undefined
    if (shuffleEnabled.value) {
      let randomIndex: number
      randomIndex = Math.floor(Math.random() * currentQueue.value.length)
      let whileCounter = 0
      while (previousIndexes.includes(randomIndex) && whileCounter < 50) {
        randomIndex = Math.floor(Math.random() * currentQueue.value.length)
        whileCounter++
      }
      previousIndexes.push(randomIndex)
      nextTrack = currentQueue.value[randomIndex]
      currentQueuePosition.value = randomIndex
      if (nextTrack !== undefined) {
        setCurrentlyPlayingTrack(nextTrack)
      }
    }
    else if (repeatStatus.value === 'all' && currentIndex === currentQueue.value.length - 1) {
      nextTrack = currentQueue.value[0]
      currentQueuePosition.value = 0
      if (nextTrack !== undefined) {
        setCurrentlyPlayingTrack(nextTrack)
      }
    }
    else if (repeatStatus.value === '1') {
      nextTrack = currentQueue.value[currentIndex]
      currentQueuePosition.value = currentIndex
      if (nextTrack !== undefined) {
        setCurrentlyPlayingTrack(nextTrack)
        if (audioElement.value) {
          audioElement.value.currentTime = 0
          void audioElement.value.play()
        }
      }
    }
    else {
      if (currentIndex < currentQueue.value.length - 1) {
        nextTrack = currentQueue.value[currentIndex + 1]
        currentQueuePosition.value = currentIndex + 1
      }
      if (nextTrack !== undefined) {
        setCurrentlyPlayingTrack(nextTrack)
        return nextTrack
      }
    }
  }
}

export async function handlePreviousTrack(): Promise<SubsonicSong | undefined> {
  if (currentQueue.value && currentQueue.value.length) {
    const currentIndex = currentQueuePosition.value
    let prevTrack: SubsonicSong | undefined
    if (shuffleEnabled.value) {
      if (previousIndexes.length) {
        previousIndexes.pop()
      }
      const previousIndex = previousIndexes.at(-1) ?? 0
      const prevTrack = currentQueue.value[previousIndex]
      currentQueuePosition.value = previousIndex
      if (prevTrack !== undefined) {
        setCurrentlyPlayingTrack(prevTrack)
      }
      return prevTrack
    }
    if (repeatStatus.value === 'all') {
      prevTrack = currentQueue.value.at(-1)
      currentQueuePosition.value = currentQueue.value.length - 1
      return prevTrack
    }
    else if (repeatStatus.value === '1') {
      prevTrack = currentQueue.value[currentIndex]
      return prevTrack
    }
    else {
      if (currentIndex > 0) {
        prevTrack = currentQueue.value[currentIndex - 1]
        currentQueuePosition.value = currentIndex - 1
      }
      else {
        prevTrack = currentQueue.value.at(-1)
        currentQueuePosition.value = currentQueue.value.length - 1
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
  if (chromecastConnected.value && castPlayerController.value) {
    castPlayerController.value.playOrPause()
  }

  if (!audioElement.value) {
    console.error('Audio element not found')
    return
  }
  else if (!(currentQueue.value && currentQueue.value.length) && routeTracks.value.length) {
    setCurrentQueue(routeTracks.value)
  }
  else if ((currentQueue.value && currentQueue.value.length) && !currentlyPlayingItem.value.track) {
    setCurrentQueue(currentQueue.value)
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
  if (chromecastConnected.value && castPlayerController.value) {
    castPlayerController.value.stop()
  }

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

export function updateProgress() {
  if (!audioElement.value) {
    return
  }
  currentTime.value = audioElement.value.currentTime

  if (!playcountPosted.value && currentTime.value >= currentHalfwayPoint) {
    void postPlaycount(currentlyPlayingItem.value.track?.musicBrainzId ?? currentlyPlayingItem.value.podcastEpisode?.streamId ?? '')
    playcountPosted.value = true
  }
}

export function seek(seekSeconds: number) {
  if (chromecastConnected.value && castPlayer.value !== null && castPlayerController.value) {
    castPlayer.value.currentTime = castPlayer.value.currentTime + seekSeconds
    castPlayerController.value.seek()
  }
  elementSeek(seekSeconds)
}
