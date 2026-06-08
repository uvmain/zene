<script setup lang="ts">
import type { SubsonicPodcastChannel } from '~/types/subsonicPodcasts'
import { artSizes, getCoverArtUrl, onImageError } from '~/logic/common'

const props = defineProps({
  podcast: { type: Object as PropType<SubsonicPodcastChannel>, required: true },
})

const coverArtUrl = computed(() => {
  return getCoverArtUrl(props.podcast.coverArt, artSizes.size200)
})

const genres = computed(() => {
  if (props.podcast.episode.length > 0) {
    return props.podcast.episode[0].genres?.map(genre => genre.name) || []
  }
  return []
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
              class="my-auto border-muted rounded-md aspect-square cursor-pointer shadow-background-500 shadow-md dark:shadow-background-900"
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
              <Genres v-if="genres.length > 0" :genre-strings="genres" :row-limit="1" />
              <!-- <PlayButton class="flex justify-start" :podcast-episode="podcast.episode[0]" /> -->
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
