<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core'
import md5 from 'blueimp-md5'
import { fetchApiKeysWithTokenAndSalt, fetchNewApiKeyWithTokenAndSalt } from '~/composables/backendFetch'

const router = useRouter()

const username = ref('')
const token = ref('')
const salt = ref('')
const password = ref('')
const loading = ref(false)
const error = ref<string | null>(null)

const apiKey = useLocalStorage('apiKey', '')

const signInDisabled = computed(() => {
  return username.value.length < 1 || password.value.length < 1 || loading.value
})

async function login() {
  error.value = ''
  loading.value = true
  try {
    salt.value = Math.random().toString(36).slice(2, 10)
    token.value = md5(password.value + salt.value)

    const data = await fetchApiKeysWithTokenAndSalt(username.value, token.value, salt.value)
    if (!data || !data.apiKeys) {
      throw new Error('Login failed')
    }

    if (data.apiKeys.apiKey.length === 0) {
      const newApiKey = await fetchNewApiKeyWithTokenAndSalt(username.value, token.value, salt.value)
      apiKey.value = newApiKey
      router.push('/')
    }
    else {
      const existingApiKey = data?.apiKeys.apiKey[0]?.api_key
      apiKey.value = existingApiKey
      router.push('/')
    }
  }
  catch (e: any) {
    error.value = e?.message || 'Login failed'
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="my-auto flex flex-col items-center justify-center gap-6 background-1 p-4 lg:flex-row lg:gap-12">
    <img
      class="size-full max-w-400px opacity-90"
      src="/minidisk.svg"
      alt="Logo"
      width="200"
      height="200"
    />
    <div class="text-muted">
      <h1>Login to Zene</h1>
      <form class="grid w-400px gap-2" @submit.prevent="login">
        <label for="username">
          Username
        </label>
        <input
          id="username"
          v-model="username"
          type="text"
          class="border-1 border-primary2 rounded background-2 py-2 pl-10 font-semibold focus:border-primary2 dark:border-opacity-60 focus:border-solid focus:shadow-primary2 hover:shadow-lg focus:outline-none"
          autocomplete="username"
          required
          @input="error = null"
        />
        <label for="password">
          Password
        </label>
        <input
          id="password"
          v-model="password"
          type="password"
          class="border-1 border-primary2 rounded background-2 py-2 pl-10 font-semibold opacity-100 focus:border-primary2 dark:border-opacity-60 focus:border-solid focus:shadow-primary2 hover:shadow-lg focus:outline-none"
          autocomplete="current-password"
          required
          @input="error = null"
        />
        <ZButton class="mx-auto mt-4" :disabled="signInDisabled">
          {{ loading ? 'Signing in…' : 'Sign in' }}
        </ZButton>
      </form>
      <div v-if="error" class="corner-cut my-4 flex justify-center background-2 p-2">
        <p class="mt-4 text-red-500">
          {{ error }}
        </p>
      </div>
    </div>
  </div>
</template>
