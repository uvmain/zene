<script setup lang="ts">
import type { ArtistMetadata } from '../types'
import { useSearch } from '../composables/useSearch'

const props = defineProps({
  artist: { type: Object as PropType<ArtistMetadata>, required: true },
})

const { closeSearch } = useSearch()

const router = useRouter()

function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}

function navigate() {
  closeSearch()
  router.push(`/artists/${props.artist.musicbrainz_artist_id}`)
}
</script>

<template>
  <div class="w-30 flex flex-col cursor-pointer gap-2" @click="navigate()">
    <div class="group size-30">
      <img
        class="h-full w-full rounded-md object-cover"
        :src="artist.image_url"
        @error="onImageError"
      />
      <div class="relative">
        <PlayButton
          size="small"
          :artist="artist"
          class="invisible absolute bottom-2 right-1 z-10 group-hover:visible"
        />
      </div>
    </div>
    <div class="text-nowrap text-xs text-gray-300">
      {{ artist.artist }}
    </div>
  </div>
</template>
