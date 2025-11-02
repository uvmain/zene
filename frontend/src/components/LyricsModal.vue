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
  <teleport v-if="lyricsRef.length > 0" to="body">
    <div class="fixed inset-0 z-50 flex justify-center overflow-y-scroll bg-white/20 p-4 backdrop-blur-lg dark:bg-black/20">
      <div>
        <ZButton class="absolute right-4 top-4" @click="$emit('close')">
          Close
        </ZButton>
        <div v-if="lyricsRef" class="flex flex-col gap-2 text-center text-muted">
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
  </teleport>
</template>
