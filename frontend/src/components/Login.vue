<script setup lang="ts">
import type { SessionCheck } from '~/types/auth'
import { useAuth } from '~/composables/useAuth'
import { useBackendFetch } from '~/composables/useBackendFetch'

const router = useRouter()
const { backendFetchRequest } = useBackendFetch()
const { checkIfLoggedIn, userLoginState } = useAuth()

async function login(username: string, password: string) {
  const formData = new FormData()
  formData.append('username', username)
  formData.append('password', password)

  const response = await backendFetchRequest('login', {
    body: formData,
    method: 'POST',
  })
  const jsonData = await response.json() as SessionCheck
  userLoginState.value = jsonData.loggedIn
  if (jsonData.loggedIn === true) {
    router.push('/')
  }
}

const username = ref('')
const password = ref('')

onBeforeMount(() => {
  checkIfLoggedIn()
})
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
