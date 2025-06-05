import type { SessionCheck } from '../types/auth'
import { useSessionStorage } from '@vueuse/core'
import { backendFetchRequest } from './fetchFromBackend'

export const userLoginState = useSessionStorage('untrustedLoginState', false)

export async function checkIfLoggedIn(): Promise<boolean> {
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

export async function logout() {
  await backendFetchRequest('logout', {
    method: 'GET',
    credentials: 'include',
  })
  userLoginState.value = false
}
