<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { currentQueue } from '~/logic/playbackQueue'

const tracks = ref<SubsonicSong[]>()

watch(currentQueue, async () => {
  tracks.value = currentQueue.value
}, { immediate: true })

</script>

<template>
  <Tracks v-if="tracks" :tracks="tracks" />
</template>
