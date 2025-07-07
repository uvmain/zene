<script setup lang="ts">
const audioUrl = ref('https://zene.ianbaron.com/api/tracks/c4750da7-7e37-4347-aad8-0beb2ac84ad1/stream?quality=160')
const audioEl = ref<HTMLAudioElement | null>(null)
const error = ref<string | null>(null)
const session = ref<cast.framework.CastSession | null>(null)

const DEFAULT_APP_ID = chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID

function castAudio() {
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

  const mediaInfo = new chrome.cast.media.MediaInfo(audioUrl.value, 'audio/aac')
  const request = new chrome.cast.media.LoadRequest(mediaInfo)

  session.value
    .loadMedia(request)
    .then(() => console.log('Media loaded to cast device'))
    .catch(err => console.error('Error loading media:', err))
}

function initializeCast() {
  const context = cast.framework.CastContext.getInstance()
  context.setOptions({
    receiverApplicationId: DEFAULT_APP_ID,
    autoJoinPolicy: chrome.cast.AutoJoinPolicy.ORIGIN_SCOPED,
  })
  console.log('CastContext initialized')
}

onMounted(() => {
  console.log('Waiting for Cast SDK...')

  // Hook for when the SDK becomes available (in case it's not already)
  window.__onGCastApiAvailable = (isAvailable: boolean) => {
    console.log('Cast API available (async):', isAvailable)
    if (isAvailable)
      initializeCast()
    else console.warn('Cast API not available')
  }

  // 🔥 If the SDK already loaded and called __onGCastApiAvailable BEFORE this script ran
  // we have to check manually and initialize right now
  if ((window.cast && window.cast.isAvailable) || (window.chrome?.cast && window.chrome.cast.isAvailable)) {
    console.log('Cast API already available (sync)')
    initializeCast()
  }
})
</script>

<template>
  <div class="p-4 space-y-4">
    <input
      v-model="audioUrl"
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
        <google-cast-launcher />
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
