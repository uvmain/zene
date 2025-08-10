import type { SubsonicUserResponse } from '~/types/getUser'
import { useLocalStorage, useSessionStorage } from '@vueuse/core'
import { useBackendFetch } from './useBackendFetch'

const userLoginState = useSessionStorage('untrustedLoginState', false)
const userIsAdminState = useSessionStorage('untrustedUserIsAdmin', false)
const userUsername = useLocalStorage('userUsername', '')
const userSalt = useLocalStorage('userSalt', '')
const userToken = useLocalStorage('userToken', '')
const userApiKey = useLocalStorage('userApiKey', '')

export function useAuth() {
  const logout = async () => {
    userUsername.value = ''
    userSalt.value = ''
    userToken.value = ''
    userApiKey.value = ''
    userLoginState.value = false
    userIsAdminState.value = false
  }

  return {
    logout,
    userLoginState,
    userIsAdminState,
    userUsername,
    userSalt,
    userToken,
    userApiKey,
  }
}
