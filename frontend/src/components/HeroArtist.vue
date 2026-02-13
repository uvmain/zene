<script setup lang="ts">
import type { SubsonicArtist } from '~/types/subsonicArtist'
import { artSizes, cacheBustAlbumArt, getCoverArtUrl, onImageError } from '~/logic/common'

const props = defineProps({
  artist: { type: Object as PropType<SubsonicArtist>, required: true },
  genres: { type: Array as PropType<string[]>, required: true },
})

const showChangeArtModal = ref(false)

const coverArtUrl = computed(() => {
  return getCoverArtUrl(props.artist.coverArt, artSizes.size200)
})

// function actOnUpdatedArt() {
//   showChangeArtModal.value = false
//   cacheBustAlbumArt(`${props.artist.id}`)
// }
</script>

<template>
  <section class="corner-cut-large overflow-hidden">
    <div
      class="h-full w-full bg-cover bg-center"
      :style="{ backgroundImage: `url(${coverArtUrl})` }"
    >
      <div class="h-full w-full flex items-center justify-between background-grad-2 backdrop-blur-md">
        <div class="p-8">
          <div class="h-30 flex flex-row gap-2 lg:h-52 lg:gap-6">
            <img
              :src="coverArtUrl"
              class="aspect-square h-30 cursor-pointer border-muted rounded-md shadow-md shadow-zshade-500 lg:h-52 dark:shadow-zshade-900"
              loading="lazy"
              @error="onImageError"
            >
            <div class="my-auto flex flex-col gap-1 text-left lg:gap-4">
              <div class="line-clamp-1 cursor-pointer text-xl font-bold lg:text-4xl">
                {{ artist.name }}
              </div>
              <div v-if="genres.length > 0" class="hidden lg:(block flex flex-nowrap justify-start gap-2 overflow-hidden)">
                <GenreBottle v-for="genre in genres.slice(0, 8)" :key="genre" :genre="genre" />
              </div>
              <PlayButton class="flex justify-start" :artist="artist" />
            </div>
          </div>
          <div class="absolute right-2 top-2 opacity-50 hover:opacity-100">
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
  </section>
</template>
