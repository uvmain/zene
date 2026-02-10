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
      <div class="h-full w-full flex items-center justify-between background-grad-2 backdrop-blur-md">
        <div class="bg-black/10 p-8">
          <div class="h-auto flex flex-row gap-2 lg:gap-6">
            <img
              :src="coverArtUrl"
              class="my-auto aspect-square h-30 cursor-pointer border-muted rounded-md shadow-md shadow-zshade-500 lg:h-52 dark:shadow-zshade-900"
              loading="lazy"
              @error="onImageError"
            >
            <div class="my-auto flex flex-col gap-1 text-left lg:gap-4">
              <div class="text-2xl font-bold">
                {{ podcast.title }}
              </div>
              <div
                class="line-clamp-5 overflow-hidden truncate whitespace-pre-line text-pretty"
                v-html="podcast.description.replaceAll(/\n/g, '<br>')"
              />
              <div v-if="podcast.episode.length && podcast.episode[0].genres?.length > 0" class="hidden lg:(block flex flex-nowrap justify-start gap-2 overflow-hidden)">
                <GenreBottle v-for="genre in podcast.episode[0].genres.slice(0, 8)" :key="genre.name" :genre="genre.name" />
              </div>
              <PlayButton class="flex justify-start" :podcast="podcast.episode[0]" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
