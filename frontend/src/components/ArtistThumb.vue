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
  <div>
    <div class="group" @click="navigate()">
      <img
        class="size-150px cursor-pointer rounded-full object-cover"
        :src="coverArtUrl"
        :loading="loading"
        width="150"
        height="150"
        @error="onImageError"
      />
      <div class="relative">
        <PlayButton
          :artist="artist"
          class="absolute bottom-2 right-4 z-10 opacity-0 transition-all duration-300 group-hover:right-1 group-hover:opacity-100"
        />
      </div>
    </div>
    <div class="max-w-150px">
      <div class="truncate text-center text-nowrap text-sm text-primary md:text-base">
        {{ artist.name }}
      </div>
    </div>
  </div>
</template>
