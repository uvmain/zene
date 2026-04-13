<script setup lang="ts">
import type { SubsonicArtist } from '~/types/subsonicArtist'
import { artSizes, getCoverArtUrl, onImageError } from '~/logic/common'

const props = defineProps({
  artist: { type: Object as PropType<SubsonicArtist>, required: true },
  genres: { type: Array as PropType<string[]>, required: true },
})

const showChangeArtModal = ref(false)
const isStarred = ref<string | undefined>(props.artist.starred)
const modelArtist = computed(() => {
  return props.artist
})

const coverArtUrl = computed(() => {
  return getCoverArtUrl(props.artist.coverArt, artSizes.size200)
})
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
            class="border-muted rounded-md h-32 aspect-square cursor-pointer shadow-background-500 shadow-md lg:h-52 dark:shadow-background-900"
            loading="lazy"
            @error="onImageError"
          >
          <div class="my-auto text-left flex flex-col gap-1 lg:gap-4">
            <div class="text-xl font-bold line-clamp-1 lg:text-4xl">
              {{ artist.name }}
            </div>
            <div v-if="genres.length > 0" class="flex-wrap gap-2 hidden justify-start overflow-hidden md:flex">
              <GenreBottle v-for="genre in genres.slice(0, 3)" :key="genre" :genre="genre" class="hidden md:block lg:hidden" />
              <GenreBottle v-for="genre in genres.slice(0, 6)" :key="genre" :genre="genre" class="hidden lg:block xl:hidden" />
              <GenreBottle v-for="genre in genres.slice(0, 9)" :key="genre" :genre="genre" class="hidden xl:block 2xl:hidden" />
              <GenreBottle v-for="genre in genres" :key="genre" :genre="genre" class="hidden 2xl:block" />
            </div>
            <div class="mt-2 flex flex-row gap-8">
              <PlayButton class="flex justify-start" :artist="modelArtist" />
              <Fave v-model="isStarred" :musicbrainz-id="modelArtist.id" />
              <Rating v-model="modelArtist.userRating" :musicbrainz-id="modelArtist.id" />
            </div>
          </div>
        </div>
        <div class="opacity-50 right-2 top-2 absolute hover:opacity-100">
          <!-- Change Album Art section -->
          <ZButton
            :disabled="true"
            @click="showChangeArtModal = true"
          >
            <div>
              Change Art
            </div>
          </ZButton>
          <!-- <ChangeAlbumArt
              v-if="showChangeArtModal"
              :album="albumArray[index]"
              @close="showChangeArtModal = false"
              @art-updated="actOnUpdatedArt"
            /> -->
        </div>
      </div>
    </div>
  </div>
</template>
