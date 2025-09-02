<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicArtist } from '~/types/subsonicArtist'
import type { SubsonicSong } from '~/types/subsonicSong'
import { usePlaybackQueue } from '~/composables/usePlaybackQueue'

defineProps({
  artist: { type: Object as PropType<SubsonicArtist>, required: false },
  album: { type: Object as PropType<SubsonicAlbum>, required: false },
  track: { type: Object as PropType<SubsonicSong>, required: false },
  size: { type: String as PropType<'small' | 'medium' | 'large'>, default: 'medium', required: false },
})

const { play } = usePlaybackQueue()
</script>

<template>
  <div class="flex items-center">
    <button v-if="size === 'small'" class="relative size-10 translate-x-[-1rem] cursor-pointer border-1 border-white rounded-full border-solid bg-zene-600/70 p-2 text-white opacity-0 outline-none transition-all transition-all duration-200 duration-300 group-hover:translate-x-0 hover:from-zene-400 hover:to-zene-200 hover:bg-gradient-to-b group-hover:opacity-100 hover:shadow-sm hover:shadow-zene-400" @click="play(artist, album, track)">
      <icon-tabler-player-play class="text-1xl sm:text-1xl absolute inset-0 m-auto md:text-2xl" />
    </button>
    <button v-if="size === 'medium'" class="w-30 cursor-pointer border-1 border-white rounded-full border-solid bg-zene-600/70 px-4 py-2 text-xl text-white outline-none transition-all duration-200 hover:from-zene-400 hover:to-zene-200 hover:bg-gradient-to-b hover:shadow-sm hover:shadow-zene-400" @click="play(artist, album, track)">
      Play
    </button>
  </div>
</template>
