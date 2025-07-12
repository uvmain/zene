<script setup lang="ts">
import { useSearch } from '../composables/useSearch'
import { useNavbar } from '../composables/useNavbar'

const route = useRoute()
const { closeSearch } = useSearch()
const { isMobileNavOpen, closeMobileNav } = useNavbar()

const currentRoute = computed(() => {
  return route.path
})

const handleLinkClick = () => {
  closeSearch()
  closeMobileNav()
}
</script>

<template>
  <!-- Mobile overlay backdrop -->
  <div 
    v-if="isMobileNavOpen"
    class="fixed inset-0 bg-black bg-opacity-50 z-40 md:hidden"
    @click="closeMobileNav"
  />
  
  <!-- Navbar -->
  <aside 
    class="flex flex-col justify-between from-zene-600 to-zene-700 bg-gradient-to-b p-4 transition-transform duration-300 ease-in-out
           fixed inset-y-0 left-0 z-50 w-64 md:w-auto md:relative md:translate-x-0
           md:flex -translate-x-full"
    :class="{
      'translate-x-0': isMobileNavOpen,
      '-translate-x-full md:translate-x-0': !isMobileNavOpen
    }"
  >
    <!-- Mobile close button -->
    <div class="flex justify-end mb-4 md:hidden">
      <button @click="closeMobileNav" class="p-2 text-white hover:text-zene-200 transition-colors">
        <icon-tabler-x class="text-2xl" />
      </button>
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
          to="/playlist"
          class="block flex gap-x-2 rounded-lg px-3 py-2 text-white no-underline transition-all duration-200"
          :class="{ 'ml-4': currentRoute === '/playlists' }"
          @click="handleLinkClick"
        >
          <icon-tabler-playlist />
          Playlists
        </RouterLink>
      </nav>
    </div>
    <NavArt />
  </aside>
</template>
