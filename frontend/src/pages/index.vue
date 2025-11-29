<template>
  <div class="flex flex-col gap-4">
    <HeroAlbum class="hidden md:block" />
    <Albums :limit="40" sort-key="homePageAlbums" :two-rows-only="true" />
    <Artists :limit="30" sort-key="homePageArtists" :two-rows-only="true" />
    <TopGenres />
  </div>
</template>
