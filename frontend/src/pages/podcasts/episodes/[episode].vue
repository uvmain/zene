<script setup lang="ts">
import type { SubsonicPodcastEpisodesResponse } from '~/types/subsonic'
import type { SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import { openSubsonicFetchRequest } from '~/composables/backendFetch'

const route = useRoute('/podcasts/episodes/[episode]')

watch(() => route.params.episode, async () => {
  getEpisode()
})

const episode = ref<SubsonicPodcastEpisode>()

async function getEpisode() {
  const formData = new FormData()
  formData.append('id', route.params.episode.toString())
  const response = await openSubsonicFetchRequest<SubsonicPodcastEpisodesResponse>('getPodcastEpisode', {
    body: formData,
  })
  episode.value = response?.podcastEpisode
}

const coverArt = computed(() => {
  if (!episode.value)
    return ''
  return `/share/img/${episode.value.coverArt}?size=400`
})

onBeforeMount(async () => {
  await getEpisode()
})
</script>

<template>
  <div>
    <div v-if="!episode" class="text-primary">
      Podcast episode not found.
    </div>
    <div v-else class="mx-auto max-w-60dvw flex flex-col gap-6">
      <div class="flex flex-row gap-4">
        <img
          :src="coverArt"
          alt="Podcast Cover"
          class="size-70 rounded-lg object-cover"
          width="280"
          height="280"
          loading="eager"
        />
        <div class="my-auto flex flex-col gap-4">
          <div class="text-4xl font-bold">
            {{ episode.title }}
          </div>
          <div v-if="episode.genres?.length > 0" class="flex flex-wrap justify-center gap-2 lg:justify-start">
            <ZInfo v-for="genre in episode.genres?.filter(g => g.name !== '')" :key="genre.name" :text="genre.name" />
          </div>
          <PlayButton
            :podcast="episode"
            class="my-auto"
            hover-text="Play episode"
          />
        </div>
      </div>
      <div
        class="whitespace-pre-line text-pretty text-op-80"
        v-html="episode.description.replaceAll(/\n/g, '<br>')"
      />
    </div>
  </div>
</template>
