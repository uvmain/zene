import { describe, expect, it } from 'vitest'
import { useRandomSeed } from '../useRandomSeed'

const { randomSeed, refreshRandomSeed, getRandomSeed } = useRandomSeed()

describe('useRandomSeed', () => {
  describe('randomSeed', () => {
    it('should have a default value of 0', () => {
      expect(randomSeed.value).toBe(0)
    })
  })

  describe('getRandomSeed', () => {
    it('should return a new random seed', () => {
      const initialSeed = randomSeed.value
      const result = getRandomSeed()
      expect(result).not.toBe(initialSeed)
    })
  })

  describe ('refreshRandomSeed', () => {
    it('should generate a new random seed', () => {
      const initialSeed = randomSeed.value
      refreshRandomSeed()
      expect(randomSeed.value).not.toBe(initialSeed)
    })
  })
})
