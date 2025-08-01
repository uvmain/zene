<script setup lang="ts">
import type { TrackMetadata } from '../types'
import { useBackendFetch } from '../composables/useBackendFetch'

const props = defineProps({
  track: { type: Object as PropType<TrackMetadata>, required: true },
  currentSeconds: { type: Number, default: 0 },
})

defineEmits(['close'])

const { getLyrics } = useBackendFetch()

const plainLyricsRef = ref<string>('')
const syncedLyricsRef = ref<string>('')

interface SyncedLyricsLine {
  seconds: number
  lyrics: string
}

const computedSyncedLyrics = computed((): SyncedLyricsLine[] => {
  return syncedLyricsRef.value?.split('\n').map((line) => {
    const timeMatch = line.match(/^\[(\d+):(\d+)\.(\d+)\]/)
    if (timeMatch) {
      const minutes = Number.parseInt(timeMatch[1], 10)
      const seconds = Number.parseInt(timeMatch[2], 10)
      const milliseconds = Number.parseInt(timeMatch[3], 10)
      const startTimeSeconds = minutes * 60 + seconds + milliseconds / 1000
      const lyrics = line.replace(/^\[\d+:\d+\.\d+\]\s*/, '').trim()
      return {
        seconds: startTimeSeconds,
        lyrics,
      }
    }
    else {
      return null
    }
  }).filter(line => line && line.lyrics.length > 0) as SyncedLyricsLine[]
})

async function fetchLyrics() {
  try {
    const { plainLyrics, syncedLyrics } = await getLyrics(props.track.musicbrainz_track_id)
    plainLyricsRef.value = plainLyrics
    syncedLyricsRef.value = syncedLyrics
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
  <div v-if="(plainLyricsRef || syncedLyricsRef).length > 0" class="left-0 top-0 isolate z-200 max-h-80vh flex items-center justify-center overflow-y-scroll bg-black/80 p-4 backdrop-blur-2xl">
    <div>
      <div class="">
        <button class="ml-2 mt-2 rounded-full bg-zene-400/20 p-1 text-white hover:bg-zene-400/30" @click="$emit('close')">
          Close
        </button>
      </div>
      <div>
        <div v-if="computedSyncedLyrics" class="flex flex-col gap-2 text-center">
          <div
            v-for="(line, index) in computedSyncedLyrics" :key="index"
            :class="{
              'text-green': line.seconds <= currentSeconds && computedSyncedLyrics[index + 1]?.seconds > currentSeconds,
              'text-white': line.seconds >= currentSeconds,
            }"
          >
            {{ line.seconds <= currentSeconds && computedSyncedLyrics[index + 1]?.seconds > currentSeconds
              ? '▶️ '
              : '' }}
            <span>{{ line.lyrics }}</span>
          </div>
        </div>
        <p v-else class="flex flex-col gap-2 text-center">
          {{ plainLyricsRef || 'No lyrics available.' }}
        </p>
      </div>
    </div>
  </div>
</template>
