import { useLocalStorage } from '@vueuse/core'
import { useLogic } from './useLogic'

const randomSeed = useLocalStorage<number>('randomSeed', 0)
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

  return {
    randomSeed,
    refreshRandomSeed,
    getRandomSeed,
  }
}
