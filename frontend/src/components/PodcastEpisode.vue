<script setup lang="ts">
import type { SubsonicPodcastChannelsResponse } from '~/types/subsonic'
import type { SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import { downloadMediaBlob, openSubsonicFetchRequest } from '~/composables/backendFetch'
import { formatTimeFromSeconds } from '~/composables/logic'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'
import { deleteStoredEpisode, episodeIsStored, setStoredEpisode } from '~/stores/usePodcastStore'

const props = defineProps({
  episode: { type: Object as PropType<SubsonicPodcastEpisode>, required: true },
  index: { type: Number, required: true },
})

const emits = defineEmits(['updateEpisodeStatus'])

const { setCurrentlyPlayingPodcastEpisode } = usePlaybackQueue()

const episodeDownloadedLocal = ref(false)
const localDownloadClicked = ref(false)

const episodeArtUrl = computed(() => {
  return `/share/img/${props.episode.coverArt}?size=192`
})

const descriptionLinesCleaned = computed(() => {
  return props.episode.description.split('\n').filter(line => line.trim() !== '').join('<br>')
})

async function downloadEpisode() {
  if (props.episode.status !== 'completed') {
    // download to server
    const formData = new FormData()
    formData.append('id', props.episode.id)
    openSubsonicFetchRequest<SubsonicPodcastChannelsResponse>('downloadPodcastEpisode', {
      body: formData,
    })
    emits('updateEpisodeStatus', props.episode.id, 'downloading')
  }
  else {
    // download to indexedDb
    localDownloadClicked.value = true
    const blob = await downloadMediaBlob(props.episode.streamId)
    await setStoredEpisode(props.episode.streamId, blob)
    updateLocalStorageStatus()
  }
}

async function deleteEpisode() {
  if (await episodeIsStored(props.episode.streamId)) {
    await deleteStoredEpisode(props.episode.streamId)
    updateLocalStorageStatus()
    localDownloadClicked.value = false
  }
  else {
    // delete from server
    const formData = new FormData()
    formData.append('id', props.episode.id)
    openSubsonicFetchRequest<SubsonicPodcastChannelsResponse>('deletePodcastEpisode', {
      body: formData,
    })
    emits('updateEpisodeStatus', props.episode.id, 'new')
  }
}

const episodeStatusButtonText = computed(() => {
  if (props.episode.status === 'completed') {
    return episodeDownloadedLocal.value ? 'Downloaded locally' : 'Downloaded on server'
  }
  else {
    return 'Not downloaded'
  }
})

async function updateLocalStorageStatus() {
  episodeDownloadedLocal.value = await episodeIsStored(props.episode.streamId)
}

watch(() => props.episode.status, async () => {
  updateLocalStorageStatus()
})

onBeforeMount(async () => {
  updateLocalStorageStatus()
})
</script>

<template>
  <div class="corner-cut mx-auto max-w-60dvw flex flex-col gap-4 border-1 border-muted border-solid p-6 hover-background-grad-2">
    <div class="flex flex-row justify-between gap-4">
      <div class="flex flex-row justify-start gap-4">
        <img
          :src="episodeArtUrl"
          alt="Podcast Cover"
          :loading="index < 20 ? 'eager' : 'lazy'"
          class="z-1 col-span-full row-span-full my-auto size-34 rounded object-cover"
          width="192"
          height="192"
        />
        <div class="my-auto ml-1 flex flex-col gap-4">
          <div class="text-lg">
            {{ episode.title }}
          </div>
          <div>
            {{ new Date(episode.publishDate).toLocaleString() }} - {{ formatTimeFromSeconds(Number(episode.duration)) }}
          </div>
          <div class="flex flex-row gap-2">
            <ZButton
              :size12="true"
              @click="downloadEpisode()"
            >
              <Loading
                v-if="episode.status === 'downloading' || localDownloadClicked && !episodeDownloadedLocal"
                class="size-8"
              />
              <icon-nrk-download
                v-else
                class="size-8"
                :class="{
                  'text-orange': episode.status === 'completed' && !episodeDownloadedLocal,
                  'text-green': episode.status === 'completed' && episodeDownloadedLocal,
                }"
              />
            </ZButton>
            <ZButton
              v-if="episode.status === 'completed'"
              :size12="true"
              @click="deleteEpisode()"
            >
              <icon-nrk-trash class="size-8" />
            </ZButton>
            <ZInfo
              v-if="episode.status === 'completed'"
              :text="episodeStatusButtonText"
            >
            </ZInfo>
          </div>
        </div>
      </div>
      <ZButton
        class="my-auto size-14"
        @click="setCurrentlyPlayingPodcastEpisode(episode)"
      >
        <icon-nrk-media-play class="size-12 footer-icon" />
      </ZButton>
    </div>
    <div
      class="line-clamp-4 overflow-hidden text-ellipsis whitespace-normal text-pretty text-op-80"
      v-html="descriptionLinesCleaned"
    />
  </div>
</template>
