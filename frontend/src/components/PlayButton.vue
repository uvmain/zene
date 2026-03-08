<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicIndexArtist } from '~/types/subsonicArtist'
import type { SubsonicPodcastEpisode } from '~/types/subsonicPodcasts'
import type { SubsonicSong } from '~/types/subsonicSong'
import { play } from '~/logic/playbackQueue'

defineProps({
  artist: { type: Object as PropType<SubsonicIndexArtist>, required: false },
  album: { type: Object as PropType<SubsonicAlbum>, required: false },
  track: { type: Object as PropType<SubsonicSong>, required: false },
  podcastEpisode: { type: Object as PropType<SubsonicPodcastEpisode>, required: false },
})
</script>

<template>
  <div class="flex items-center align-middle">
    <ZButton
      :primary="true"
      class="group/button"
      :size12="true"
      :hover-text="`Play ${track ? 'track' : album ? 'album' : artist ? 'artist' : podcastEpisode ? 'podcast' : ''}`"
      @click="play({ artist, album, track, podcastEpisode })"
    >
      <icon-nrk-media-play
        class="footer-icon"
      />
    </ZButton>
  </div>
</template>
