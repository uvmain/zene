<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'

const { currentPlaylist } = usePlaybackQueue()
</script>

<template>
  <Tracks :tracks="currentPlaylist?.tracks ?? [] as TrackMetadataWithImageUrl[]" :show-album="true" />
</template>
