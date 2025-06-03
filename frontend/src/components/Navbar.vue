<script setup lang="ts">
import { currentlyPlayingTrack } from '../composables/globalState'

const route = useRoute()
const router = useRouter()

const currentRoute = computed(() => {
  return route.path
})

function onImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.onerror = null
  target.src = '/default-square.png'
}
</script>

<template>
  <aside class="flex flex-col justify-between from-zene-600 to-zene-700 bg-gradient-to-b p-4">
    <div class="flex flex-col space-y-6">
      <div class="flex items-center justify-center gap-x-2">
        <img class="size-12 rounded-full" src="/logo.png" alt="Logo" />
        <div class="text-2xl font-bold">
          Zene
        </div>
      </div>
      <nav class="flex flex-col gap-y-2 text-xl text-white">
        <RouterLink
          to="/"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/' }"
        >
          <icon-tabler-home />
          Home
        </RouterLink>
        <RouterLink
          to="/albums"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/albums' }"
        >
          <icon-tabler-vinyl />
          Albums
        </RouterLink>
        <RouterLink
          to="/tracks"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/tracks' }"
        >
          <icon-tabler-music />
          Tracks
        </RouterLink>
        <RouterLink
          to="/artists"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/artists' }"
        >
          <icon-tabler-users-group />
          Artists
        </RouterLink>
        <RouterLink
          to="/genres"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/genres' }"
        >
          <icon-tabler-tags />
          Genres
        </RouterLink>
        <RouterLink
          to="/genres"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/playlists' }"
        >
          <icon-tabler-playlist />
          Playlists
        </RouterLink>
      </nav>
    </div>
    <div v-if="currentlyPlayingTrack" class="flex flex-col gap-2">
      <img :src="currentlyPlayingTrack?.image_url" class="w-full cursor-pointer rounded-lg object-cover" @error="onImageError" @click="() => router.push(`/albums/${currentlyPlayingTrack?.musicbrainz_album_id}`)">
      <div class="">
        {{ currentlyPlayingTrack?.title }}
      </div>
      <div class="text-sm text-white/80">
        {{ currentlyPlayingTrack?.artist }}
      </div>
      <div class="text-sm text-white/80">
        {{ currentlyPlayingTrack?.album }}
      </div>
    </div>
  </aside>
</template>
