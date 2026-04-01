<script setup lang="ts">
import type { StructuredLyricLine } from '~/types/subsonicLyrics'
import { fetchLyrics } from '~/logic/backendFetch'
import { currentlyPlayingItem, currentTime } from '~/logic/playbackQueue'

const showLyrics = ref(false)

const lyricsRef = ref<StructuredLyricLine[]>([])

async function getLyrics() {
  if (currentlyPlayingItem.value.track?.musicBrainzId === undefined) {
    lyricsRef.value = []
    return
  }
  const lyrics = await fetchLyrics(currentlyPlayingItem.value.track?.musicBrainzId)
  if (!lyrics) {
    lyricsRef.value = []
    return
  }
  lyricsRef.value = lyrics.line.map(line => ({
    start: line.start ? line.start / 1000 : 0, // convert milliseconds to seconds, default to 0
    value: line.value,
  }))
}

watch(() => currentlyPlayingItem.value.track?.musicBrainzId, async (newTrack, oldTrack) => {
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
    <abbr :title="lyricsRef.length > 0 ? 'Show lyrics' : 'No Lyrics Available'">
      <button
        id="lyrics"
        class="text-muted font-semibold outline-none border-none bg-white/0 flex items-center justify-center"
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
    </abbr>
    <LyricsModal
      v-if="showLyrics && currentlyPlayingItem.track"
      :lyrics="lyricsRef"
      :current-seconds="currentTime"
      @close="showLyrics = false"
    />
  </div>
</template>
