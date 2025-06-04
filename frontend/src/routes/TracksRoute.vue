<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const tracks = ref<TrackMetadataWithImageUrl[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    const response = await backendFetchRequest('tracks?random=true&limit=100')
    const json = await response.json()
    tracks.value = json as TrackMetadataWithImageUrl[]
  }
  catch (err) {
    error.value = 'Failed to fetch tracks.'
    console.error(err)
  }
  finally {
    loading.value = false
  }
})
</script>

<template>
  <Tracks :tracks="tracks" :show-album="true" />
</template>
