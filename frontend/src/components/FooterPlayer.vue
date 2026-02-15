<template>
  <footer
    class="sticky bottom-0 mt-auto w-full border-0 border-t-1 border-zshade-600 border-solid background-2"
  >
    <div class="flex flex-col items-center px-8 py-4 lg:flex-row">
      <div
        class="h-full w-full flex flex-grow flex-col items-center justify-center"
      >
        <PlayerAudio />
        <PlayerProgressBar />
        <PlayerMediaControls />
      </div>
      <div class="flex flex-row items-center gap-x-4">
        <!-- <PlayerCastButton /> -->
        <PlayerLyricsButton />
        <PlayerPlaylistButton />
        <PlayerVolumeSlider />
      </div>
    </div>
  </footer>
</template>
