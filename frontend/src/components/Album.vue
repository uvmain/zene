<script setup lang="ts">
import type { AlbumMetadata } from '../types'
import { useSearch } from '../composables/useSearch'

const props = defineProps({
  album: { type: Object as PropType<AlbumMetadata>, required: true },
  size: { type: String, default: 'md' },
})

const router = useRouter()
const { closeSearch } = useSearch()

const artistAndDate = computed(() => {
  return props.album.release_date !== 'Invalid Date' ? `${props.album.artist} • ${props.album.release_date}` : props.album.artist
})

function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}

function navigateAlbum() {
  closeSearch()
  router.push(`/albums/${props.album.musicbrainz_album_id}`)
}

function navigateArtist() {
  closeSearch()
  router.push(`/artists/${props.album.musicbrainz_artist_id}`)
}
</script>

<template>
  <div v-if="props.size === 'lg'" class="h-32 w-24 md:h-40 md:w-30">
    <img class="h-24 w-24 cursor-pointer rounded-lg object-cover md:size-30" :src="album.image_url" alt="Album Cover" @error="onImageError" @click="navigateAlbum()" />
    <div class="w-24 truncate text-nowrap text-xs md:w-30 md:text-sm">
      {{ album.album }}
    </div>
    <div class="w-24 cursor-pointer truncate text-nowrap text-xs text-gray-300 md:w-30" @click="navigateArtist()">
      {{ artistAndDate }}
    </div>
  </div>
  <div v-else-if="props.size === 'xl'" class="h-full flex flex-col items-center gap-4 from-zene-600/90 via-zene-600/80 bg-gradient-to-r p-4 md:flex-row md:gap-6 md:p-10">
    <img :src="album.image_url" class="h-32 w-32 cursor-pointer rounded-lg object-cover md:size-50" @error="onImageError" @click="navigateAlbum()">
    <div class="flex flex-col gap-3 text-center md:gap-5 md:text-left">
      <div class="cursor-pointer text-xl text-white font-bold md:text-4xl" @click="navigateAlbum()">
        {{ album.album }}
      </div>
      <div class="cursor-pointer text-lg text-white md:text-xl" @click="navigateArtist()">
        {{ artistAndDate }}
      </div>
      <div v-if="album.genres.length > 0" class="flex flex-wrap justify-center gap-2 md:justify-start">
        <GenreBottle v-for="genre in album.genres" :key="genre" :genre />
      </div>
      <div class="flex justify-center md:justify-start">
        <PlayButton :album="album" />
      </div>
    </div>
  </div>
</template>
