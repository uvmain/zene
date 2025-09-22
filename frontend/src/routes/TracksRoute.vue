<script setup lang="ts">
import { fetchRandomTracks } from '~/composables/backendFetch'
import { useRouteTracks } from '~/composables/useRouteTracks'

const { routeTracks, clearRouteTracks } = useRouteTracks()

const error = ref<string | null>(null)

async function loadTracks() {
  routeTracks.value = await fetchRandomTracks(100)
}

onMounted(async () => {
  try {
    await loadTracks()
  }
  catch (err) {
    error.value = 'Failed to fetch tracks.'
    console.error(err)
  }
})

onUnmounted(() => clearRouteTracks())
</script>

<template>
  <Tracks :tracks="routeTracks" :show-album="true" />
</template>
