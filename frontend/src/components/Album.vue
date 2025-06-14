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
  <div v-if="props.size === 'lg'" class="h-40 w-30">
    <img class="size-30 cursor-pointer rounded-lg rounded-md object-cover" :src="album.image_url" alt="Album Cover" @error="onImageError" @click="navigateAlbum()" />
    <div class="w-30 truncate text-nowrap text-sm">
      {{ album.album }}
    </div>
    <div class="w-30 cursor-pointer truncate text-nowrap text-xs text-gray-300" @click="navigateArtist()">
      {{ artistAndDate }}
    </div>
  </div>
  <div v-else-if="props.size === 'xl'" class="h-full flex items-center gap-6 from-zene-600/90 via-zene-600/80 bg-gradient-to-r p-10">
    <img :src="album.image_url" class="size-50 cursor-pointer rounded-lg object-cover" @error="onImageError" @click="navigateAlbum()">
    <div class="flex flex-col gap-5">
      <div class="cursor-pointer text-4xl text-white font-bold" @click="navigateAlbum()">
        {{ album.album }}
      </div>
      <div class="cursor-pointer text-xl text-white" @click="navigateArtist()">
        {{ artistAndDate }}
      </div>
      <div v-if="album.genres.length > 0" class="flex flex-wrap gap-2">
        <GenreBottle v-for="genre in album.genres" :key="genre" :genre />
      </div>
      <PlayButton :album="album" />
    </div>
  </div>
</template>
