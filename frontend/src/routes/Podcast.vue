<script setup lang="ts">
import type { SubsonicPodcastChannelsResponse } from '~/types/subsonic'
import type { SubsonicPodcastChannel } from '~/types/subsonicPodcasts'
import { getStreamUrl, openSubsonicFetchRequest } from '~/composables/backendFetch'

const route = useRoute()

watch(() => route.params.podcast_id, async () => {
  getPodcast()
})

const podcast = ref<SubsonicPodcastChannel>()

async function getPodcast() {
  const formData = new FormData()
  // formData.append('includeEpisodes', true.toString())
  formData.append('id', route.params.podcast_id.toString())
  const response = await openSubsonicFetchRequest<SubsonicPodcastChannelsResponse>('getPodcasts', {
    body: formData,
  })
  podcast.value = response?.podcasts?.channel[0]
  podcast.value.coverArt = `/share/img/${podcast.value.coverArt}?size=400`

  for (const ep of podcast.value.episode) {
    ep.coverArt = `/share/img/${ep.coverArt}?size=400`
  }
}

async function downloadEpisode(episodeId: string) {
  if (!podcast.value)
    return
  const formData = new FormData()
  formData.append('id', episodeId)
  openSubsonicFetchRequest<SubsonicPodcastChannelsResponse>('downloadPodcastEpisode', {
    body: formData,
  })
  podcast.value.episode.find(episode => episode.id === episodeId)!.status = 'downloading'
}

function playEpisodeInNewTab(episodeId: string) {
  if (!podcast.value)
    return
  const episode = podcast.value.episode.find(ep => ep.id === episodeId)
  if (!episode || episode.status !== 'completed')
    return
  window.open(getStreamUrl('stream', new URLSearchParams({ id: episode.streamId })), '_blank')?.focus()
}

onBeforeMount(getPodcast)
</script>

<template>
  <div>
    <div v-if="!podcast" class="text-gray-500">
      Podcast not found.
    </div>
    <div v-else class="flex flex-col gap-6">
      <div class="pb-4">
        <center>
          <h2 class="mb-4 text-2xl font-bold">
            {{ podcast.title }}
          </h2>
        </center>
        <div class="mx-auto max-w-60dvw flex flex-row gap-4 align-top">
          <img
            :src="podcast.coverArt"
            alt="Podcast Cover"
            class="size-70 rounded-lg object-cover"
          />
          <div class="line-clamp-12 max-h-70 overflow-hidden text-ellipsis whitespace-normal text-pretty text-white text-op-80">
            {{ podcast.description }}
          </div>
        </div>
      </div>
      <div v-if="podcast.lastRefresh === ''" class="mx-auto">
        Episodes are being refreshed...
      </div>
      <div v-for="episode in podcast.episode" :key="episode.id">
        <div class="mx-auto max-w-60dvw flex flex-row justify-start gap-4 align-top transition duration-200 hover:scale-101">
          <div class="grid items-end justify-items-end">
            <img
              :src="episode.coverArt"
              alt="Podcast Cover"
              class="z-1 col-span-full row-span-full my-auto h-48 w-48 rounded-lg object-cover"
            />
            <div class="z-2 col-span-full row-span-full m-2 size-12 hover:text-zene-200">
              <icon-tabler-progress-down v-if="episode.status === 'downloading'" class="size-8 rounded bg-dark bg-opacity-50 p-2" />
              <icon-tabler-play
                v-else-if="episode.status === 'completed'"
                class="size-8 rounded bg-dark bg-opacity-50 p-2 outline-3 outline-green outline-solid -outline-offset-3"
                @click="playEpisodeInNewTab(episode.id)"
              />
              <icon-tabler-download
                v-else
                class="size-8 rounded bg-dark bg-opacity-50 p-2"
                @click="downloadEpisode(episode.id)"
              />
            </div>
          </div>
          <div class="my-auto">
            <div class="text-lg font-semibold">
              {{ episode.title }}
            </div>
            <div class="line-clamp-6 max-h-30 overflow-hidden text-ellipsis whitespace-normal text-pretty text-white text-op-80">
              {{ episode.description }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
