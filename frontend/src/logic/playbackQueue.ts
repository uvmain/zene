import type { PlayItem } from '~/types'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicIndexArtist } from '~/types/subsonicArtist'
import type { SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import type { SubsonicSong } from '~/types/subsonicSong'
import { audioElement, clearActiveAudio, seek as elementSeek, playWhenReady } from '~/logic/audioElement'
import { fetchAlbum, fetchArtistTopSongs, fetchRandomTracks } from '~/logic/backendFetch'
import { getAuthenticatedTrackUrl } from '~/logic/common'
import { postPlaycount } from '~/logic/playerUtils'
import { routeTracks } from '~/logic/routeTracks'
import { repeatStatus, shuffleEnabled } from '~/stores/main'
import { episodeIsStored, getStoredEpisode } from '~/stores/podcastStore'
import { debugLog } from './logger'

export const currentlyPlayingItem = ref<PlayItem>({})
export const currentQueue = ref<SubsonicSong[]>([])
export const currentQueuePosition = ref<number>(0)
export const isPlaying = ref(false)
export const playcountPosted = ref(false)
export const currentTime = ref(0)
export const trackUrl = ref('')
export const previousIndexes = ref<number[]>([])

interface QueueTransition {
  track: SubsonicSong
  queuePosition: number
  previousIndexes: number[]
  restartCurrentTrack?: boolean
}

interface PlaybackOptions {
  autoPlay?: boolean
  queuePosition?: number
  previousIndexes?: number[]
  restartCurrentTrack?: boolean
}

let currentHalfwayPoint = 0

function setHalfwayPoint(playItem: PlayItem) {
  if (playItem.track) {
    currentHalfwayPoint = playItem.track.duration / 2
    return
  }

  if (playItem.podcastEpisode) {
    currentHalfwayPoint = Number.parseInt(playItem.podcastEpisode.duration, 10) / 2
    return
  }

  currentHalfwayPoint = 0
}

function applyPlaybackState(playItem: PlayItem, src: string, options: PlaybackOptions = {}) {
  currentlyPlayingItem.value = playItem
  trackUrl.value = src
  playcountPosted.value = false
  setHalfwayPoint(playItem)

  if (options.queuePosition !== undefined) {
    currentQueuePosition.value = options.queuePosition
  }

  if (options.previousIndexes) {
    previousIndexes.value = [...options.previousIndexes]
  }
}

async function resolvePodcastUrl(episode: SubsonicPodcastEpisode): Promise<string> {
  const stored = await episodeIsStored(episode.streamId)
  if (!stored) {
    return getAuthenticatedTrackUrl(episode.streamId, true)
  }

  const blob = await getStoredEpisode(episode.streamId)
  return URL.createObjectURL(blob)
}

async function startPlayback(playItem: PlayItem, src: string, options: PlaybackOptions = {}) {
  const { autoPlay = true, restartCurrentTrack = false } = options

  applyPlaybackState(playItem, src, options)

  if (!autoPlay) {
    return
  }

  if (restartCurrentTrack && audioElement.value) {
    audioElement.value.currentTime = 0
  }

  let started = false

  if (!started) {
    started = await playWhenReady(playItem, src)
  }

  if (!started) {
    isPlaying.value = false
  }
}

async function playTrack(track: SubsonicSong, options: PlaybackOptions = {}) {
  const src = getAuthenticatedTrackUrl(track.musicBrainzId)
  await startPlayback({ track }, src, options)
}

async function playQueueTransition(transition: QueueTransition, options: PlaybackOptions = {}) {
  await playTrack(transition.track, {
    ...options,
    queuePosition: transition.queuePosition,
    previousIndexes: transition.previousIndexes,
    restartCurrentTrack: transition.restartCurrentTrack,
  })
}

function findRandomQueueIndex(queueLength: number, priorIndexes: number[]): number {
  let randomIndex = Math.floor(Math.random() * queueLength)
  let whileCounter = 0

  while (priorIndexes.includes(randomIndex) && whileCounter < 50) {
    randomIndex = Math.floor(Math.random() * queueLength)
    whileCounter++
  }

  return randomIndex
}

function getNextQueueTransition(): QueueTransition | undefined {
  if (currentQueue.value.length === 0) {
    debugLog('No queue or empty queue when trying to get next track')
    return undefined
  }

  const queue = currentQueue.value
  const currentIndex = currentQueuePosition.value

  if (shuffleEnabled.value) {
    const nextPreviousIndexes = [...previousIndexes.value, currentIndex]
    if (nextPreviousIndexes.length >= 50) {
      nextPreviousIndexes.shift()
    }

    const randomIndex = findRandomQueueIndex(queue.length, nextPreviousIndexes)
    return {
      track: queue[randomIndex],
      queuePosition: randomIndex,
      previousIndexes: nextPreviousIndexes,
    }
  }

  if (repeatStatus.value === 'all' && currentIndex === queue.length - 1) {
    return {
      track: queue[0],
      queuePosition: 0,
      previousIndexes: [...previousIndexes.value],
    }
  }

  if (repeatStatus.value === '1') {
    return {
      track: queue[currentIndex],
      queuePosition: currentIndex,
      previousIndexes: [...previousIndexes.value],
      restartCurrentTrack: true,
    }
  }

  if (currentIndex < queue.length - 1) {
    return {
      track: queue[currentIndex + 1],
      queuePosition: currentIndex + 1,
      previousIndexes: [...previousIndexes.value],
    }
  }

  debugLog('Reached end of queue with no repeat when trying to get next track')

  return undefined
}

function getPreviousQueueTransition(): QueueTransition | undefined {
  if (currentQueue.value.length === 0) {
    return undefined
  }

  const queue = currentQueue.value
  const currentIndex = currentQueuePosition.value

  if (shuffleEnabled.value && previousIndexes.value.length > 0) {
    const nextPreviousIndexes = [...previousIndexes.value]
    const previousIndex = nextPreviousIndexes.pop() as number
    return {
      track: queue[previousIndex],
      queuePosition: previousIndex,
      previousIndexes: nextPreviousIndexes,
    }
  }

  if (repeatStatus.value === '1') {
    return {
      track: queue[currentIndex],
      queuePosition: currentIndex,
      previousIndexes: [...previousIndexes.value],
      restartCurrentTrack: true,
    }
  }

  if (currentIndex > 0) {
    return {
      track: queue[currentIndex - 1],
      queuePosition: currentIndex - 1,
      previousIndexes: [...previousIndexes.value],
    }
  }

  if (repeatStatus.value === 'all' || currentIndex === 0) {
    return {
      track: queue.at(-1) as SubsonicSong,
      queuePosition: queue.length - 1,
      previousIndexes: [...previousIndexes.value],
    }
  }

  return undefined
}

export function handlePlay(track: SubsonicSong) {
  if (currentQueue.value?.some(queueTrack => queueTrack.id === track.id)) {
    void setCurrentlyPlayingTrack(track)
  }
  else if (routeTracks.value?.some(queueTrack => queueTrack.musicBrainzId === track.musicBrainzId)) {
    setCurrentQueue(routeTracks.value)
    void setCurrentlyPlayingTrack(track)
  }
  else {
    void play({ track })
  }
}

export async function setCurrentlyPlayingTrack(track: SubsonicSong, autoPlay = true) {
  const queuePosition = currentQueue.value?.indexOf(track)
  await playTrack(track, {
    autoPlay,
    queuePosition: queuePosition !== undefined && queuePosition >= 0 ? queuePosition : undefined,
  })
}

export async function setCurrentlyPlayingPodcast(episode: SubsonicPodcastEpisode) {
  const src = await resolvePodcastUrl(episode)
  clearQueue()
  await startPlayback({ podcastEpisode: episode }, src)
}

function setCurrentQueue(tracks: SubsonicSong[]) {
  clearQueue()
  if (tracks.length === 0) {
    return
  }

  let index = 0
  if (shuffleEnabled.value) {
    index = Math.floor(Math.random() * tracks.length)
  }

  currentQueue.value = tracks
  currentQueuePosition.value = index
  void setCurrentlyPlayingTrack(tracks[index])
}

export function clearQueue() {
  previousIndexes.value = []
  currentQueue.value = []
}

interface PlayOptions {
  artist?: SubsonicIndexArtist
  album?: SubsonicAlbum
  track?: SubsonicSong
  podcastEpisode?: SubsonicPodcastEpisode
}

export async function play(playOptions: PlayOptions) {
  if (playOptions.track) {
    await setCurrentlyPlayingTrack(playOptions.track)
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
  else {
    setCurrentQueue(routeTracks.value ?? [])
  }
}

export async function getRandomTracks(size: number = 10): Promise<SubsonicSong[]> {
  const randomTracks = await fetchRandomTracks({ limit: size })
  setCurrentQueue(randomTracks)
  return randomTracks
}

export async function handleNextTrack(): Promise<SubsonicSong | undefined> {
  const transition = getNextQueueTransition()

  if (transition) {
    await playQueueTransition(transition)
    return transition.track
  }
  else {
    return undefined
  }
}

export async function handlePreviousTrack(): Promise<SubsonicSong | undefined> {
  const transition = getPreviousQueueTransition()
  if (!transition) {
    return undefined
  }

  await playQueueTransition(transition)
  return transition.track
}

export function togglePlayback() {
  if (!audioElement.value) {
    console.error('Audio element not found')
    return
  }

  if (currentQueue.value.length === 0 && routeTracks.value.length) {
    setCurrentQueue(routeTracks.value)
    return
  }

  if (currentQueue.value.length > 0 && !currentlyPlayingItem.value.track && !currentlyPlayingItem.value.podcastEpisode) {
    setCurrentQueue(currentQueue.value)
    return
  }

  if (isPlaying.value) {
    audioElement.value.pause()
    isPlaying.value = false
    return
  }

  if (currentlyPlayingItem.value.track || currentlyPlayingItem.value.podcastEpisode) {
    void audioElement.value.play()
    isPlaying.value = true
  }
}

export function stopPlayback() {
  if (audioElement.value) {
    if (!isPlaying.value || audioElement.value.currentTime === 0) {
      currentlyPlayingItem.value = {}
      trackUrl.value = ''
      clearQueue()
      clearActiveAudio()
    }
    else {
      audioElement.value.pause()
      audioElement.value.load()
    }
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
  elementSeek(seekSeconds)
}
