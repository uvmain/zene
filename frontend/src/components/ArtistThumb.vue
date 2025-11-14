<script setup lang="ts">
import type { LoadingAttribute } from '../types'
import type { SubsonicIndexArtist } from '~/types/subsonicArtist'
import { getCoverArtUrl, onImageError } from '~/composables/logic'
import { useSearch } from '~/composables/useSearch'

const props = defineProps({
  artist: { type: Object as PropType<SubsonicIndexArtist>, required: true },
  index: { type: Number, default: 0 },
})

const { closeSearch } = useSearch()

const router = useRouter()

const coverArtUrl = computed(() => {
  return getCoverArtUrl(props.artist.coverArt, 120)
})

const loading = computed<LoadingAttribute>(() => {
  return props.index < 10 ? 'eager' : 'lazy'
})

function navigate() {
  closeSearch()
  const artistId = props.artist.id ?? props.artist.musicBrainzId
  router.push(`/artists/${artistId}`)
}
</script>

<template>
  <div class="group" @click="navigate()">
    <img
      class="aspect-square w-full cursor-pointer object-cover"
      :src="coverArtUrl"
      :loading="loading"
      width="150"
      height="150"
      @error="onImageError"
    />
    <div class="relative">
      <PlayButton
        :artist="artist"
        class="absolute bottom-2 right-10 z-10 opacity-0 transition-all duration-300 group-hover:right-6 group-hover:opacity-100"
      />
    </div>
    <div class="max-w-150px truncate text-nowrap text-sm text-primary">
      {{ artist.name }}
    </div>
  </div>
</template>
