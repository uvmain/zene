<script setup lang="ts">
import md5 from 'md5'
import QRCode from 'qrcode'
import { createNewApiKeyWithTokenAndSalt, fetchApiKeysWithTokenAndSalt } from '~/logic/backendFetch'
import { completeQRLoginChallenge, createQRLoginChallenge, getQRLoginChallengeStatus } from '~/logic/qrLogin'
import { apiKey } from '~/stores/main'

const router = useRouter()
const route = useRoute()

const username = ref('')
const token = ref('')
const salt = ref('')
const password = ref('')
const loading = ref(false)
const error = ref<string | null>(null)
const passwordRef = useTemplateRef('passwordRef')
const qrCanvas = useTemplateRef('qrCanvas')
const qrLoading = ref(false)
const qrError = ref('')
const qrExpiresAt = ref<Date | null>(null)
const qrSecondsRemaining = ref(0)
const computerSecret = ref('')
let qrPollTimer: ReturnType<typeof setTimeout> | undefined
let qrCountdownTimer: ReturnType<typeof setInterval> | undefined

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
      await redirectAfterLogin()
    }
    else {
      const existingApiKey = data?.apiKeys.apiKey[0]?.api_key
      apiKey.value = existingApiKey
      await redirectAfterLogin()
    }
  }
  catch (e: any) {
    error.value = e?.message || 'Login failed'
  }
  finally {
    loading.value = false
  }
}

async function redirectAfterLogin() {
  const redirect = typeof route.query.redirect === 'string'
    ? route.query.redirect
    : '/'
  await router.push(redirect.startsWith('/') && !redirect.startsWith('//')
    ? redirect
    : '/')
}

function stopQRLogin() {
  if (qrPollTimer) {
    clearTimeout(qrPollTimer)
  }
  if (qrCountdownTimer) {
    clearInterval(qrCountdownTimer)
  }
  qrPollTimer = undefined
  qrCountdownTimer = undefined
  computerSecret.value = ''
}

function updateQRCountdown() {
  if (!qrExpiresAt.value) {
    return
  }
  qrSecondsRemaining.value = Math.max(0, Math.ceil((qrExpiresAt.value.getTime() - Date.now()) / 1000))
  if (qrSecondsRemaining.value === 0) {
    qrError.value = 'This QR code has expired.'
    stopQRLogin()
  }
}

async function pollQRLogin() {
  if (!computerSecret.value) {
    return
  }
  try {
    const challenge = await getQRLoginChallengeStatus(computerSecret.value)
    if (challenge.state === 'approved') {
      const completed = await completeQRLoginChallenge(computerSecret.value)
      apiKey.value = completed.apiKey
      stopQRLogin()
      await redirectAfterLogin()
      return
    }
    if (challenge.state === 'denied' || challenge.state === 'expired' || challenge.state === 'consumed') {
      qrError.value = challenge.state === 'denied'
        ? 'This QR login request was denied.'
        : 'This QR code has expired.'
      stopQRLogin()
      return
    }
    qrPollTimer = setTimeout(pollQRLogin, 1500)
  }
  catch (cause) {
    qrError.value = cause instanceof Error
      ? cause.message
      : 'Could not check QR login status.'
    stopQRLogin()
  }
}

async function startQRLogin() {
  stopQRLogin()
  qrLoading.value = true
  qrError.value = ''
  try {
    const challenge = await createQRLoginChallenge()
    computerSecret.value = challenge.computerSecret
    qrExpiresAt.value = new Date(challenge.expiresAt)
    await nextTick()
    const qrURL = new URL('/qr-login', window.location.origin)
    qrURL.hash = new URLSearchParams({ challenge: challenge.phoneSecret }).toString()
    if (!qrCanvas.value) {
      throw new Error('QR code display is unavailable.')
    }
    await QRCode.toCanvas(qrCanvas.value, qrURL.toString(), { width: 240, margin: 1 })
    updateQRCountdown()
    qrCountdownTimer = setInterval(updateQRCountdown, 1000)
    qrPollTimer = setTimeout(pollQRLogin, 1500)
  }
  catch (cause) {
    qrError.value = cause instanceof Error
      ? cause.message
      : 'Could not create QR login.'
    stopQRLogin()
  }
  finally {
    qrLoading.value = false
  }
}

onUnmounted(stopQRLogin)
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
      <div class="pt-4 border-background-300 border-t-1 border-solid flex flex-col gap-3 items-center">
        <ZButton :disabled="qrLoading" @click="startQRLogin">
          {{ qrLoading ? 'Creating QR code…' : 'Log in with QR code' }}
        </ZButton>
        <template v-if="computerSecret">
          <canvas ref="qrCanvas" class="p-2 bg-white" aria-label="QR code for login" />
          <p>Scan this code on a device where you are already logged in.</p>
          <p>Expires in {{ qrSecondsRemaining }} seconds.</p>
        </template>
        <p v-if="qrError" class="text-red-500">
          {{ qrError }}
        </p>
      </div>
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
  @apply text-muted p-2 border-1 border-background-300 corner-cut border-solid focus:outline-none focus:ring-2 focus:ring-main-500;
  @apply bg-background-100 dark:bg-background-800;
}
</style>
