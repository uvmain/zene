<script setup lang="ts">
import type { StructuredLyricLine } from '../types/subsonicLyrics'
import type { SubsonicSong } from '../types/subsonicSong'
import { fetchLyrics } from '../composables/backendFetch'

const props = defineProps({
  track: { type: Object as PropType<SubsonicSong>, required: true },
  currentSeconds: { type: Number, default: 0 },
})

defineEmits(['close'])

const lyricsRef = ref<StructuredLyricLine[]>([])

async function getLyrics() {
  const lyrics = await fetchLyrics(props.track.musicBrainzId)
  lyricsRef.value = lyrics.line.map(line => ({
    start: line.start ? line.start / 1000 : 0, // convert milliseconds to seconds, default to 0
    value: line.value,
  }))
}

watch(props.track, async (newTrack, oldTrack) => {
  if (newTrack !== oldTrack) {
    await getLyrics()
  }
})

onMounted(async () => {
  await getLyrics()
})
</script>

<template>
  <div v-if="lyricsRef.length > 0" class="left-0 top-0 isolate z-200 max-h-80vh flex items-center justify-center overflow-y-scroll bg-black/80 p-4 backdrop-blur-2xl">
    <div>
      <div class="">
        <button class="bg-zgray-400/20 hover:bg-zgray-400/30 ml-2 mt-2 p-1" @click="$emit('close')">
          Close
        </button>
      </div>
      <div>
        <div v-if="lyricsRef" class="flex flex-col gap-2 text-center">
          <div
            v-for="(line, index) in lyricsRef" :key="index"
            :class="{
              'text-green': line.start && line.start <= currentSeconds && (lyricsRef[index + 1]?.start ?? 0) > currentSeconds,
              '': (line.start ?? 0) >= currentSeconds,
            }"
          >
            {{ (line.start ?? 0) <= currentSeconds && (lyricsRef[index + 1]?.start ?? 0) > currentSeconds
              ? '▶️ '
              : '' }}
            <span>{{ line.value }}</span>
          </div>
        </div>
        <p v-else class="flex flex-col gap-2 text-center">
          {{ lyricsRef || 'No lyrics available.' }}
        </p>
      </div>
    </div>
  </div>
</template>
