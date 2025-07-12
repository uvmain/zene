<script setup lang="ts">
import { useSearch } from '../composables/useSearch'

const route = useRoute()
const { closeSearch } = useSearch()

const currentRoute = computed(() => {
  return route.path
})
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
          @click="closeSearch()"
        >
          <icon-tabler-home />
          Home
        </RouterLink>
        <RouterLink
          to="/albums"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/albums' }"
          @click="closeSearch()"
        >
          <icon-tabler-vinyl />
          Albums
        </RouterLink>
        <RouterLink
          to="/tracks"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/tracks' }"
          @click="closeSearch()"
        >
          <icon-tabler-music />
          Tracks
        </RouterLink>
        <RouterLink
          to="/artists"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/artists' }"
          @click="closeSearch()"
        >
          <icon-tabler-users-group />
          Artists
        </RouterLink>
        <RouterLink
          to="/genres"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/genres' }"
          @click="closeSearch()"
        >
          <icon-tabler-tags />
          Genres
        </RouterLink>
        <RouterLink
          to="/playlist"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/playlists' }"
          @click="closeSearch()"
        >
          <icon-tabler-playlist />
          Playlists
        </RouterLink>
        <RouterLink
          to="/debug"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/debug' }"
          @click="closeSearch()"
        >
          <icon-tabler-bug />
          Debug
        </RouterLink>
      </nav>
    </div>
    <NavArt />
  </aside>
</template>
