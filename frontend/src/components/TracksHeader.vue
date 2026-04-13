<template>
  <div
    class="track-grid text-lg text-muted mb-2 pl-2 pr-6"
  >
    <div class="text-center cursor-pointer">
      #
    </div>
    <div class="cursor-pointer">
      Title
    </div>
    <div class="flex cursor-pointer items-center justify-center">
      <icon-nrk-clock class="text-sm" />
    </div>
    <div class="cursor-pointer">
      Album
    </div>
    <div class="hidden cursor-pointer lg:block">
      Genres
    </div>
    <div class="text-center cursor-pointer">
      Year
    </div>
    <div class="flex hidden cursor-pointer cursor-pointer items-center justify-center xl:block">
      <icon-nrk-star class="text-base" />
    </div>
    <div class="hidden cursor-pointer xl:block">
      Plays
    </div>
    <div class="text-center">
      <icon-nrk-more />
    </div>
  </div>
</template>
