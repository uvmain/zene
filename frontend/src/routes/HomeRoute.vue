<script setup lang="ts">
const topTracks = ref()
</script>

<template>
  <div>
    <HeroAlbum />
    <section class="grid grid-cols-2 gap-6">
      <RecentlyAddedAlbums class="cols-1" />
      <TopGenres class="cols-1" />
      <section class="grid grid-cols-2 gap-6">
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
    </section>
  </div>
</template>
