<script setup lang="ts">
import { closeMobileNav, isMobileNavOpen } from '~/logic/navbar'
</script>

<template>
  <!-- Mobile overlay backdrop -->
  <div
    v-if="isMobileNavOpen"
    class="inset-0 fixed z-40 lg:hidden"
    @click="closeMobileNav"
  />

  <!-- Navbar -->
  <aside
    class="p-4 background-2 flex flex-col max-h-100dvh w-full transition-transform duration-300 ease-in-out inset-y-0 left-0 fixed z-50 lg:(flex w-auto relative)"
    :class="{
      'translate-x-0': isMobileNavOpen,
      '-translate-x-full lg:translate-x-0': !isMobileNavOpen,
    }"
  >
    <!-- Mobile close button -->
    <div class="mb-4 flex justify-start lg:hidden">
      <ZButton
        :size12="true"
        hover-text="Close Menu"
        @click="closeMobileNav()"
      >
        <icon-nrk-close class="text-2xl" />
      </ZButton>
    </div>

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
      <nav class="px-2 flex flex-col gap-y-4">
        <NavLink route-name="Home" route-prop="/" />
        <NavLink route-name="Albums" route-prop="/albums" />
        <NavLink route-name="Artists" route-prop="/artists" />
        <NavLink route-name="Genres" route-prop="/genres" />
        <NavLink route-name="Tracks" route-prop="/tracks" />
        <NavLink route-name="Playlists" route-prop="/playlists" />
        <NavLink route-name="Radio" route-prop="/radio" />
        <NavLink route-name="Podcasts" route-prop="/podcasts" />
        <NavLink route-name="Queue" route-prop="/queue" />
        <NavLink route-name="Visualizer" route-prop="/visualizer" />
      </nav>
    </div>
    <NavArt class="max-w-50" />
  </aside>
</template>

<style scoped lang="css">
:deep(a) {
  text-decoration: none;
}
</style>
