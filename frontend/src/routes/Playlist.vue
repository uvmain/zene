<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'

const { currentPlaylist } = usePlaybackQueue()

const tracks = computed(() => currentPlaylist?.value?.tracks ?? [] as TrackMetadataWithImageUrl[])
</script>

<template>
  <Tracks :tracks="tracks" :show-album="true" />
</template>
