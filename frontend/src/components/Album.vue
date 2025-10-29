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
  index: { type: Number, default: 0 },
})

const router = useRouter()
const { closeSearch } = useSearch()

const showChangeArtModal = ref(false)

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

const loading = computed<LoadingAttribute>(() => {
  return props.index < 10 ? 'eager' : 'lazy'
})

const updatedTime = ref<Date | null>(null)

const coverArtUrlSm = computed(() => {
  return updatedTime.value ? `${getCoverArtUrl(props.album.id, 120)}&time=${updatedTime.value.getTime()}` : `${getCoverArtUrl(props.album.id, 120)}`
})

const coverArtUrlMd = computed(() => {
  return updatedTime.value ? `${getCoverArtUrl(props.album.id, 200)}&time=${updatedTime.value.getTime()}` : `${getCoverArtUrl(props.album.id, 200)}`
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
  updatedTime.value = new Date()
}
</script>

<template>
  <div>
    <div v-if="size === 'sm'" class="group">
      <img
        class="h-24 w-24 cursor-pointer object-cover md:size-30"
        :src="coverArtUrlSm"
        alt="Album Cover"
        :loading="loading"
        width="120"
        height="120"
        @error="onImageError" @click="navigateAlbum()"
      />
      <div class="relative">
        <PlayButton
          :album="album"
          class="absolute bottom-2 right-1 z-10 opacity-0 transition-all duration-300 group-hover:opacity-100"
        />
      </div>
      <div class="w-24 truncate text-nowrap text-xs text-primary md:w-30 md:text-sm">
        {{ album.name }}
      </div>
      <div class="w-24 cursor-pointer truncate text-nowrap text-xs md:w-30" @click="navigateArtist()">
        {{ artistAndDate }}
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
      <div class="flex flex-col gap-2 text-center md:gap-5 md:text-left">
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
