<template>
  <div class="flex flex-col gap-4">
    <HeroAlbum class="hidden lg:block" />
    <Albums sort-key="homePageAlbums" :limit-rows="true" />
    <Artists sort-key="homePageArtists" :limit-rows="true" />
    <TopGenres />
  </div>
</template>
