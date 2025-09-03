<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'

const { currentQueue } = usePlaybackQueue()

const tracks = computed(() => currentQueue?.value?.tracks ?? [] as SubsonicSong[])
</script>

<template>
  <div>
    <h2 class="px-2 text-lg font-semibold">
      Queue
    </h2>
    <Tracks :tracks="tracks" :show-album="true" />
  </div>
</template>
