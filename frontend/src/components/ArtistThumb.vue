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
  <div class="flex flex-col gap-2">
    <div class="size-30">
      <img
        class="h-full w-full rounded-md object-cover"
        :src="artist.image_url"
        @error="onImageError"
      />
    </div>
    <div class="text-nowrap text-xs text-gray-300">
      {{ artist.artist }}
    </div>
  </div>
</template>
