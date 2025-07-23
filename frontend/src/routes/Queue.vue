<script setup lang="ts">
import type { TrackMetadataWithImageUrl } from '../types'
import { usePlaybackQueue } from '../composables/usePlaybackQueue'
import { usePlaycounts } from '../composables/usePlaycounts'

const { currentQueue } = usePlaybackQueue()
const { playcount_updated_musicbrainz_track_id } = usePlaycounts()

const tracks = computed(() => currentQueue?.value?.tracks ?? [] as TrackMetadataWithImageUrl[])

watch(playcount_updated_musicbrainz_track_id, (newTrack) => {
  tracks.value?.forEach((track) => {
    if (track.musicbrainz_track_id === newTrack) {
      track.user_play_count = track.user_play_count + 1
      track.global_play_count = track.global_play_count + 1
    }
  })
})
</script>

<template>
  <div>
    <h2 class="px-2 text-lg font-semibold">
      Queue
    </h2>
    <Tracks :tracks="tracks" :show-album="true" />
  </div>
</template>
