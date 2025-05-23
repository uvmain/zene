<script setup lang="ts">
import type { ArtistMetadata } from '../types'

defineProps({
  artist: { type: Object as PropType<ArtistMetadata>, required: true },
})

function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}
</script>

<template>
  <div>
    <img
      class="w-full rounded-md"
      :src="artist.image_url"
      @error="onImageError"
    />
    <div class="text-nowrap text-xs text-gray-300">
      {{ artist.artist }}
    </div>
  </div>
</template>
