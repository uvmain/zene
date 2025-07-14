<script setup lang="ts">
import { useNavbar } from '../composables/useNavbar'
import { useSearch } from '../composables/useSearch'

const route = useRoute()
const { closeSearch } = useSearch()
const { isMobileNavOpen, closeMobileNav } = useNavbar()

const currentRoute = computed(() => {
  return route.path
})

function handleLinkClick() {
  closeSearch()
  closeMobileNav()
}
</script>

<template>
  <!-- Mobile overlay backdrop -->
  <div
    v-if="isMobileNavOpen"
    class="fixed inset-0 z-40 bg-black bg-opacity-50 md:hidden"
    @click="closeMobileNav"
  />

  <!-- Navbar -->
  <aside
    class="fixed inset-y-0 left-0 z-50 w-64 flex flex-col from-zene-600 to-zene-700 bg-gradient-to-b p-4 transition-transform duration-300 ease-in-out md:relative md:w-auto md:flex"
    :class="{
      'translate-x-0': isMobileNavOpen,
      '-translate-x-full md:translate-x-0': !isMobileNavOpen,
    }"
  >
    <!-- Mobile close button -->
    <div class="mb-4 flex justify-start md:hidden">
      <icon-tabler-x class="text-2xl text-white transition-colors hover:text-zene-200" @click="closeMobileNav" />
    </div>

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
          @click="handleLinkClick"
        >
          <icon-tabler-home />
          Home
        </RouterLink>
        <RouterLink
          to="/albums"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/albums' }"
          @click="handleLinkClick"
        >
          <icon-tabler-vinyl />
          Albums
        </RouterLink>
        <RouterLink
          to="/tracks"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/tracks' }"
          @click="handleLinkClick"
        >
          <icon-tabler-music />
          Tracks
        </RouterLink>
        <RouterLink
          to="/artists"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/artists' }"
          @click="handleLinkClick"
        >
          <icon-tabler-users-group />
          Artists
        </RouterLink>
        <RouterLink
          to="/genres"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/genres' }"
          @click="handleLinkClick"
        >
          <icon-tabler-tags />
          Genres
        </RouterLink>
        <RouterLink
          to="/playlists"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/playlists' }"
          @click="handleLinkClick"
        >
          <icon-tabler-playlist />
          Playlists
        </RouterLink>
      </nav>
    </div>
    <NavArt class="mt-auto" />
  </aside>
</template>
