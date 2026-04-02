<template>
  <footer
    class="border-0 border-t-1 border-zshade-600 border-solid background-2 w-full bottom-0 sticky"
  >
    <div class="px-8 py-4 flex flex-col items-center lg:flex-row">
      <div
        class="flex flex-grow flex-col h-full w-full items-center justify-center"
      >
        <PlayerAudio />
        <PlayerProgressBar />
        <PlayerMediaControls />
      </div>
      <div class="flex flex-row gap-4 items-center lg:gap-6">
        <!-- <GoogleCastButton /> -->
        <PlayerVisualizerButton />
        <PlayerLyricsButton />
        <PlayerQueueButton />
        <PlayerVolumeSlider />
      </div>
    </div>
  </footer>
</template>
