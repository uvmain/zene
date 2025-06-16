<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { useBackendFetch } from '../composables/useBackendFetch'
import { useRouteTracks } from '../composables/useRouteTracks'

const { routeTracks, clearRouteTracks } = useRouteTracks()
const { backendFetchRequest } = useBackendFetch()

const tracks = ref<TrackMetadataWithImageUrl[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    const response = await backendFetchRequest('tracks?random=true&limit=100')
    const json = await response.json() as TrackMetadataWithImageUrl[]
    tracks.value = json
    routeTracks.value = json
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
