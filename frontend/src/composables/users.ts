import type { SubsonicResponse, SubsonicUserResponse, SubsonicUsersResponse } from '~/types/subsonic'
import type { SubsonicUser } from '~/types/subsonicUser'
import { openSubsonicFetchRequest } from './backendFetch'

export async function fetchCurrentUser(): Promise<SubsonicUser> {
  const response = await openSubsonicFetchRequest('getUser') as SubsonicUserResponse
  return response.user
}

export async function fetchUsers(): Promise<SubsonicUser[]> {
  const response = await openSubsonicFetchRequest('getUsers') as SubsonicUsersResponse
  return response.users.user
}

export async function createUser(user: SubsonicUser): Promise<void> {
  const formData = new FormData()
  formData.append('username', user.username)
  formData.append('adminRole', user.adminRole.toString())
  formData.append('password', user.password ?? '')
  formData.append('email', user.email)
  const response = await openSubsonicFetchRequest('createUser', {
    body: formData,
  })
  if (response.status !== 'ok') {
    throw new Error(response.error?.message ?? 'Unknown error')
  }
}

export async function updateUser(user: SubsonicUser): Promise<void> {
  const formData = new FormData()
  formData.append('username', user.username)
  formData.append('adminRole', user.adminRole.toString())
  if (user.password !== undefined && user.password.length > 0) {
    formData.append('password', user.password)
  }
  if (user.email) {
    formData.append('email', user.email)
  }
  const response = await openSubsonicFetchRequest('updateUser.view', {
    body: formData,
  })
  if (response.status !== 'ok') {
    throw new Error(response.error?.message ?? 'Unknown error')
  }
}

export async function deleteUser(user: SubsonicUser) {
  const formData = new FormData()
  formData.append('username', user.username)
  const response = await openSubsonicFetchRequest('deleteUser.view', {
    body: formData,
  })
  if (response.status !== 'ok') {
    throw new Error(response.error?.message ?? 'Unknown error')
  }
}
