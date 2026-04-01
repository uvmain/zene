<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { currentQueue } from '~/logic/playbackQueue'

const tracks = ref<SubsonicSong[]>([])

watch(currentQueue, async () => {
  tracks.value = currentQueue?.value ?? [] as SubsonicSong[]
})

onMounted(() => {
  tracks.value = currentQueue?.value ?? [] as SubsonicSong[]
})
</script>

<template>
  <div>
    <Tracks v-if="tracks.length" :tracks="tracks" :show-album="true" />
    <h2 v-else class="text-lg font-semibold px-2">
      No tracks in queue..
    </h2>
  </div>
</template>
