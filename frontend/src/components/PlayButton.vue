<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicIndexArtist } from '~/types/subsonicArtist'
import type { SubsonicSong } from '~/types/subsonicSong'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'

defineProps({
  artist: { type: Object as PropType<SubsonicIndexArtist>, required: false },
  album: { type: Object as PropType<SubsonicAlbum>, required: false },
  track: { type: Object as PropType<SubsonicSong>, required: false },
})

const { play } = usePlaybackQueue()
</script>

<template>
  <div class="flex items-center align-middle">
    <ZButton
      :primary="true"
      class="group/button"
      :size12="true"
      @click="play(artist, album, track)"
    >
      <icon-nrk-media-play
        class="footer-icon"
      />
    </ZButton>
  </div>
</template>
