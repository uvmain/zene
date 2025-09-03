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
const error = ref('')

const apiKey = useLocalStorage('apiKey', '')

async function login() {
  error.value = ''
  loading.value = true
  try {
    salt.value = Math.random().toString(36).slice(2, 10)
    token.value = md5(password.value + salt.value)

    const data = await fetchApiKeysWithTokenAndSalt(username.value, token.value, salt.value)

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
  <div class="login">
    <h1>Login to Zene</h1>
    <form @submit.prevent="login">
      <label>
        Username
        <input v-model="username" autocomplete="username" />
      </label>
      <label>
        Password
        <input v-model="password" type="password" autocomplete="current-password" />
      </label>
      <button :disabled="loading">
        {{ loading ? 'Signing in…' : 'Sign in' }}
      </button>
      <p v-if="error" class="error">
        {{ error }}
      </p>
    </form>
  </div>
</template>

<style scoped>
.login { max-width: 420px; margin: 6rem auto; padding: 2rem; border: 1px solid #e1e1e1; border-radius: 12px; background-color: #7e0a361c; }
form { display: grid; gap: 1rem; }
label { display: grid; gap: 0.25rem; font-weight: 600; }
input { padding: 0.5rem 0.75rem; border: 1px solid #ccc; border-radius: 8px; }
button { padding: 0.6rem 0.9rem; border-radius: 8px; border: none; background:#42b883; color: white; font-weight: 700; cursor: pointer; }
.error { color: #c62828; }
</style>
