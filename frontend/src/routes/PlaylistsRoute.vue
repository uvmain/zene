<script setup lang="ts">
import type { TokenResponse } from '~/types/auth'
import { useDebug } from '~/composables/useDebug'

const { debugLog } = useDebug()

const audioUrl = ref('https://static.ianbaron.com/dc2b0ca2-ff9c-41ec-9672-58a96f5e58bd-160.aac')
const audioEl = ref<HTMLAudioElement | null>(null)
const error = ref<string | null>(null)
const session = ref<cast.framework.CastSession | null>(null)
const temporaryToken = ref<TokenResponse | null>(null)

async function getMimeType(): Promise<string> {
  await fetch(audioUrl.value)
    .then((response) => {
      const contentType = response.headers.get('content-type') ?? ''
      debugLog(contentType)
      return contentType
    })
    .catch((err) => {
      debugLog(`fetch failed: ${err}`)
    })
  return ''
}

async function castAudio() {
  const context = cast.framework.CastContext.getInstance()
  session.value = context.getCurrentSession()

  if (!session.value) {
    error.value = 'No cast session available'
    return
  }

  if (!audioUrl.value) {
    error.value = 'Please enter a valid audio URL.'
    return
  }

  const contentType = await (getMimeType())

  let requestUrl: string
  if (audioUrl.value.includes('?')) {
    requestUrl = `${audioUrl.value}&token=${temporaryToken.value?.token}`
  }
  else {
    requestUrl = `${audioUrl.value}?token=${temporaryToken.value?.token}`
  }
  const mediaInfo = new chrome.cast.media.MediaInfo(requestUrl, contentType)
  const request = new chrome.cast.media.LoadRequest(mediaInfo)

  session.value
    .loadMedia(request)
    .then(() => debugLog('Media loaded to cast device'))
    .catch(err => console.error('Error loading media:', err))
}

function initializeCast() {
  const context = cast.framework.CastContext.getInstance()
  context.setOptions({
    receiverApplicationId: chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
    autoJoinPolicy: chrome.cast.AutoJoinPolicy.ORIGIN_SCOPED,
  })
  debugLog('CastContext initialized')
}

onMounted(async () => {
  debugLog('Waiting for Cast SDK...')

  // Hook for when the SDK becomes available (in case it's not already)
  window.__onGCastApiAvailable = (isAvailable: boolean) => {
    debugLog(`Cast API available (async):', ${isAvailable}`)
    if (isAvailable)
      initializeCast()
    else console.warn('Cast API not available')
  }

  // if the SDK already loaded and called __onGCastApiAvailable BEFORE this script ran
  // we have to check manually and initialize right now
  if ((window.cast && window.cast.isAvailable) || (window.chrome?.cast && window.chrome.cast.isAvailable)) {
    debugLog('Cast API already available (sync)')
    initializeCast()
  }
})
</script>

<template>
  <div class="p-4 space-y-4">
    <input
      v-model="audioUrl"
      name="audio url"
      type="text"
      placeholder="Enter audio URL"
      class="w-full border rounded px-2 py-1"
    />

    <div v-if="error">
      {{ error }}
    </div>

    <audio ref="audioEl" controls :src="audioUrl" class="w-full" />

    <!-- Cast button -->
    <div class="flex items-center space-x-2">
      <div class="inline-block size-24px">
        <google-cast-launcher style="inline-block" />
      </div>
      <button class="rounded bg-blue-600 px-4 py-2 text-white" @click="castAudio">
        Cast to Device
      </button>
    </div>
  </div>
</template>

<style scoped>
google-cast-launcher {
  @apply size-24px inline-block cursor-pointer bg-white;
}
</style>
