<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { currentQueue } from '~/logic/playbackQueue'
import { getStoredKV, setStoredKV } from '~/stores/keyValueIdbStore'

const tracks = ref<SubsonicSong[]>([])

watch(currentQueue, async () => {
  tracks.value = currentQueue.value
  await setStoredKV('queue', JSON.stringify(tracks.value))
})

onBeforeMount(async () => {
  const storedQueue = await getStoredKV('queue')
  if (storedQueue) {
    tracks.value = JSON.parse(storedQueue) as SubsonicSong[]
  }
  else {
    tracks.value = currentQueue.value
  }
})
</script>

<template>
  <Tracks :tracks="tracks" />
</template>
