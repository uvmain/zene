<script setup lang="ts">
import type { SubsonicPodcastChannelsResponse } from '~/types/subsonic'
import type { SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import { openSubsonicFetchRequest } from '~/composables/backendFetch'
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

async function downloadEpisodeOnServer(episode: SubsonicPodcastEpisode) {
  const formData = new FormData()
  formData.append('id', episode.id)
  openSubsonicFetchRequest<SubsonicPodcastChannelsResponse>('downloadPodcastEpisode', {
    body: formData,
  })
  emits('updateEpisodeStatus', episode.id, 'downloading')
}

const episodeDownloadedLocal = ref(false)

async function updateLocalStorageStatus() {
  episodeDownloadedLocal.value = await episodeIsStored(props.episode.id)
}

watch(() => props.episode.status, async () => {
  updateLocalStorageStatus()
})

onBeforeMount(async () => {
  updateLocalStorageStatus()
})
</script>

<template>
  <div class="corner-cut mx-auto max-w-60dvw flex flex-col gap-4 border-1 border-muted border-solid p-4">
    <div class="flex flex-row justify-start gap-4 align-top transition duration-150 hover:scale-101">
      <img
        :src="episodeArtUrl"
        alt="Podcast Cover"
        :loading="index < 20 ? 'eager' : 'lazy'"
        class="z-1 col-span-full row-span-full my-auto size-40 rounded object-cover"
        width="192"
        height="192"
      />
      <div class="my-auto flex flex-col gap-4">
        <div class="text-lg font-semibold">
          {{ episode.title }}
        </div>
        <div>
          {{ new Date(episode.publishDate).toLocaleString() }}
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
            :size12="true"
            class="flex flex-col items-center gap-1"
            :hover-text="episode.status === 'completed' ? 'downloaded to server' : 'download to server'"
            @click="downloadEpisodeOnServer(episode)"
          >
            <Loading v-if="episode.status === 'downloading'" />
            <icon-nrk-download
              v-else
              class="size-8"
              :class="{
                'text-orange': episode.status === 'completed' && !episodeDownloadedLocal,
                'text-green': episode.status === 'completed' && episodeDownloadedLocal,
              }"
            />
          </ZButton>
        </div>
      </div>
    </div>
    <div
      class="line-clamp-4 overflow-hidden text-ellipsis whitespace-normal text-pretty text-op-80"
      v-html="episode.description.replaceAll(/\n/g, '<br>')"
    />
  </div>
</template>
