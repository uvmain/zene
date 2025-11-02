<script setup lang="ts">
import type { StructuredLyricLine } from '../types/subsonicLyrics'
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchLyrics } from '../composables/backendFetch'

const props = defineProps({
  currentlyPlayingTrack: { type: Object as PropType<SubsonicSong | null>, default: null },
  currentTime: { type: Number, default: 0 },
})

const showLyrics = ref(false)

const lyricsRef = ref<StructuredLyricLine[]>([])

async function getLyrics() {
  if (props.currentlyPlayingTrack?.musicBrainzId === undefined) {
    lyricsRef.value = []
    return
  }
  const lyrics = await fetchLyrics(props.currentlyPlayingTrack.musicBrainzId)
  lyricsRef.value = lyrics.line.map(line => ({
    start: line.start ? line.start / 1000 : 0, // convert milliseconds to seconds, default to 0
    value: line.value,
  }))
}

watch(() => props.currentlyPlayingTrack?.musicBrainzId, async (newTrack, oldTrack) => {
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
      class="h-10 w-10 flex cursor-pointer items-center justify-center border-none bg-white/0 text-muted font-semibold outline-none md:h-12 md:w-12 sm:h-10 sm:w-10"
      @click="showLyrics = !showLyrics"
    >
      <icon-nrk-mening class="footer-icon" />
    </button>
    <LyricsModal
      v-if="showLyrics && currentlyPlayingTrack"
      :lyrics="lyricsRef"
      :current-seconds="currentTime"
      @close="showLyrics = false"
    />
  </div>
</template>
