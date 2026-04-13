<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchRandomTracks } from '~/logic/backendFetch'
import { generateSeed } from '~/logic/common'
import { clearRouteTracks, routeTracks } from '~/logic/routeTracks'
import { randomTracksSeed } from '~/logic/store'

const tracks = ref<SubsonicSong[]>()

async function getData() {
  const randomTracks = await fetchRandomTracks({ seed: randomTracksSeed.value })
  tracks.value = tracks.value?.concat(randomTracks) ?? randomTracks
  routeTracks.value = tracks.value
}

onMounted(async () => {
  randomTracksSeed.value = generateSeed()
  await getData()
})

onUnmounted(() => clearRouteTracks())
</script>

<template>
  <Tracks v-if="tracks" :tracks="tracks" />
</template>
