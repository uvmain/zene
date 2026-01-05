<script setup lang="ts">
import type { SubsonicPodcastChannelsResponse } from '~/types/subsonic'
import type { SubsonicPodcastChannel, SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import { openSubsonicFetchRequest, useServerSentEventsForPodcast } from '~/composables/backendFetch'

const route = useRoute()
const router = useRouter()

const showDeleteChannelModal = ref(false)
const showRefreshEpisodesModal = ref(false)

watch(() => route.params.podcast, async () => {
  getPodcast()
})

const podcast = ref<SubsonicPodcastChannel>()

async function getPodcast() {
  const formData = new FormData()
  formData.append('includeEpisodes', true.toString())
  formData.append('id', route.params.podcast.toString())
  const response = await openSubsonicFetchRequest<SubsonicPodcastChannelsResponse>('getPodcasts', {
    body: formData,
  })
  podcast.value = response?.podcasts?.channel[0]
}

const channelCoverArt = computed(() => {
  if (!podcast.value)
    return ''
  return `/share/img/${podcast.value.coverArt}?size=400`
})

const descriptionLinesCleaned = computed(() => {
  return podcast.value?.description.split('\n').filter(line => line.trim() !== '').join('<br>')
})

function confirmDeletePodcast() {
  deletePodcastChannel()
  showDeleteChannelModal.value = false
}

async function refreshPodcastEpisodes() {
  if (!podcast.value)
    return
  const formData = new FormData()
  formData.append('id', podcast.value.id)
  showRefreshEpisodesModal.value = true
  await openSubsonicFetchRequest<SubsonicPodcastChannelsResponse>('refreshPodcast', {
    body: formData,
  })
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

function onMessageReceived(data: any) {
  if (!podcast.value) {
    podcast.value = data[0]
  }
  else if (podcast.value.episode.length < data[0].episode.length) {
    podcast.value.episode = data[0].episode
  }
  else {
    data[0].episode.forEach((newEpisode: SubsonicPodcastEpisode) => {
      const existingEpisode = podcast.value!.episode.find(oldEpisode => oldEpisode.id === newEpisode.id)
      if (existingEpisode) {
        existingEpisode.status = newEpisode.status
      }
    })
  }
}

function onErrorReceived(error: any) {
  console.error('SSE Error Received:', error)
}

function updateEpisodeStatus(episodeId: string, status: string) {
  if (!podcast.value)
    return
  const episode = podcast.value.episode.find(ep => ep.id === episodeId)
  if (episode) {
    episode.status = status
  }
}

onBeforeMount(async () => {
  await getPodcast()
  useServerSentEventsForPodcast(route.params.podcast.toString(), onMessageReceived, onErrorReceived)
})
</script>

<template>
  <div>
    <div v-if="!podcast" class="text-primary">
      Podcast not found.
    </div>
    <div v-else class="mx-auto max-w-60dvw flex flex-col gap-6">
      <!-- header -->
      <div>
        <div class="group relative flex flex-row gap-4 align-top">
          <div class="absolute right-0 flex flex-row gap-2 opacity-0 group-hover:opacity-100">
            <ZButton @click="refreshPodcastEpisodes">
              Refresh episodes
            </ZButton>
            <ZButton @click="showDeleteChannelModal = true">
              Delete podcast channel
            </ZButton>
          </div>
          <img
            :src="channelCoverArt"
            alt="Podcast Cover"
            class="size-70 object-cover"
            width="280"
            height="280"
            loading="eager"
          />
          <div class="my-auto flex flex-col gap-4">
            <div class="text-4xl font-bold">
              {{ podcast.title }}
            </div>
            <div
              class="line-clamp-6 max-h-70 overflow-hidden text-ellipsis whitespace-pre-line text-pretty text-op-80"
              v-html="descriptionLinesCleaned"
            />
            <div v-if="podcast.episode.length && podcast.episode[0].genres?.length > 0" class="flex flex-wrap justify-center gap-2 lg:justify-start">
              <ZInfo v-for="genre in podcast.episode[0].genres?.filter(g => g.name !== '')" :key="genre.name" :text="genre.name" />
            </div>
            <div>
              Source: <a :href="podcast.url" class="text-primary hover:underline" target="_blank">{{ podcast.url }}</a>
            </div>
          </div>
        </div>
      </div>
      <!-- episodes -->
      <div v-if="podcast.lastRefresh === ''">
        Episodes are being refreshed and will appear shortly...
      </div>
      <PodcastEpisode
        v-for="(episode, index) in podcast.episode"
        :key="episode.id"
        :episode="episode"
        :index="index"
        @update-episode-status="updateEpisodeStatus"
      />
    </div>
    <!-- delete channel modal -->
    <Modal
      :show-modal="showDeleteChannelModal"
      modal-text="Are you sure you want to delete this podcast channel?"
    >
      <template #buttons>
        <div class="flex flex-row gap-4">
          <ZButton aria-label="Close" @click="showDeleteChannelModal = false">
            Cancel
          </ZButton>
          <ZButton class="bg-red-600" @click="confirmDeletePodcast">
            Delete
          </ZButton>
        </div>
      </template>
    </Modal>
    <!-- refresh episodes modal -->
    <Modal
      :show-modal="showRefreshEpisodesModal"
      modal-text="Episodes are now being refreshed. Please reload the page later to see updated episodes."
    >
      <template #buttons>
        <ZButton aria-label="Close" @click="showRefreshEpisodesModal = false">
          Okay
        </ZButton>
      </template>
    </Modal>
  </div>
</template>
