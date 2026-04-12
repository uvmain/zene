<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { currentQueue } from '~/logic/playbackQueue'

const tracks = ref<SubsonicSong[]>([])

watch(currentQueue, async () => {
  tracks.value = currentQueue?.value ?? [] as SubsonicSong[]
})

onBeforeMount(() => {
  tracks.value = currentQueue?.value ?? [] as SubsonicSong[]
})
</script>

<template>
  <Tracks :tracks="tracks" />
</template>
