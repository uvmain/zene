<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { currentQueue } from '~/logic/playbackQueue'
import { queueStore } from '~/logic/store'

const tracks = ref<SubsonicSong[]>([])

watch(currentQueue, async () => {
  tracks.value = currentQueue?.value ?? [] as SubsonicSong[]
  queueStore.value = tracks.value
})

onBeforeMount(() => {
  if (queueStore.value.length > 0 && (currentQueue.value === undefined || currentQueue.value.length === 0)) {
    tracks.value = queueStore.value as SubsonicSong[]
  }
  else {
    tracks.value = currentQueue?.value ?? [] as SubsonicSong[]
  }
})
</script>

<template>
  <Tracks :tracks="tracks" />
</template>
