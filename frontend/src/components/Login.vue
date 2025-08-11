<script setup lang="ts">
import type { SubsonicUserResponse } from '~/types/subsonicUser'
import { md5 } from 'js-md5'
import { useAuth } from '~/composables/useAuth'
import { useBackendFetch } from '~/composables/useBackendFetch'

const router = useRouter()
const { openSubsonicFetchRequest } = useBackendFetch()
const { userLoginState, userIsAdminState, userUsername, userSalt, userToken } = useAuth()

function generateSalt(length = 6) {
  const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let result = ''
  const array = new Uint32Array(length)
  crypto.getRandomValues(array)
  for (let i = 0; i < length; i++) {
    result += chars[array[i] % chars.length]
  }
  return result
}

function md5Hash(password: string, salt: string): string {
  const messageToHash = password + salt
  const hash = md5.create()
  hash.update(messageToHash)
  return hash.hex()
}

async function generateSaltAndTokenForOpenSubsonic(username: string, password: string): Promise<void> {
  const salt = generateSalt()
  const token = md5Hash(password, salt)
  userUsername.value = username
  userSalt.value = salt
  userToken.value = token
}

async function login(username: string, password: string) {
  await generateSaltAndTokenForOpenSubsonic(username, password)
  const response = await openSubsonicFetchRequest('getUser.view')
  const jsonData = await response.json() as SubsonicUserResponse
  userLoginState.value = jsonData['subsonic-response'].status === 'ok'
  userIsAdminState.value = jsonData['subsonic-response'].user.adminRole === 'true'
  if (userLoginState.value) {
    router.push('/')
  }
}

const username = ref('')
const password = ref('')
</script>

<template>
  <div class="fixed left-0 top-0 z-999 size-full animate-fade-in backdrop-blur-xl">
    <div class="mx-auto mb-auto mt-150px max-w-sm w-full border rounded px-6 pb-6 pt-4 md:w-300px md:px-30px md:pb-30px md:pt-20px">
      <div class="flex flex-col gap-4 p-4 md:w-300 md:p-6">
        <form class="flex flex-col gap-3">
          <div class="flex flex-col gap-2 md:flex-row md:items-center md:gap-2">
            <label for="username" class="min-w-20 text-sm md:text-base">Username:</label>
            <input id="username" v-model="username" type="text" name="username" autocomplete="username" class="border rounded px-3 py-2 md:px-2 md:py-1">
          </div>
          <div class="flex flex-col gap-2 md:flex-row md:items-center md:gap-2">
            <label for="password" class="min-w-20 text-sm md:text-base">Password:</label>
            <input id="password" v-model="password" type="password" name="password" autocomplete="current-password" class="border rounded px-3 py-2 md:px-2 md:py-1" @keydown.enter="login(username, password)">
          </div>
        </form>
      </div>
      <div class="flex justify-center gap-4">
        <button aria-label="login" class="min-h-11 rounded bg-zene-400 px-6 py-3 text-white font-semibold transition-colors hover:bg-zene-200 md:px-4 md:py-2" @click="login(username, password)">
          Login
        </button>
      </div>
    </div>
  </div>
</template>
