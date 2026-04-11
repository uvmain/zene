<script setup lang="ts">
import { isMobileNavOpen, toggleMobileNav } from '~/logic/navbar'
import { apiKey } from '~/logic/store'
</script>

<template>
  <div v-if="apiKey" class="flex flex-row">
    <!-- Mobile overlay backdrop -->
    <icon-nrk-media-ffw
      class="size-10 rotate-y-0 transition-all duration-300 bottom-4 left-2 absolute z-100 lg:hidden"
      :class="{
        'rotate-y-180 text-accent-400 translate-x-2': isMobileNavOpen,
        'text-primary': !isMobileNavOpen,
      }"
      @click="toggleMobileNav()"
    />

    <!-- Navbar -->
    <aside
      class="p-4 background-2 flex flex-col max-h-100dvh w-full transition-transform duration-300 ease-in-out inset-y-0 left-0 fixed z-50 overflow-y-auto lg:(flex max-w-200px relative)"
      :class="{
        'translate-x-0': isMobileNavOpen,
        '-translate-x-full lg:translate-x-0': !isMobileNavOpen,
      }"
    >
      <div class="flex flex-col space-y-6">
        <RouterLink class="px-2 flex gap-x-2 items-center" to="/">
          <img
            class="opacity-90 size-20 lg:size-12"
            src="/minidisk.svg"
            alt="Logo"
          />
          <div class="text-3xl text-muted font-bold lg:text-2xl">
            Zene
          </div>
        </RouterLink>
        <nav class="px-2 flex flex-col gap-y-4 lg:gap-y-4 md:gap-y-2">
          <NavLink route-name="Home" route-prop="/" />
          <NavLink route-name="Search" route-prop="/search" />
          <NavLink route-name="Albums" route-prop="/albums" />
          <NavLink route-name="Artists" route-prop="/artists" />
          <NavLink route-name="Genres" route-prop="/genres" />
          <NavLink route-name="Tracks" route-prop="/tracks" />
          <NavLink route-name="Playlists" route-prop="/playlists" />
          <NavLink route-name="Radio" route-prop="/radio" />
          <NavLink route-name="Podcasts" route-prop="/podcasts" />
          <NavLink route-name="Queue" route-prop="/queue" />
          <NavLink route-name="Visualizer" route-prop="/visualizer" />
          <NavLink route-name="Settings" route-prop="/settings" />
        </nav>
      </div>
      <NavArt class="max-w-200px" />
    </aside>
  </div>
</template>

<style scoped lang="css">
:deep(a) {
  text-decoration: none;
}
</style>
