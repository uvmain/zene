<script setup lang="ts">
import type { AlbumMetadata, ArtistMetadata, TrackMetadata, TrackMetadataWithImageUrl } from '../types'
import { getAlbumTracks } from '../composables/fetchFromBackend'
import { setCurrentlyPlayingTrack } from '../composables/globalState'
import { trackWithImageUrl } from '../composables/logic'

const props = defineProps({
  artist: { type: Object as PropType<ArtistMetadata>, required: false },
  album: { type: Object as PropType<AlbumMetadata>, required: false },
  track: { type: Object as PropType<TrackMetadata | TrackMetadataWithImageUrl>, required: false },
})

const trackMetadataWithImageUrl = computed(() => {
  return !props.track ? undefined : trackWithImageUrl(props.track)
})

async function play() {
  if (trackMetadataWithImageUrl.value) {
    setCurrentlyPlayingTrack(trackMetadataWithImageUrl.value)
  }
  else if (props.album) {
    const tracks = await getAlbumTracks(props.album.musicbrainz_album_id)
    setCurrentlyPlayingTrack(tracks[0])
  }
}
</script>

<template>
  <button class="w-30 border-1 border-white rounded-full border-solid bg-zene-600/70 px-4 py-2 text-xl text-white outline-none hover:bg-zene-200" @click="play()">
    Play
  </button>
</template>
