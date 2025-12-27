<script setup lang="ts">
import type { SubsonicPodcastChannelsResponse } from '~/types/subsonic'
import type { SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import { downloadMediaBlob, openSubsonicFetchRequest } from '~/composables/backendFetch'
import { formatTimeFromSeconds } from '~/composables/logic'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'
import { deleteStoredEpisode, episodeIsStored, getStoredEpisode, setStoredEpisode } from '~/stores/usePodcastStore'

const props = defineProps({
  episode: { type: Object as PropType<SubsonicPodcastEpisode>, required: true },
  index: { type: Number, required: true },
})

const emits = defineEmits(['updateEpisodeStatus'])

const { setCurrentlyPlayingPodcastEpisode } = usePlaybackQueue()

const episodeArtUrl = computed(() => {
  return `/share/img/${props.episode.coverArt}?size=192`
})

async function downloadEpisodeOnServer() {
  const formData = new FormData()
  formData.append('id', props.episode.id)
  openSubsonicFetchRequest<SubsonicPodcastChannelsResponse>('downloadPodcastEpisode', {
    body: formData,
  })
  emits('updateEpisodeStatus', props.episode.id, 'downloading')
}

async function downloadEpisode() {
  if (props.episode.status !== 'completed') {
    downloadEpisodeOnServer()
  }
  else {
    const blob = await downloadMediaBlob(props.episode.streamId)
    await setStoredEpisode(props.episode.streamId, blob)
    updateLocalStorageStatus()
  }
}

const episodeDownloadedLocal = ref(false)

const episodeStatusButtonText = computed(() => {
  if (props.episode.status === 'completed') {
    return episodeDownloadedLocal.value ? 'Downloaded Locally' : 'Download locally'
  }
  else {
    return 'Download to server'
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
  <div class="border-primary corner-cut mx-auto max-w-60dvw flex flex-col gap-4 border-1 border-solid p-6 hover-background-grad-2">
    <div class="flex flex-row justify-start gap-4">
      <img
        :src="episodeArtUrl"
        alt="Podcast Cover"
        :loading="index < 20 ? 'eager' : 'lazy'"
        class="z-1 col-span-full row-span-full my-auto size-34 rounded object-cover"
        width="192"
        height="192"
      />
      <div class="my-auto flex flex-col gap-4">
        <div class="text-lg font-semibold">
          {{ episode.title }}
        </div>
        <div>
          {{ new Date(episode.publishDate).toLocaleString() }} - {{ formatTimeFromSeconds(Number(episode.duration)) }}
        </div>
        <div class="flex flex-row gap-2">
          <ZButton
            :size12="true"
            hover-text="play episode"
            @click="setCurrentlyPlayingPodcastEpisode(episode)"
          >
            <icon-nrk-media-play class="size-8 footer-icon" />
          </ZButton>
          <ZButton
            class="flex flex-col items-center gap-1"
            :hover-text="episode.status === 'completed' ? 'downloaded to server' : 'download to server'"
            @click="downloadEpisode()"
          >
            <Loading v-if="episode.status === 'downloading'" class="size-8" />
            <div
              v-else
              class="h-8 flex items-center justify-center text-wrap"
              :class="{
                'text-orange': episode.status === 'completed' && !episodeDownloadedLocal,
                'text-green': episode.status === 'completed' && episodeDownloadedLocal,
              }"
            >
              {{ episodeStatusButtonText }}
            </div>
          </ZButton>
        </div>
      </div>
    </div>
    <div
      class="line-clamp-4 overflow-hidden text-ellipsis whitespace-normal text-pretty text-op-80"
      v-html="episode.description"
    />
  </div>
</template>
