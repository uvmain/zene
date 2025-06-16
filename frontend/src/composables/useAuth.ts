import type { SessionCheck } from '../types/auth'
import { useSessionStorage } from '@vueuse/core'
import { useBackendFetch } from './useBackendFetch'

const userLoginState = useSessionStorage('untrustedLoginState', false)
const { backendFetchRequest } = useBackendFetch()

export function useAuth() {
  const checkIfLoggedIn = async (): Promise<boolean> => {
    try {
      const response = await backendFetchRequest('check-session', {
        method: 'GET',
        credentials: 'include',
      })
      const json = await response.json() as SessionCheck
      userLoginState.value = json.loggedIn
      return userLoginState.value
    }
    catch {
      userLoginState.value = false
      return false
    }
  }

  const logout = async () => {
    await backendFetchRequest('logout', {
      method: 'GET',
      credentials: 'include',
    })
    userLoginState.value = false
  }

  return {
    checkIfLoggedIn,
    logout,
    userLoginState,
  }
}
