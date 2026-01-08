<script setup lang="ts">
import type { LoadingAttribute } from '../types'
import type { SubsonicAlbum } from '../types/subsonicAlbum'
import { cacheBustAlbumArt, getCoverArtUrl, onImageError, parseReleaseDate } from '~/composables/logic'
import { useSearch } from '../composables/useSearch'

type AlbumSize = 'sm' | 'md' | 'lg'

const props = defineProps({
  album: { type: Object as PropType<SubsonicAlbum>, required: true },
  size: { type: String as PropType<AlbumSize>, default: 'sm' },
  showChangeArtButton: { type: Boolean, default: false },
  showArtist: { type: Boolean, default: true },
  showDate: { type: Boolean, default: true },
  index: { type: Number, default: 0 },
})

const artist = computed(() => {
  return props.album.displayAlbumArtist ?? props.album.artist ?? props.album.displayArtist ?? 'Unknown Artist'
})

const router = useRouter()
const { closeSearch } = useSearch()

const showChangeArtModal = ref(false)
const artUpdatedTime = ref<string | undefined>(undefined)

const artistAndDate = computed(() => {
  if (props.album.releaseDate) {
    return `${artist.value} • ${parseReleaseDate(props.album.releaseDate)}`
  }
  else if (props.album.year) {
    return `${artist.value} • ${props.album.year}`
  }
  else {
    return artist.value
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
  return getCoverArtUrl(`${props.album.id}`, 150, artUpdatedTime.value)
})

const coverArtUrlMd = computed(() => {
  return getCoverArtUrl(`${props.album.id}`, 200, artUpdatedTime.value)
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
  cacheBustAlbumArt(`${props.album.id}`)
  artUpdatedTime.value = Date.now().toString()
}
</script>

<template>
  <div>
    <div v-if="size === 'sm'" class="group">
      <img
        class="aspect-square h-full w-full cursor-pointer border-muted"
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
          class="absolute bottom-2 right-1 z-10 opacity-0 transition-all duration-100 group-hover:opacity-100"
        />
      </div>
      <div class="max-w-150px">
        <div v-if="showArtist" class="truncate text-nowrap text-lg text-primary lg:text-base">
          {{ album.title || album.name }}
        </div>
        <div v-if="showArtist && showDate" class="cursor-pointer truncate text-nowrap text-sm" @click="navigateArtist()">
          {{ artistAndDate }}
        </div>
        <div v-else-if="showArtist && !showDate" class="cursor-pointer truncate text-nowrap text-sm" @click="navigateArtist()">
          {{ artist }}
        </div>
        <div v-if="!showArtist && showDate" class="cursor-pointer truncate text-nowrap text-sm lg:text-base" @click="navigateArtist()">
          {{ albumAndDate }}
        </div>
        <div v-else-if="!showArtist && !showDate" class="cursor-pointer truncate text-nowrap text-sm lg:text-base" @click="navigateArtist()">
          {{ album.title }}
        </div>
      </div>
    </div>
    <!-- medium -->
    <div v-else-if="props.size === 'md'" class="group corner-cut relative background-grad-2 p-3 lg:(corner-cut-large p-8)">
      <div class="h-30 flex flex-row gap-2 lg:h-52 lg:gap-6">
        <img
          :src="coverArtUrlMd"
          class="h-30 cursor-pointer border-muted lg:h-52"
          loading="lazy"
          @error="onImageError"
          @click="navigateAlbum()"
        >
        <div class="my-auto flex flex-col gap-1 text-left lg:gap-4">
          <div class="line-clamp-1 cursor-pointer text-xl font-bold lg:text-4xl" @click="navigateAlbum()">
            {{ album.name }}
          </div>
          <div class="cursor-pointer text-lg lg:text-xl" @click="navigateArtist()">
            {{ artistAndDate }}
          </div>
          <div v-if="album.genres?.length > 0" class="hidden lg:(block flex flex-nowrap justify-start gap-2 overflow-hidden)">
            <GenreBottle v-for="genre in album.genres.filter(g => g.name !== '').slice(0, 8)" :key="genre.name" :genre="genre.name" />
          </div>
          <PlayButton class="flex justify-start" :album="album" />
        </div>
      </div>
      <!-- Change Album Art section -->
      <div v-if="showChangeArtButton" class="absolute right-2 top-2">
        <ZButton
          class="opacity-0 group-hover:opacity-100"
          @click="showChangeArtModal = true"
        >
          <div>
            Change Album Art
          </div>
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
