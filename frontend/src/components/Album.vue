<script setup lang="ts">
import type { AlbumMetadata } from '../types'

const props = defineProps({
  album: { type: Object as PropType<AlbumMetadata>, required: true },
  size: { type: String, default: 'md' },
})

const artistAndDate = computed(() => {
  return props.album.release_date !== 'Invalid Date' ? `${props.album.artist} • ${props.album.release_date}` : props.album.artist
})

function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}
</script>

<template>
  <div v-if="props.size === 'lg'">
    <img class="w-full rounded-md" :src="album.image_url" alt="Album Cover" @error="onImageError" />
    <div class="text-nowrap text-sm">
      {{ album.album }}
    </div>
    <div class="text-nowrap text-xs text-gray-300">
      {{ album.album_artist }}
    </div>
  </div>
  <div v-else-if="props.size === 'xl'" class="h-full flex items-center gap-6 from-zene-600/90 via-zene-600/80 bg-gradient-to-r p-10">
    <img :src="album.image_url" class="size-50 rounded-lg object-cover" @error="onImageError">
    <div class="flex flex-col gap-5">
      <div class="text-4xl text-white font-bold">
        {{ album.album }}
      </div>
      <div class="text-xl text-white">
        {{ artistAndDate }}
      </div>
      <div v-if="album.genres.length > 0" class="flex flex-wrap gap-2">
        <GenreBottle v-for="genre in album.genres" :key="genre" :genre />
      </div>
      <button class="w-30 border-1 border-white rounded-full border-solid bg-zene-600/70 px-4 py-2 text-xl text-white outline-none hover:bg-zene-200">
        Play
      </button>
    </div>
  </div>
</template>
