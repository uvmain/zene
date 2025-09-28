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
      class="group corner-cut size-12 cursor-pointer border-none bg-zene-400 px-3 py-1 text-white outline-none transition-all duration-200 hover:from-zene-400 hover:to-zene-200 hover:bg-gradient-to-b hover:shadow-sm hover:shadow-zene-400"
      @click="play(artist, album, track)"
    >
      <icon-nrk-media-play
        class="translate-x--1 translate-y-1 text-xl transition-all duration-200 group-hover:translate-x-1 md:text-2xl"
      />
    </button>
  </div>
</template>
