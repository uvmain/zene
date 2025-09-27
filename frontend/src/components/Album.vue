<script setup lang="ts">
import type { SubsonicAlbum } from '../types/subsonicAlbum'
import { getCoverArtUrl, onImageError, parseReleaseDate } from '~/composables/logic'
import { useSearch } from '../composables/useSearch'

const props = defineProps({
  album: { type: Object as PropType<SubsonicAlbum>, required: true },
  size: { type: String, default: 'md' },
})

const router = useRouter()
const { closeSearch } = useSearch()

const artistAndDate = computed(() => {
  return `${props.album.artist} • ${parseReleaseDate(props.album.releaseDate)}`
})

const coverArtUrl = computed(() => {
  return getCoverArtUrl(props.album.id, 200)
})

const coverArtUrlLarge = computed(() => {
  return getCoverArtUrl(props.album.id, 400)
})

function navigateAlbum() {
  closeSearch()
  router.push(`/albums/${props.album.id}`)
}

function navigateArtist() {
  closeSearch()
  router.push(`/artists/${props.album.artistId}`)
}
</script>

<template>
  <div>
    <div v-if="size === 'lg'" class="group h-32 w-24 md:h-40 md:w-30">
      <img
        class="h-24 w-24 cursor-pointer rounded-lg object-cover md:size-30"
        :src="coverArtUrl"
        alt="Album Cover"
        loading="lazy"
        width="200"
        height="200"
        @error="onImageError" @click="navigateAlbum()"
      />
      <div class="relative">
        <PlayButton
          size="small"
          :album="album"
          class="invisible absolute bottom-2 right-1 z-10 group-hover:visible"
        />
      </div>
      <div class="w-24 truncate text-nowrap text-xs md:w-30 md:text-sm">
        {{ album.name }}
      </div>
      <div class="w-24 cursor-pointer truncate text-nowrap text-xs text-gray-300 md:w-30" @click="navigateArtist()">
        {{ artistAndDate }}
      </div>
    </div>
    <div v-else-if="props.size === 'xl'" class="h-full flex flex-col items-center gap-2 from-zene-600/90 via-zene-600/80 bg-gradient-to-r p-3 md:flex-row md:gap-6 md:p-10">
      <img
        :src="coverArtUrlLarge"
        class="h-24 w-24 cursor-pointer rounded-lg object-cover md:size-50"
        loading="lazy"
        width="400"
        height="400"
        @error="onImageError"
        @click="navigateAlbum()"
      >
      <div class="flex flex-col gap-2 text-center md:gap-5 md:text-left">
        <div class="cursor-pointer text-lg text-white font-bold md:text-4xl" @click="navigateAlbum()">
          {{ album.name }}
        </div>
        <div class="cursor-pointer text-sm text-white md:text-xl" @click="navigateArtist()">
          {{ artistAndDate }}
        </div>
        <div v-if="album.genres.length > 0" class="flex flex-wrap justify-center gap-2 md:justify-start">
          <GenreBottle v-for="genre in album.genres" :key="genre.name" :genre="genre.name" />
        </div>
        <div class="flex justify-center md:justify-start">
          <PlayButton :album="album" />
        </div>
      </div>
    </div>
  </div>
</template>
