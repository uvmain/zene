import type { SessionCheck } from '../types/auth'
import { useSessionStorage } from '@vueuse/core'
import { backendFetchRequest } from './fetchFromBackend'

const userLoginState = useSessionStorage('untrustedLoginState', false)

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

export async function login(username: string, password: string) {
  const formData = new FormData()
  formData.append('username', username)
  formData.append('password', password)

  const response = await backendFetchRequest('login', {
    body: formData,
    method: 'POST',
  })
  const jsonData = await response.json() as SessionCheck
  userLoginState.value = jsonData.loggedIn
}

export async function logout() {
  await backendFetchRequest('logout', {
    method: 'GET',
    credentials: 'include',
  })
  userLoginState.value = false
}
