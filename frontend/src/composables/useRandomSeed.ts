import { useSessionStorage } from '@vueuse/core'
import { useLogic } from './useLogic'

const randomSeed = useSessionStorage<number>('randomSeed', 0)
const randomAlbumSeed = useSessionStorage<number>('randomAlbumSeed', 0)
const randomArtistSeed = useSessionStorage<number>('randomArtistSeed', 0)
const { getRandomInteger } = useLogic()

export function useRandomSeed() {
  const refreshRandomSeed = (): number => {
    randomSeed.value = getRandomInteger()
    return randomSeed.value
  }

  const getRandomSeed = (): number => {
    if (randomSeed.value === 0) {
      randomSeed.value = getRandomInteger()
    }
    return randomSeed.value
  }

  const refreshRandomAlbumSeed = (): number => {
    randomAlbumSeed.value = getRandomInteger()
    return randomAlbumSeed.value
  }

  const getRandomAlbumSeed = (): number => {
    if (randomAlbumSeed.value === 0) {
      randomAlbumSeed.value = getRandomInteger()
    }
    return randomSeed.value
  }

  const refreshRandomArtistSeed = (): number => {
    randomArtistSeed.value = getRandomInteger()
    return randomArtistSeed.value
  }

  const getRandomArtistSeed = (): number => {
    if (randomArtistSeed.value === 0) {
      randomArtistSeed.value = getRandomInteger()
    }
    return randomArtistSeed.value
  }

  return {
    randomSeed,
    refreshRandomSeed,
    getRandomSeed,
    randomAlbumSeed,
    refreshRandomAlbumSeed,
    getRandomAlbumSeed,
    randomArtistSeed,
    refreshRandomArtistSeed,
    getRandomArtistSeed,
  }
}
