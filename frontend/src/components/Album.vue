<script setup lang="ts">
import type { LoadingAttribute } from '../types'
import type { SubsonicAlbum } from '../types/subsonicAlbum'
import { albumArtSizes, getCoverArtUrl, onImageError, parseReleaseDate } from '~/logic/common'

const props = defineProps({
  album: { type: Object as PropType<SubsonicAlbum>, required: true },
  showArtist: { type: Boolean, default: true },
  showDate: { type: Boolean, default: true },
  index: { type: Number, default: 0 },
})

const router = useRouter()

const artist = computed(() => {
  return props.album.displayAlbumArtist ?? props.album.artist ?? props.album.displayArtist ?? 'Unknown Artist'
})

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

const coverArtUrl = computed(() => {
  return getCoverArtUrl(`${props.album.id}`, albumArtSizes.size150)
})

function navigateAlbum() {
  router.push(`/albums/${props.album.id}`)
}

function navigateArtist() {
  router.push(`/artists/${props.album.artistId}`)
}
</script>

<template>
  <div>
    <div class="group max-w-150px">
      <img
        class="aspect-square cursor-pointer border-muted rounded-md shadow-sm shadow-zshade-500 dark:shadow-zshade-900"
        :src="coverArtUrl"
        alt="Album Cover"
        :loading="loading"
        width="150"
        height="150"
        @error="onImageError"
        @click="navigateAlbum()"
      />
      <div class="relative">
        <PlayButton
          :album="album"
          class="absolute bottom-2 right-1 z-10 scale-50 opacity-0 transition-all duration-200 group-hover:scale-100 group-hover:opacity-100"
        />
      </div>
      <div>
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
  </div>
</template>
