<template>
  <div class="flex flex-col gap-4">
    <HeroAlbum class="hidden md:block" />
    <Albums :limit="20" sort-key="homePageAlbums" />
    <Artists :limit="10" sort-key="homePageArtists" />
    <TopGenres />
  </div>
</template>
