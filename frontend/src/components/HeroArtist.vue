<script setup lang="ts">
import type { SubsonicArtist } from '~/types/subsonicArtist'
import { artSizes, cacheBustArt, getCoverArtUrl, onImageError } from '~/logic/common'

const props = defineProps({
  artist: { type: Object as PropType<SubsonicArtist>, required: true },
  genres: { type: Array as PropType<string[]>, required: true },
})

const showChangeArtModal = ref(false)
const isStarred = ref<string | undefined>(props.artist.starred)
const artUpdatedTime = ref<string | undefined>(undefined)

const modelArtist = computed(() => {
  return props.artist
})

const coverArtUrl = computed(() => {
  return getCoverArtUrl(props.artist.coverArt, artSizes.size200, artUpdatedTime.value)
})

const artistRoute = computed(() => {
  return `/artists/${props.artist.id}`
})

function actOnUpdatedArt() {
  showChangeArtModal.value = false
  cacheBustArt(`${modelArtist.value.id}`)
  artUpdatedTime.value = Date.now().toString()
}
</script>

<template>
  <div
    class="corner-cut-large w-full shadow-background-500 shadow-md bg-cover bg-center dark:shadow-background-950"
    :style="{ backgroundImage: `url(${coverArtUrl})` }"
  >
    <div class="corner-cut background-grad-2 backdrop-blur-md lg:(corner-cut-large)">
      <div class="p-4 lg:p-8">
        <div class="flex flex-row gap-4 items-center">
          <img
            :src="coverArtUrl"
            class="border-muted rounded-md h-32 aspect-square cursor-pointer shadow-background-500 shadow-md object-cover lg:h-52 dark:shadow-background-900"
            loading="eager"
            @error="onImageError"
          >
          <div class="my-auto text-left flex flex-col gap-1 lg:gap-4">
            <div class="text-xl font-bold line-clamp-1 lg:text-4xl">
              {{ artist.name }}
            </div>
            <Genres v-if="genres.length > 0" :genre-strings="genres" :row-limit="1" />
            <div class="mt-2 flex flex-row gap-4 lg:gap-6">
              <PlayButton class="flex justify-start" :artist="modelArtist" :playing-route="artistRoute" />
              <Fave v-model="isStarred" :musicbrainz-id="modelArtist.id" />
              <Rating v-model="modelArtist.userRating" :musicbrainz-id="modelArtist.id" />
              <ChangeArtIcon @click="showChangeArtModal = true" />
            </div>
          </div>
        </div>
      </div>
    </div>
    <ChangeArtistArt
      v-if="showChangeArtModal"
      :artist="modelArtist"
      @close="showChangeArtModal = false"
      @art-updated="actOnUpdatedArt"
    />
  </div>
</template>
