import type { SessionCheck } from '~/types/auth'
import { useSessionStorage } from '@vueuse/core'
import { useBackendFetch } from './useBackendFetch'

const userLoginState = useSessionStorage('untrustedLoginState', false)
const userIsAdminState = useSessionStorage('untrustedUserIsAdmin', false)
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
      userIsAdminState.value = json.isAdmin
      return userLoginState.value
    }
    catch {
      userLoginState.value = false
      userIsAdminState.value = false
      return false
    }
  }

  const logout = async () => {
    try {
      await backendFetchRequest('logout', {
        method: 'GET',
        credentials: 'include',
      })
      userLoginState.value = false
    }
    catch {
      // Don't update userLoginState if logout request fails
    }
  }

  return {
    checkIfLoggedIn,
    logout,
    userLoginState,
    userIsAdminState,
  }
}
