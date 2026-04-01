<script setup lang="ts">
import md5 from 'md5'
import { createNewApiKeyWithTokenAndSalt, fetchApiKeysWithTokenAndSalt } from '~/logic/backendFetch'
import { apiKey } from '~/logic/store'

const router = useRouter()

const username = ref('')
const token = ref('')
const salt = ref('')
const password = ref('')
const loading = ref(false)
const error = ref<string | null>(null)

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
      const newApiKey = await createNewApiKeyWithTokenAndSalt(username.value, token.value, salt.value)
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
  <div class="my-auto p-4 background-1 flex flex-col gap-6 items-center justify-center lg:(flex-row gap-12)">
    <img
      class="opacity-90 size-full max-w-400px"
      src="/minidisk.svg"
      alt="Logo"
      width="200"
      height="200"
    />
    <div class="text-muted">
      <h1>Login to Zene</h1>
      <form class="gap-2 grid w-400px" @submit.prevent="login">
        <label for="username">
          Username
        </label>
        <input
          id="username"
          v-model="username"
          type="text"
          class="font-semibold py-2 pl-10 border-1 border-primary2 rounded background-2 focus:outline-none focus:border-primary2 dark:border-opacity-60 focus:border-solid focus:shadow-primary2 hover:shadow-lg"
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
          class="font-semibold py-2 pl-10 border-1 border-primary2 rounded background-2 opacity-100 focus:outline-none focus:border-primary2 dark:border-opacity-60 focus:border-solid focus:shadow-primary2 hover:shadow-lg"
          autocomplete="current-password"
          required
          @input="error = null"
        />
        <ZButton class="mx-auto mt-4" :disabled="signInDisabled">
          {{ loading ? 'Signing in…' : 'Sign in' }}
        </ZButton>
      </form>
      <div v-if="error" class="my-4 p-2 corner-cut background-2 flex justify-center">
        <p class="text-red-500 mt-4">
          {{ error }}
        </p>
      </div>
    </div>
  </div>
</template>
