<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { trackWithImageUrl } from '../composables/logic'
import { useBackendFetch } from '../composables/useBackendFetch'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'
import { useRouteTracks } from '../composables/useRouteTracks'

const { routeTracks, clearRouteTracks } = useRouteTracks()
const { backendFetchRequest } = useBackendFetch()
const { getRandomSeed } = usePlaybackQueue()

const tracks = ref<TrackMetadataWithImageUrl[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    const randomSeed = getRandomSeed()
    const response = await backendFetchRequest(`tracks?random=${randomSeed}&limit=100`)
    const json = await response.json() as TrackMetadataWithImageUrl[]

    const tracksWithImages = json.map((track) => {
      const newTrack = trackWithImageUrl(track)
      return newTrack
    })
    tracks.value = tracksWithImages
    routeTracks.value = tracksWithImages
  }
  catch (err) {
    error.value = 'Failed to fetch tracks.'
    console.error(err)
  }
  finally {
    loading.value = false
  }
})

onUnmounted(() => clearRouteTracks())
</script>

<template>
  <Tracks :tracks="tracks" :show-album="true" />
</template>
