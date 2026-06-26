<script setup lang="ts">
import { isMobileNavOpen } from '~/logic/navbar'
import { nowPlayingRoute } from '~/logic/playbackQueue'
import { apiKey } from '~/stores/main'
import NavOpenerMobile from './NavOpenerMobile.vue'
</script>

<template>
  <div v-if="apiKey" class="flex flex-row">
    <aside
      class="p-4 background-2 flex flex-col max-h-100dvh w-full transition-transform duration-300 ease-in-out inset-y-0 left-0 fixed z-50 overflow-y-auto gap-4 lg:(flex w-200px relative)"
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
        <nav class="flex flex-col gap-y-1rem">
          <NavLink route-name="Home" route-prop="/">
            <template #icon>
              <icon-nrk-home class="navicon" />
            </template>
          </NavLink>
          <NavLink route-name="Search" route-prop="/search">
            <template #icon>
              <icon-nrk-search class="navicon" />
            </template>
          </NavLink>
          <NavLink route-name="Albums" route-prop="/albums">
            <template #icon>
              <icon-nrk-media-direkte-notlive class="navicon" />
            </template>
          </NavLink>
          <NavLink route-name="Artists" route-prop="/artists">
            <template #icon>
              <icon-nrk-hardware-microphone class="navicon" />
            </template>
          </NavLink>
          <NavLink route-name="Genres" route-prop="/genres">
            <template #icon>
              <icon-nrk-category class="navicon" />
            </template>
          </NavLink>
          <NavLink route-name="Tracks" route-prop="/tracks">
            <template #icon>
              <icon-nrk-media-quaver class="navicon" />
            </template>
          </NavLink>
          <NavLink route-name="Playlists" route-prop="/playlists">
            <template #icon>
              <icon-nrk-media-playlist class="navicon" />
            </template>
          </NavLink>
          <NavLink route-name="Radio" route-prop="/radio">
            <template #icon>
              <icon-nrk-hardware-radio class="navicon" />
            </template>
          </NavLink>
          <NavLink route-name="Podcasts" route-prop="/podcasts">
            <template #icon>
              <icon-nrk-dialogue class="navicon" />
            </template>
          </NavLink>
          <NavLink route-name="Visualizer" route-prop="/visualizer">
            <template #icon>
              <icon-nrk-media-soundwave class="navicon" />
            </template>
          </NavLink>
          <NavLink route-name="Settings" route-prop="/settings">
            <template #icon>
              <icon-nrk-settings class="navicon" />
            </template>
          </NavLink>
          <NavLink route-name="Now Playing" :route-prop="nowPlayingRoute" :disable-indicator="true">
            <template #icon>
              <icon-nrk-bulleted-list class="navicon" />
            </template>
          </NavLink>
        </nav>
      </div>
      <NavOpenerMobile class="mt-auto relative lg:hidden" />
      <NavArt class="max-w-200px hidden lg:block" />
    </aside>
  </div>
</template>

<style scoped lang="css">
:deep(a) {
  text-decoration: none;
}
.navicon {
  @apply mx-2 size-4 text-muted;
}
</style>
