<script setup lang="ts">
import type { LoadingAttribute } from '~/types'
import type { SubsonicIndexArtist } from '~/types/subsonicArtist'
import { artSizes, getCoverArtUrl, onImageError } from '~/logic/common'

const props = defineProps({
  artist: { type: Object as PropType<SubsonicIndexArtist>, required: true },
  index: { type: Number, default: 0 },
})

const router = useRouter()

const coverArtUrl = computed(() => {
  return getCoverArtUrl(props.artist.coverArt, artSizes.size150)
})

const loading = computed<LoadingAttribute>(() => {
  return props.index < 10 ? 'eager' : 'lazy'
})

function navigateArtist() {
  router.push(`/artists/${props.artist.id ?? props.artist.musicBrainzId}`)
}
</script>

<template>
  <div class="flex flex-col gap-2">
    <div class="group grid max-w-150px cursor-pointer" @click="navigateArtist()">
      <img
        class="rounded-full col-span-full row-span-full aspect-square shadow-md shadow-background-500 z-1 object-cover dark:shadow-background-950"
        :src="coverArtUrl"
        :loading="loading"
        width="150"
        height="150"
        @error="onImageError"
      />
      <PlayButton
        :artist="artist"
        class="m-auto pr-1 opacity-0 col-span-full row-span-full scale-50 duration-200 z-2 group-hover:opacity-100 group-hover:scale-100"
      />
    </div>
    <div class="max-w-150px">
      <div class="text-sm text-primary text-center text-nowrap truncate lg:text-base">
        {{ artist.name }}
      </div>
    </div>
  </div>
</template>
