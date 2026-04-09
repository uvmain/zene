<script setup lang="ts">
import type { LoadingAttribute } from '~/types'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { artSizes, getCoverArtUrl, onImageError, parseReleaseDate } from '~/logic/common'

const props = defineProps({
  album: { type: Object as PropType<SubsonicAlbum>, required: true },
  showArtist: { type: Boolean, default: true },
  showDate: { type: Boolean, default: true },
  trackTitle: { type: String, required: false },
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
  return getCoverArtUrl(`${props.album.id}`, artSizes.size150)
})

function navigateAlbum() {
  router.push(`/albums/${props.album.id}`)
}

function navigateArtist() {
  router.push(`/artists/${props.album.artistId}`)
}
</script>

<template>
  <div class="flex flex-col gap-2">
    <div class="group grid max-w-150px cursor-pointer" @click="navigateAlbum()">
      <img
        class="col-span-full row-span-full aspect-square shadow-md shadow-background-500 z-1 object-cover dark:shadow-background-950"
        :src="coverArtUrl"
        :loading="loading"
        width="150"
        height="150"
        @error="onImageError"
      />
      <PlayButton
        :album="album"
        class="m-auto pr-1 opacity-0 col-span-full row-span-full scale-50 duration-200 z-2 group-hover:opacity-100 group-hover:scale-100"
      />
    </div>
    <div class="max-w-150px">
      <div v-if="trackTitle" class="text-lg text-primary text-nowrap truncate">
        {{ trackTitle }}
      </div>
      <div v-if="showArtist" class="text-lg text-primary text-nowrap truncate lg:text-base">
        {{ album.title || album.name }}
      </div>
      <div v-if="showArtist && showDate" class="text-sm cursor-pointer text-nowrap truncate" @click="navigateArtist()">
        {{ artistAndDate }}
      </div>
      <div v-else-if="showArtist && !showDate" class="text-sm cursor-pointer text-nowrap truncate" @click="navigateArtist()">
        {{ artist }}
      </div>
      <div v-if="!showArtist && showDate" class="text-sm cursor-pointer text-nowrap truncate lg:text-base" @click="navigateArtist()">
        {{ albumAndDate }}
      </div>
      <div v-else-if="!showArtist && !showDate" class="text-sm cursor-pointer text-nowrap truncate lg:text-base" @click="navigateArtist()">
        {{ album.title }}
      </div>
    </div>
  </div>
</template>
