<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'

const { currentQueue } = usePlaybackQueue()

const tracks = computed(() => currentQueue?.value?.tracks ?? [] as TrackMetadataWithImageUrl[])
</script>

<template>
  <Tracks :tracks="tracks" :show-album="true" />
</template>
