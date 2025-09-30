<template>
  <div class="flex flex-col gap-4">
    <HeroAlbum />
    <Albums :limit="20" />
    <Artists :limit="10" />
    <TopGenres />
  </div>
</template>
