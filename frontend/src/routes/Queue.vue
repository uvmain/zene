<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'

const { currentQueue } = usePlaybackQueue()

const tracks = computed(() => currentQueue?.value?.tracks ?? [] as TrackMetadataWithImageUrl[])
</script>

<template>
  <div>
    <h2 class="px-2 text-lg font-semibold">
      Queue
    </h2>
    <Tracks :tracks="tracks" :show-album="true" />
  </div>
</template>
