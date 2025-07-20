import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useAuth } from '../useAuth'

const { checkIfLoggedIn, logout, userLoginState } = useAuth()

beforeEach(() => {
  vi.clearAllMocks()
})

describe('useAuth', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should be defined', () => {
    expect(useAuth).toBeDefined()
  })

  describe('userLoginState', () => {
    it('should have a default value of false', () => {
      expect(userLoginState.value).toBe(false)
    })
  })

  describe('checkIfLoggedIn', () => {
    it('should return true if user is logged in', async () => {
      const mockResponse = { loggedIn: true }
      vi.spyOn(globalThis, 'fetch').mockResolvedValue({
        json: async () => Promise.resolve(mockResponse),
      } as Response)

      const result = await checkIfLoggedIn()
      expect(result).toBe(true)
      expect(userLoginState.value).toBe(true)
    })

    it('should return false if user is not logged in', async () => {
      const mockResponse = { loggedIn: false }
      vi.spyOn(globalThis, 'fetch').mockResolvedValue({
        json: async () => Promise.resolve(mockResponse),
      } as Response)

      const result = await checkIfLoggedIn()
      expect(result).toBe(false)
      expect(userLoginState.value).toBe(false)
    })

    it('should handle errors and set userLoginState to false', async () => {
      vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'))

      const result = await checkIfLoggedIn()
      expect(result).toBe(false)
      expect(userLoginState.value).toBe(false)
    })
  })

  describe('logout', () => {
    it('should log out the user and update userLoginState', async () => {
      const mockResponse = { success: true }
      vi.spyOn(globalThis, 'fetch').mockResolvedValue({
        json: async () => Promise.resolve(mockResponse),
      } as Response)

      await logout()
      expect(userLoginState.value).toBe(false)
    })

    it('should handle errors and not update userLoginState', async () => {
      vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'))
      userLoginState.value = true
      await logout()
      expect(userLoginState.value).toBe(true)
    })
  })
})
