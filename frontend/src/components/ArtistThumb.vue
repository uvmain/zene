<script setup lang="ts">
import type { SubsonicIndexArtist } from '~/types/subsonicArtist'
import { getCoverArtUrl, onImageError } from '~/composables/logic'
import { useSearch } from '~/composables/useSearch'

const props = defineProps({
  artist: { type: Object as PropType<SubsonicIndexArtist>, required: true },
})

const { closeSearch } = useSearch()

const router = useRouter()

const coverArtUrl = computed(() => {
  return getCoverArtUrl(props.artist.coverArt, 200)
})

function navigate() {
  closeSearch()
  router.push(`/artists/${props.artist.musicBrainzId}`)
}
</script>

<template>
  <div class="w-30 flex flex-col cursor-pointer gap-2" @click="navigate()">
    <div class="group size-30">
      <img
        class="h-full w-full object-cover"
        :src="coverArtUrl"
        loading="lazy"
        width="200"
        height="200"
        @error="onImageError"
      />
      <div class="relative">
        <PlayButton
          :artist="artist"
          class="invisible absolute bottom-2 right-1 z-10 group-hover:visible"
        />
      </div>
    </div>
    <div class="text-nowrap text-xs text-zgray-200">
      {{ artist.name }}
    </div>
  </div>
</template>
