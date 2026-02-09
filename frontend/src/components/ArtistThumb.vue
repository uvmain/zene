<script setup lang="ts">
import type { LoadingAttribute } from '../types'
import type { SubsonicIndexArtist } from '~/types/subsonicArtist'
import { getCoverArtUrl, onImageError } from '~/composables/logic'

const props = defineProps({
  artist: { type: Object as PropType<SubsonicIndexArtist>, required: true },
  index: { type: Number, default: 0 },
})

const router = useRouter()

const coverArtUrl = computed(() => {
  return getCoverArtUrl(props.artist.coverArt, 120)
})

const loading = computed<LoadingAttribute>(() => {
  return props.index < 10 ? 'eager' : 'lazy'
})

function navigateArtist() {
  router.push(`/artists/${props.artist.id ?? props.artist.musicBrainzId}`)
}
</script>

<template>
  <div>
    <div class="group grid cursor-pointer" @click="navigateArtist()">
      <img
        class="z-1 col-span-full row-span-full rounded-full object-cover"
        :src="coverArtUrl"
        :loading="loading"
        width="150"
        height="150"
        @error="onImageError"
      />
      <PlayButton
        :artist="artist"
        class="z-2 col-span-full row-span-full m-auto scale-50 pr-1 opacity-0 duration-200 group-hover:scale-100 group-hover:opacity-100"
      />
    </div>
    <div class="max-w-150px">
      <div class="truncate text-center text-nowrap text-sm text-primary lg:text-base">
        {{ artist.name }}
      </div>
    </div>
  </div>
</template>
