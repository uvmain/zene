<script setup lang="ts">
import type { AlbumMetadata, ArtistMetadata, TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'

defineProps({
  artist: { type: Object as PropType<ArtistMetadata>, required: false },
  album: { type: Object as PropType<AlbumMetadata>, required: false },
  track: { type: Object as PropType<TrackMetadata | TrackMetadataWithImageUrl>, required: false },
})

const { play } = usePlaybackQueue()
</script>

<template>
  <button class="w-30 border-1 border-white rounded-full border-solid bg-zene-600/70 px-4 py-2 text-xl text-white outline-none transition-all duration-200 hover:from-zene-400 hover:to-zene-200 hover:bg-gradient-to-b hover:shadow-sm hover:shadow-zene-400" @click="play(artist, album, track)">
    Play
  </button>
</template>
