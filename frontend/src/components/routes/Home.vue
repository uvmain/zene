<script setup lang="ts">
const topTracks = ref()
</script>

<template>
  <div>
    <HeroAlbum />
    <section class="grid grid-cols-2 gap-6">
      <div>
        <RecentlyAddedAlbums />
        <section class="grid grid-cols-2 gap-6">
          <TopGenres />

          <!-- Top Tracks -->
          <div>
            <h2 class="mb-2 text-lg font-semibold">
              Top Tracks
            </h2>
            <ul class="space-y-2">
              <li v-for="(track, i) in topTracks" :key="track.title" class="flex justify-between text-sm">
                <span>{{ i + 1 }}. {{ track.name }}</span>
                <span>{{ track.artist }}</span>
              </li>
            </ul>
          </div>
        </section>
      </div>

      <!-- Player -->
      <div class="rounded bg-gray-800 p-4">
        <div class="text-center">
          <div class="mx-auto mb-2 h-24 w-24 bg-blue-500"></div>
          <div class="font-semibold">
            Track name
          </div>
          <div class="text-sm text-gray-400">
            Artist · Album
          </div>
        </div>
        <div class="mt-4 flex justify-between text-xs">
          <span>2:34</span><span>3:21</span>
        </div>
        <div class="my-1 h-1 rounded bg-gray-600"></div>
        <div class="mt-2 flex justify-center space-x-4">
          <button>
            <icon-tabler-player-skip-back-filled />
          </button>
          <button>
            <icon-tabler-player-play-filled class="text-3xl" />
          </button>
          <button>
            <icon-tabler-player-skip-forward-filled />
          </button>
        </div>
        <button class="mt-2 text-xs underline">
          LYRICS
        </button>
      </div>
    </section>
  </div>
</template>
