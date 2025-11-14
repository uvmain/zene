<script setup lang="ts">
import type { LoadingAttribute } from '../types'
import type { SubsonicAlbum } from '../types/subsonicAlbum'
import { getCoverArtUrl, onImageError, parseReleaseDate } from '~/composables/logic'
import { useSearch } from '../composables/useSearch'

type AlbumSize = 'sm' | 'md' | 'lg'

const props = defineProps({
  album: { type: Object as PropType<SubsonicAlbum>, required: true },
  size: { type: String as PropType<AlbumSize>, default: 'sm' },
  showChangeArtButton: { type: Boolean, default: false },
  showArtist: { type: Boolean, default: true },
  index: { type: Number, default: 0 },
})

const router = useRouter()
const { closeSearch } = useSearch()

const showChangeArtModal = ref(false)
const artUpdatedTime = ref<string | undefined>(undefined)

const artistAndDate = computed(() => {
  const artist = props.album.displayAlbumArtist ?? props.album.displayArtist ?? props.album.artist ?? 'Unknown Artist'
  if (props.album.releaseDate) {
    return `${artist} • ${parseReleaseDate(props.album.releaseDate)}`
  }
  else if (props.album.year) {
    return `${artist} • ${props.album.year}`
  }
  else {
    return artist
  }
})

const albumAndDate = computed(() => {
  const album = props.album.title || props.album.name || 'Unknown Album'
  if (props.album.releaseDate) {
    return `${album} • ${parseReleaseDate(props.album.releaseDate)}`
  }
  else if (props.album.year) {
    return `${album} • ${props.album.year}`
  }
  else {
    return album
  }
})

const loading = computed<LoadingAttribute>(() => {
  return props.index < 10 ? 'eager' : 'lazy'
})

const coverArtUrlSm = computed(() => {
  return getCoverArtUrl(props.album.id, 150, artUpdatedTime.value)
})

const coverArtUrlMd = computed(() => {
  return getCoverArtUrl(props.album.id, 200, artUpdatedTime.value)
})

function navigateAlbum() {
  closeSearch()
  router.push(`/albums/${props.album.id}`)
}

function navigateArtist() {
  closeSearch()
  router.push(`/artists/${props.album.artistId}`)
}

function actOnUpdatedArt() {
  showChangeArtModal.value = false
  // cache bust
  fetch(getCoverArtUrl(props.album.id, 120), { method: 'POST', credentials: 'include' })
  fetch(getCoverArtUrl(props.album.id, 200), { method: 'POST', credentials: 'include' })
  artUpdatedTime.value = Date.now().toString()
}
</script>

<template>
  <div>
    <div v-if="size === 'sm'" class="group">
      <img
        class="aspect-square h-full w-full cursor-pointer object-cover"
        :src="coverArtUrlSm"
        alt="Album Cover"
        :loading="loading"
        width="150"
        height="150"
        @error="onImageError" @click="navigateAlbum()"
      />
      <div class="relative">
        <PlayButton
          :album="album"
          class="absolute bottom-2 right-4 z-10 opacity-0 transition-all duration-300 group-hover:right-1 group-hover:opacity-100"
        />
      </div>
      <div class="max-w-150px">
        <div v-if="showArtist" class="truncate text-nowrap text-sm text-primary md:text-base">
          {{ album.title || album.name }}
        </div>
        <div v-if="showArtist" class="cursor-pointer truncate text-nowrap text-sm" @click="navigateArtist()">
          {{ artistAndDate }}
        </div>
        <div v-if="!showArtist" class="cursor-pointer truncate text-nowrap text-sm md:text-base" @click="navigateArtist()">
          {{ albumAndDate }}
        </div>
      </div>
    </div>
    <div v-else-if="props.size === 'md'" class="group corner-cut-large relative h-full flex flex-col items-center gap-2 background-grad-2 p-3 md:flex-row md:gap-6 md:p-10">
      <img
        :src="coverArtUrlMd"
        class="size-24 cursor-pointer object-cover md:size-52"
        loading="lazy"
        width="200"
        height="200"
        @error="onImageError"
        @click="navigateAlbum()"
      >
      <div class="flex flex-col gap-2 text-center md:gap-4 md:text-left">
        <div class="cursor-pointer text-lg font-bold md:text-4xl" @click="navigateAlbum()">
          {{ album.name }}
        </div>
        <div class="cursor-pointer text-sm md:text-xl" @click="navigateArtist()">
          {{ artistAndDate }}
        </div>
        <div v-if="album.genres?.length > 0" class="flex justify-center gap-2 overflow-hidden md:flex-nowrap md:justify-start">
          <GenreBottle v-for="genre in album.genres.filter(g => g.name !== '').slice(0, 8)" :key="genre.name" :genre="genre.name" />
        </div>
        <div class="flex justify-center md:justify-start">
          <PlayButton :album="album" />
        </div>
      </div>
      <!-- Change Album Art section -->
      <div v-if="showChangeArtButton">
        <ZButton
          class="absolute right-2 top-2 opacity-0 group-hover:opacity-100"
          @click="showChangeArtModal = true"
        >
          Change Album Art
        </ZButton>
        <ChangeAlbumArt
          v-if="showChangeArtModal"
          :album="album"
          @close="showChangeArtModal = false"
          @art-updated="actOnUpdatedArt"
        />
      </div>
    </div>
  </div>
</template>
