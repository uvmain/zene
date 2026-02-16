<script setup lang="ts">
import type { StructuredLyricLine } from '~/types/subsonicLyrics'
import { fetchLyrics } from '~/logic/backendFetch'
import { currentlyPlayingTrack, currentTime } from '~/logic/playbackQueue'

const showLyrics = ref(false)

const lyricsRef = ref<StructuredLyricLine[]>([])

async function getLyrics() {
  if (currentlyPlayingTrack.value?.musicBrainzId === undefined) {
    lyricsRef.value = []
    return
  }
  const lyrics = await fetchLyrics(currentlyPlayingTrack.value?.musicBrainzId)
  if (!lyrics) {
    lyricsRef.value = []
    return
  }
  lyricsRef.value = lyrics.line.map(line => ({
    start: line.start ? line.start / 1000 : 0, // convert milliseconds to seconds, default to 0
    value: line.value,
  }))
}

watch(() => currentlyPlayingTrack.value?.musicBrainzId, async (newTrack, oldTrack) => {
  if (newTrack !== oldTrack) {
    await getLyrics()
  }
})

onMounted(async () => {
  await getLyrics()
})
</script>

<template>
  <div>
    <button
      id="lyrics"
      class="h-10 w-10 flex items-center justify-center border-none bg-white/0 text-muted font-semibold outline-none lg:h-12 lg:w-12 sm:h-10 sm:w-10"
      :class="{
        'cursor-not-allowed': lyricsRef.length === 0,
        'cursor-pointer': lyricsRef.length > 0,
      }"
      :disabled="lyricsRef.length === 0"
      @click="showLyrics = !showLyrics"
    >
      <icon-nrk-mening
        :class="{
          'footer-icon-disabled': lyricsRef.length === 0,
          'footer-icon': lyricsRef.length > 0,
        }"
      />
    </button>
    <LyricsModal
      v-if="showLyrics && currentlyPlayingTrack"
      :lyrics="lyricsRef"
      :current-seconds="currentTime"
      @close="showLyrics = false"
    />
  </div>
</template>
