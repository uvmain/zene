import { describe, expect, it, vi } from 'vitest'
import { useDebug } from '../useDebug'

const { debugLog, toggleDebug, useDebugBool } = useDebug()

describe('useDebug', () => {
  afterAll(() => {
    useDebugBool.value = false // Reset debug state after tests
  })

  it('should be defined', () => {
    expect(useDebug).toBeDefined()
  })

  describe('toggleDebug', () => {
    it('should toggle debug state', () => {
      expect(useDebugBool.value).toBe(false)
      toggleDebug()
      expect(useDebugBool.value).toBe(true)
    })
  })

  it('should log a text string to console', () => {
    const consoleSpy = vi.spyOn(console, 'log')
    const text = 'Test log message'
    debugLog(text)
    expect(consoleSpy).toHaveBeenCalledWith(`[DEBUG] ${text}`)
  })
})
