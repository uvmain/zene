<script setup lang="ts">
import { getRandomTrack } from '../composables/randomTrack'

const topTracks = ref()
const randomTrack = ref()

onMounted(async () => {
  randomTrack.value = await getRandomTrack()
})
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

      <Player :track="randomTrack" />
    </section>
  </div>
</template>
