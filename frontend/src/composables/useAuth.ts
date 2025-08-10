import type { SubsonicUserResponse } from '~/types/getUser'
import { useLocalStorage, useSessionStorage } from '@vueuse/core'
import { useBackendFetch } from './useBackendFetch'

const userLoginState = useSessionStorage('untrustedLoginState', false)
const userIsAdminState = useSessionStorage('untrustedUserIsAdmin', false)
const userUsername = useLocalStorage('userUsername', '')
const userSalt = useLocalStorage('userSalt', '')
const userToken = useLocalStorage('userToken', '')
const userApiKey = useLocalStorage('userApiKey', '')
const router = useRouter()

const { openSubsonicFetchRequest } = useBackendFetch()

export function useAuth() {
  const checkIfLoggedIn = async (): Promise<boolean> => {
    try {
      const queryParams = new URLSearchParams()
      if (userApiKey.value) {
        queryParams.append('apiKey', userApiKey.value)
      }
      else if (userSalt.value && userToken.value) {
        queryParams.append('s', userSalt.value)
        queryParams.append('t', userToken.value)
      }
      else {
        await router.push('/login')
      }
      queryParams.append('u', userUsername.value)

      queryParams.append('f', 'json')
      queryParams.append('v', '1.16.0')
      queryParams.append('c', 'zene-frontend')

      const response = await openSubsonicFetchRequest('getUser.view', {
        method: 'GET',
        queryParams,
      })

      const json = await response.json() as SubsonicUserResponse
      const subsonicResponse = json['subsonic-response']
      if (subsonicResponse.error) {
        throw new Error(subsonicResponse.error.message)
      }
      userLoginState.value = subsonicResponse.status === 'ok'
      userIsAdminState.value = subsonicResponse.user.adminRole === 'true'
      return userLoginState.value
    }
    catch {
      userLoginState.value = false
      userIsAdminState.value = false
      return false
    }
  }

  const logout = async () => {
    userUsername.value = ''
    userSalt.value = ''
    userToken.value = ''
    userApiKey.value = ''
    userLoginState.value = false
    userIsAdminState.value = false
  }

  return {
    checkIfLoggedIn,
    logout,
    userLoginState,
    userIsAdminState,
    userUsername,
    userSalt,
    userToken,
    userApiKey,
  }
}
