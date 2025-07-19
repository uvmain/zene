import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useDebug } from '../useDebug'

const { debugLog, toggleDebug, useDebugBool } = useDebug()

// Mock backend fetch if needed
vi.mock('../useBackendFetch', () => ({
  useBackendFetch: () => ({
    backendFetchRequest: vi.fn().mockResolvedValue({
      json: async () => Promise.resolve({}),
      ok: true,
    }),
  }),
}))

describe('useDebug', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })
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
