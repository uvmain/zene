<script setup lang="ts">
import type { SubsonicPodcastChannelsResponse } from '~/types/subsonic'
import type { SubsonicPodcastChannel } from '~/types/subsonicPodcasts'
import { getStreamUrl, openSubsonicFetchRequest } from '~/composables/backendFetch'

const route = useRoute()
const router = useRouter()

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

const showDeleteChannelModal = ref(false)
function confirmDeletePodcast() {
  deletePodcastChannel()
  showDeleteChannelModal.value = false
}

async function deletePodcastChannel() {
  if (!podcast.value)
    return
  const formData = new FormData()
  formData.append('id', podcast.value.id)
  const response = await openSubsonicFetchRequest<SubsonicPodcastChannelsResponse>('deletePodcastChannel', {
    body: formData,
  })
  if (response?.status === 'ok') {
    router.push('/podcasts')
  }
}

onBeforeMount(getPodcast)
</script>

<template>
  <div>
    <div v-if="!podcast" class="text-zgray-200">
      Podcast not found.
    </div>
    <div v-else class="flex flex-col gap-6">
      <!-- header -->
      <div class="pb-4">
        <div class="group relative mx-auto max-w-60dvw flex flex-row gap-4 align-top">
          <button
            class="z-button absolute right-0 opacity-0 group-hover:opacity-100"
            @click="showDeleteChannelModal = true"
          >
            Delete Podcast
          </button>
          <img
            :src="podcast.coverArt"
            alt="Podcast Cover"
            class="size-70 object-cover"
          />
          <div class="my-auto flex flex-col gap-4">
            <div class="mb-4 text-2xl font-bold">
              {{ podcast.title }}
            </div>
            <div
              class="line-clamp-6 max-h-70 overflow-hidden text-ellipsis whitespace-pre-line text-pretty text-op-80"
              v-html="podcast.description.replaceAll(/\n/g, '<br>')"
            />
            <div v-if="podcast.episode[0].genres?.length > 0" class="flex flex-wrap justify-center gap-2 md:justify-start">
              <GenreBottle v-for="genre in podcast.episode[0].genres" :key="genre.name" :genre="genre.name" />
            </div>
          </div>
        </div>
      </div>
      <!-- episodes -->
      <div v-if="podcast.lastRefresh === ''" class="mx-auto">
        Episodes are being refreshed...
      </div>
      <div v-for="episode in podcast.episode" :key="episode.id">
        <div class="mx-auto max-w-60dvw flex flex-row justify-start gap-4 align-top transition duration-150 hover:scale-101">
          <div class="grid items-end justify-items-end">
            <img
              :src="episode.coverArt"
              alt="Podcast Cover"
              class="z-1 col-span-full row-span-full my-auto h-48 w-48 object-cover"
            />
            <div class="z-2 col-span-full row-span-full m-2 size-12 hover:text-zene-200">
              <icon-nrk-progress v-if="episode.status === 'downloading'" class="size-8 bg-dark bg-opacity-50 p-2" />
              <icon-nrk-media-play
                v-else-if="episode.status === 'completed'"
                class="size-8 bg-dark bg-opacity-50 p-2 outline-3 outline-green outline-solid -outline-offset-3"
                @click="playEpisodeInNewTab(episode.id)"
              />
              <icon-nrk-download
                v-else
                class="size-8 bg-dark bg-opacity-50 p-2"
                @click="downloadEpisode(episode.id)"
              />
            </div>
          </div>
          <div class="my-auto">
            <div class="text-lg font-semibold">
              {{ episode.title }}
            </div>
            <div
              class="line-clamp-6 max-h-30 overflow-hidden text-ellipsis whitespace-normal text-pretty text-op-80"
              v-html="episode.description.replaceAll(/\n/g, '<br>')"
            />
          </div>
        </div>
      </div>
    </div>
    <div v-if="showDeleteChannelModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div class="w-80 bg-dark p-6 text-center shadow-lg">
        <div class="mb-4 text-lg font-semibold">
          Are you sure you want to delete this podcast channel?
        </div>
        <div class="mt-6 flex justify-center gap-4">
          <button class="bg-gray-400 px-4 py-2 hover:bg-gray-500" @click="showDeleteChannelModal = false">
            Cancel
          </button>
          <button class="bg-red-600 px-4 py-2 hover:bg-red-700" @click="confirmDeletePodcast">
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
