<script setup lang="ts">
import { closeMobileNav, isMobileNavOpen } from '~/logic/navbar'
</script>

<template>
  <!-- Mobile overlay backdrop -->
  <div
    v-if="isMobileNavOpen"
    class="fixed inset-0 z-40 lg:hidden"
    @click="closeMobileNav"
  />

  <!-- Navbar -->
  <aside
    class="fixed inset-y-0 left-0 z-50 w-full flex flex-col background-2 p-4 transition-transform duration-300 ease-in-out lg:relative lg:w-auto lg:flex"
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
      <RouterLink class="flex items-center gap-x-2 px-2" to="/">
        <img
          class="size-20 opacity-90 lg:size-12"
          src="/minidisk.svg"
          alt="Logo"
        />
        <div class="text-3xl text-muted font-bold lg:text-2xl">
          Zene
        </div>
      </RouterLink>
      <nav class="flex flex-col gap-y-4 px-2">
        <NavLink route-name="Home" route-prop="/" />
        <NavLink route-name="Albums" route-prop="/albums" />
        <NavLink route-name="Artists" route-prop="/artists" />
        <NavLink route-name="Genres" route-prop="/genres" />
        <NavLink route-name="Tracks" route-prop="/tracks" />
        <NavLink route-name="Playlists" route-prop="/playlists" />
        <NavLink route-name="Radio" route-prop="/radio" />
        <NavLink route-name="Podcasts" route-prop="/podcasts" />
        <NavLink route-name="Queue" route-prop="/queue" />
      </nav>
    </div>
    <NavArt />
  </aside>
</template>

<style scoped lang="css">
:deep(a) {
  text-decoration: none;
}
</style>
