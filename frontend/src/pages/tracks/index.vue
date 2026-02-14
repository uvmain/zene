<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchRandomTracks } from '~/logic/backendFetch'
import { generateSeed } from '~/logic/common'
import { clearRouteTracks, routeTracks } from '~/logic/routeTracks'
import { randomTracksSeed } from '~/logic/store'

const tracks = ref<SubsonicSong[]>()
const canLoadMore = ref<boolean>(true)
const limit = ref<number>(100)
const offset = ref<number>(0)

async function getData() {
  const randomTracks = await fetchRandomTracks(limit.value, offset.value, randomTracksSeed.value)
  tracks.value = tracks.value?.concat(randomTracks) ?? randomTracks
  routeTracks.value = tracks.value
  offset.value += randomTracks.length

  if (randomTracks.length < limit.value) {
    canLoadMore.value = false
  }
}

onMounted(async () => {
  randomTracksSeed.value = generateSeed()
  await getData()
})

onUnmounted(() => clearRouteTracks())
</script>

<template>
  <Tracks v-if="tracks" :tracks="tracks" :show-album="true" :observer-enabled="canLoadMore" @observer-visible="getData" />
</template>
