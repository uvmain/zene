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
    <button
      class="z-button size-12"
      @click="play(artist, album, track)"
    >
      <icon-nrk-media-play
        class="text-xl text-primary"
      />
    </button>
  </div>
</template>
