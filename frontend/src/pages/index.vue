<template>
  <div class="flex flex-col">
    <HeroAlbum class="hidden md:block" />
    <Albums :limit="50" sort-key="homePageAlbums" :limit-rows="true" />
    <Artists :limit="50" sort-key="homePageArtists" :limit-rows="true" />
    <TopGenres />
  </div>
</template>
