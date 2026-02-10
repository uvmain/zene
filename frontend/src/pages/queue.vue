<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { currentQueue } from '~/logic/playbackQueue'

const tracks = ref<SubsonicSong[]>([])

watch(currentQueue, async () => {
  tracks.value = currentQueue?.value?.tracks ?? [] as SubsonicSong[]
})

onMounted(() => {
  tracks.value = currentQueue?.value?.tracks ?? [] as SubsonicSong[]
})
</script>

<template>
  <div>
    <h2 v-if="tracks.length" class="px-2 text-lg font-semibold">
      Queue
    </h2>
    <h2 v-else class="px-2 text-lg font-semibold">
      No tracks in queue..
    </h2>
    <Tracks v-if="tracks.length" :tracks="tracks" :show-album="true" />
  </div>
</template>
