<script setup lang="ts">
import { useNavbar } from '~/composables/useNavbar'

const { isMobileNavOpen, closeMobileNav } = useNavbar()
</script>

<template>
  <!-- Mobile overlay backdrop -->
  <div
    v-if="isMobileNavOpen"
    class="fixed inset-0 z-40 md:hidden"
    @click="closeMobileNav"
  />

  <!-- Navbar -->
  <aside
    class="fixed inset-y-0 left-0 z-50 w-64 flex flex-col background-2 p-4 transition-transform duration-300 ease-in-out md:relative md:w-auto md:flex"
    :class="{
      'translate-x-0': isMobileNavOpen,
      '-translate-x-full md:translate-x-0': !isMobileNavOpen,
    }"
  >
    <!-- Mobile close button -->
    <div class="mb-4 flex justify-start md:hidden">
      <button
        class="size-12"
        @click="closeMobileNav()"
      >
        <icon-nrk-close class="text-2xl" />
      </button>
    </div>

    <div class="flex flex-col space-y-6">
      <div class="flex items-center gap-x-2 px-2">
        <img class="size-12" src="/logo.png" alt="Logo" />
        <div class="text-2xl font-bold">
          Zene
        </div>
      </div>
      <nav class="flex flex-col gap-y-1 px-2 lg:text-xl">
        <NavLink route-name="Home" route-prop="/" />
        <NavLink route-name="Albums" route-prop="/albums" />
        <NavLink route-name="Tracks" route-prop="/tracks" />
        <NavLink route-name="Artists" route-prop="/artists" />
        <NavLink route-name="Genres" route-prop="/genres" />
        <NavLink route-name="Playlists" route-prop="/playlists" />
        <NavLink route-name="Radio" route-prop="/radio" />
        <NavLink route-name="Podcasts" route-prop="/podcasts" />
      </nav>
    </div>
    <NavArt class="mt-auto" />
  </aside>
</template>
