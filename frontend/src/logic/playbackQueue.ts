import type { Queue } from '~/types'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicIndexArtist } from '~/types/subsonicArtist'
import type { SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchAlbum, fetchArtistTopSongs, fetchRandomTracks } from '~/logic/backendFetch'

export const currentlyPlayingTrack = ref<SubsonicSong | undefined>()
export const currentlyPlayingPodcastEpisode = ref<SubsonicPodcastEpisode | undefined>()
export const currentQueue = ref<Queue | undefined>()

const previousIndexes = ref<number[]>([])

export function resetCurrentlyPlayingTrack() {
  currentlyPlayingTrack.value = undefined
}

export function setCurrentlyPlayingTrack(track: SubsonicSong) {
  if (currentQueue.value) {
    const index = currentQueue.value.tracks.indexOf(track)
    currentQueue.value.position = index
  }
  currentlyPlayingTrack.value = track
}

export function setCurrentlyPlayingPodcastEpisode(episode: SubsonicPodcastEpisode) {
  currentlyPlayingTrack.value = undefined
  currentlyPlayingPodcastEpisode.value = episode
  clearQueue()
}

export function setCurrentQueue(tracks: SubsonicSong[], playRandomTrack: boolean = false) {
  currentQueue.value = {
    tracks,
    position: 0,
  }
  let index = 0
  if (playRandomTrack && tracks.length > 0) {
    index = Math.floor(Math.random() * tracks.length)
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

export async function getRandomTracks(size: number = 10, shuffled = false): Promise<SubsonicSong[]> {
  previousIndexes.value = []
  const randomTracks = await fetchRandomTracks(size)
  setCurrentQueue(randomTracks, shuffled)
  return randomTracks
}

export async function getNextTrack(): Promise<SubsonicSong | undefined> {
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

export async function getShuffledTrack(): Promise<SubsonicSong | undefined> {
  let nextTrack: SubsonicSong | undefined
  if (currentQueue.value && currentQueue.value.tracks.length) {
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
  }
  if (nextTrack !== undefined) {
    setCurrentlyPlayingTrack(nextTrack)
  }
  return nextTrack
}

export async function getPreviousTrack(): Promise<SubsonicSong | undefined> {
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

export async function getPreviousShuffledTrack(): Promise<SubsonicSong | undefined> {
  if (currentQueue.value && currentQueue.value.tracks.length) {
    previousIndexes.value.pop()
    const previousIndex = previousIndexes.value[previousIndexes.value.length - 1]
    const prevTrack = currentQueue.value.tracks[previousIndex]
    currentQueue.value.position = previousIndex
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
