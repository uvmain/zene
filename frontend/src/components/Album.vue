<script setup lang="ts">
import type { AlbumMetadata } from '../types'

const props = defineProps({
  album: { type: Object as PropType<AlbumMetadata>, required: true },
  size: { type: String, default: 'md' },
})

const router = useRouter()

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
    <img class="size-25 cursor-pointer rounded-lg rounded-md object-cover" :src="album.image_url" alt="Album Cover" @error="onImageError" @click="() => router.push(`/albums/${album.musicbrainz_album_id}`)" />
    <div class="text-nowrap text-sm">
      {{ album.album }}
    </div>
    <div class="cursor-pointer text-nowrap text-xs text-gray-300" @click="() => router.push(`/artists/${album.musicbrainz_artist_id}`)">
      {{ artistAndDate }}
    </div>
  </div>
  <div v-else-if="props.size === 'xl'" class="h-full flex items-center gap-6 from-zene-600/90 via-zene-600/80 bg-gradient-to-r p-10">
    <img :src="album.image_url" class="size-50 rounded-lg object-cover" @error="onImageError">
    <div class="flex flex-col gap-5">
      <div class="cursor-pointer text-4xl text-white font-bold" @click="() => router.push(`/aLbums/${album.musicbrainz_album_id}`)">
        {{ album.album }}
      </div>
      <div class="cursor-pointer text-xl text-white" @click="() => router.push(`/artists/${album.musicbrainz_artist_id}`)">
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
