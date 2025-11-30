<template>
  <div class="flex flex-col">
    <HeroAlbum class="hidden md:block" />
    <Albums :limit="50" sort-key="homePageAlbums" :two-rows-only="true" />
    <Artists :limit="50" sort-key="homePageArtists" :two-rows-only="true" />
    <TopGenres />
  </div>
</template>
