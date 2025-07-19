import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useLogic } from '../useLogic'

const { inputText, closeSearch, niceDate, formatTime } = useLogic()

describe('logic', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('closeSearch', () => {
    it('should clear the search input', () => {
      inputText.value = 'test'
      closeSearch()
      expect(inputText.value).toBe('')
    })

    it('should export closeSearch function', () => {
      expect(closeSearch).toBeDefined()
      expect(typeof closeSearch).toBe('function')
    })
  })

  it('should export niceDate function', () => {
    expect(niceDate).toBeDefined()
    expect(typeof niceDate).toBe('function')
  })

  it('should export formatTime function', () => {
    expect(formatTime).toBeDefined()
    expect(typeof formatTime).toBe('function')
  })

  describe('niceDate', () => {
    it('should format valid date correctly', () => {
      const result = niceDate('2023-01-15T10:30:00Z')
      expect(result).toMatch(/\d{2}\/\d{2}\/\d{4}/)
    })

    it('should return "Invalid Date" for invalid date', () => {
      const result = niceDate('invalid-date')
      expect(result).toBe('Invalid Date')
    })
  })

  describe('formatTime', () => {
    it('should format time correctly', () => {
      expect(formatTime(65)).toBe('1:05')
      expect(formatTime(125)).toBe('2:05')
      expect(formatTime(0)).toBe('0:00')
    })
  })
})
