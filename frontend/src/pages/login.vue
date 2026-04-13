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
const passwordRef = useTemplateRef('passwordRef')

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
  <div class="my-auto p-4 background-1 flex flex-col gap-6 max-w-screen items-center justify-center lg:(flex-row gap-12)">
    <img
      class="opacity-90 size-full max-w-400px"
      src="/minidisk.svg"
      alt="Logo"
      width="200"
      height="200"
    />
    <div class="text-muted flex flex-col gap-6 items-center justify-center lg:justify-start">
      <div class="text-xl font-bold">
        Login to Zene
      </div>
      <form class="gap-2 grid w-400px" @submit.prevent="login">
        <label for="username">
          Username
        </label>
        <input
          id="username"
          v-model="username"
          type="text"
          class="input"
          autocomplete="username"
          required
          @input="error = null"
          @keydown.enter.prevent="passwordRef?.focus()"
        />
        <label for="password">
          Password
        </label>
        <input
          id="password"
          ref="passwordRef"
          v-model="password"
          type="password"
          class="input"
          autocomplete="current-password"
          required
          @input="error = null"
          @keydown.enter.prevent="login()"
        />
      </form>
      <ZButton :disabled="signInDisabled" @click="login()">
        <div class="text-xl lg:text-base">
          {{ loading ? 'Signing in…' : 'Sign in' }}
        </div>
      </ZButton>
      <div v-if="error" class="my-4 p-2 corner-cut background-2 flex justify-center">
        <p class="text-red-500 mt-4">
          {{ error }}
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.input {
  @apply text-muted p-2 border-1 border-background-300 corner-cut border-solid focus:outline-none focus:ring-2 focus:ring-primary-500;
  @apply bg-background-100 dark:bg-background-800;
}
</style>
