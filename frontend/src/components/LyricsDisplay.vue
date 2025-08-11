<script setup lang="ts">
import type { TrackMetadata } from '../types'
import type { StructuredLyricLine } from '../types/subsonicLyrics'
import { useBackendFetch } from '../composables/useBackendFetch'

const props = defineProps({
  track: { type: Object as PropType<TrackMetadata>, required: true },
  currentSeconds: { type: Number, default: 0 },
})

defineEmits(['close'])

const { getLyrics } = useBackendFetch()

const lyricsRef = ref<StructuredLyricLine[]>([])

async function fetchLyrics() {
  try {
    const response = await getLyrics(props.track.musicbrainz_track_id)
    if (response['subsonic-response'].lyricsList.structuredLyrics[0].synced) {
      lyricsRef.value = response['subsonic-response'].lyricsList.structuredLyrics[0].line.map((line: StructuredLyricLine) => ({
        start: line.start ? line.start / 1000 : undefined, // convert milliseconds to seconds
        value: line.value,
      }))
    }
    else {
      lyricsRef.value = response['subsonic-response'].lyricsList.structuredLyrics[0].line.map((line: StructuredLyricLine) => ({
        value: line.value,
      }))
    }
  }
  catch (error) {
    console.error('Failed to fetch lyrics:', error)
  }
}

watch(props.track, async (newTrack, oldTrack) => {
  if (newTrack !== oldTrack) {
    await fetchLyrics()
  }
})

onMounted(async () => {
  await fetchLyrics()
})
</script>

<template>
  <div v-if="lyricsRef.length > 0" class="left-0 top-0 isolate z-200 max-h-80vh flex items-center justify-center overflow-y-scroll bg-black/80 p-4 backdrop-blur-2xl">
    <div>
      <div class="">
        <button class="ml-2 mt-2 rounded-full bg-zene-400/20 p-1 text-white hover:bg-zene-400/30" @click="$emit('close')">
          Close
        </button>
      </div>
      <div>
        <div v-if="lyricsRef" class="flex flex-col gap-2 text-center">
          <div
            v-for="(line, index) in lyricsRef" :key="index"
            :class="{
              'text-green': line.start && line.start <= currentSeconds && (lyricsRef[index + 1]?.start ?? 0) > currentSeconds,
              'text-white': (line.start ?? 0) >= currentSeconds,
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
