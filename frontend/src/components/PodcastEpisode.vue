<script lang="ts" setup>
import type { SubsonicPodcastChannelsResponse } from '~/types/subsonic'
import type { SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import { openSubsonicFetchRequest } from '~/composables/backendFetch'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'
import { deleteStoredEpisode, getStoredEpisode, setStoredEpisode } from '~/stores/usePodcastStore'

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
</script>

<template>
  <div class="mx-auto max-w-60dvw flex flex-row justify-start gap-4 align-top transition duration-150 hover:scale-101">
    <div class="grid items-end justify-items-end">
      <img
        :src="episodeArtUrl"
        alt="Podcast Cover"
        :loading="index < 20 ? 'eager' : 'lazy'"
        class="z-1 col-span-full row-span-full my-auto h-48 w-48 object-cover"
        width="192"
        height="192"
      />
    </div>
    <div class="flex flex-col justify-between gap-4">
      <div class="text-lg font-semibold">
        {{ episode.title }}
      </div>
      <div
        class="line-clamp-4 overflow-hidden text-ellipsis whitespace-normal text-pretty text-op-80"
        v-html="episode.description.replaceAll(/\n/g, '<br>')"
      />
      <div class="flex flex-row gap-2">
        <ZButton
          :size12="true"
          @click="setCurrentlyPlayingPodcastEpisode(episode)"
        >
          <icon-nrk-media-play class="size-8 footer-icon" />
        </ZButton>
        <ZButton
          :size12="true"
          :disabled="episode.status === 'downloading'"
          @click="downloadEpisodeOnServer(episode)"
        >
          <Loading v-if="episode.status === 'downloading'" />
          <icon-nrk-download
            v-else
            class="size-8"
            :class="{ 'text-orange': episode.status === 'completed' }"
          />
        </ZButton>
      </div>
    </div>
  </div>
</template>
