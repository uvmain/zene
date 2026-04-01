<script setup lang="ts">
import type { SubsonicPodcastChannel } from '~/types/subsonicPodcasts'
import { artSizes, getCoverArtUrl, onImageError } from '~/logic/common'

const props = defineProps({
  podcast: { type: Object as PropType<SubsonicPodcastChannel>, required: true },
})

const coverArtUrl = computed(() => {
  return getCoverArtUrl(props.podcast.coverArt, artSizes.size200)
})
</script>

<template>
  <section class="corner-cut-large cursor-pointer overflow-hidden">
    <div
      class="h-full w-full bg-cover bg-center"
      :style="{ backgroundImage: `url(${coverArtUrl})` }"
    >
      <div class="flex h-full w-full items-center justify-between background-grad-2 backdrop-blur-md">
        <div class="p-8 bg-black/10">
          <div class="flex flex-row gap-2 h-auto lg:gap-6">
            <img
              :src="coverArtUrl"
              class="my-auto border-muted rounded-md aspect-square cursor-pointer shadow-md shadow-zshade-500 dark:shadow-zshade-900"
              loading="lazy"
              width="200"
              height="200"
              @error="onImageError"
            >
            <div class="my-auto text-left flex flex-col gap-1 lg:gap-4">
              <div class="text-2xl font-bold">
                {{ podcast.title }}
              </div>
              <div
                class="whitespace-pre-line text-pretty truncate overflow-hidden line-clamp-5"
                v-html="podcast.description.replaceAll(/\n/g, '<br>')"
              />
              <div v-if="podcast.episode.length && podcast.episode[0].genres?.length > 0" class="hidden lg:(flex flex-nowrap gap-2 justify-start overflow-hidden)">
                <GenreBottle v-for="genre in podcast.episode[0].genres.slice(0, 8)" :key="genre.name" :genre="genre.name" />
              </div>
              <!-- <PlayButton class="flex justify-start" :podcast-episode="podcast.episode[0]" /> -->
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
